package watcher

import (
	"context"
	"docklog/internal/filter"
	"docklog/internal/notifier"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"time"

	containertypes "github.com/docker/docker/api/types/container"
)

type Container struct {
	Id    string
	Names []string
}

type LogEvent struct {
	ContainerId   string
	ContainerName string
	TimeStamp     string
	LogLine       string
	SourceStream  string
}

// starts the log watcher, receiving logs from watchContainers through the out channel
func Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := make(chan LogEvent)

	err := watchContainers(ctx, out)
	if err != nil {
		log.Printf("failed to watch containers: %v", err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case e := <-out:
			if filter.IsError(e.LogLine, e.SourceStream) {
				notifier.Notify(e.ContainerName, e.TimeStamp, e.SourceStream, e.LogLine)
			}
		}
	}
}

// watches for logs from all running containers, spawning a goroutine for each
func watchContainers(ctx context.Context, out chan<- LogEvent) error {
	cli := getDockerClient()
	defer cli.Close()

	containers := listContainers()
	for _, container := range containers {
		go func(c Container) {
			logs, err := cli.ContainerLogs(ctx, container.Id, containertypes.LogsOptions{
				ShowStdout: true,
				ShowStderr: true,
				Follow:     true,
				Since:      "0s",
			})
			if err != nil {
				log.Printf("failed to get logs for container %s: %v", c.Id, err)
				return
			}

			name := c.Names[0]

			err = readDemuxedLogs(ctx, logs, c, out)
			if err != nil {
				log.Printf("failed to read logs for container %s: %v", name, err)
			}
		}(container)
	}

	return nil
}

// reads logs from a multiplexed stream, demultiplexing stdout and stderr and returning them as LogEvents
func readDemuxedLogs(ctx context.Context, stream io.Reader, container Container, out chan<- LogEvent) error {
	header := make([]byte, 8)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, err := io.ReadFull(stream, header)
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("failed to read header: %w", err)
			}

			size := binary.BigEndian.Uint32(header[4:8])
			if size == 0 {
				continue
			}

			streamType := getStreamType(header[0])
			payload := make([]byte, size)
			_, err = io.ReadFull(stream, payload)
			if err != nil {
				return fmt.Errorf("failed to read payload: %w", err)
			}

			out <- LogEvent{
				ContainerId:   container.Id,
				ContainerName: container.Names[0],
				TimeStamp:     time.Now().Format("2006-01-02 15:04:05"),
				LogLine:       string(payload),
				SourceStream:  streamType,
			}
		}
	}
}

// returns a string representation of the stream type
func getStreamType(streamType byte) string {
	switch streamType {
	case 1:
		return "stdout"
	case 2:
		return "stderr"
	default:
		return "unknown"
	}
}

// returns a list of running containers, with their ID and names
func listContainers() []Container {
	cli := getDockerClient()
	defer cli.Close()

	containers, err := cli.ContainerList(context.Background(), containertypes.ListOptions{})
	if err != nil {
		panic(err)
	}

	containerList := make([]Container, 0)
	for _, container := range containers {
		containerList = append(containerList, Container{Id: container.ID, Names: container.Names})
	}

	return containerList
}
