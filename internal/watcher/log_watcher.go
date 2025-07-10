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
// it also receives new containers to watch, from start and restart events
func logWatcher(e <-chan Container, ctx context.Context) {
	out := make(chan LogEvent)

	log.Println("starting container watchers")

	// lists all running containers at startup
	containers := listActiveContainers()
	log.Printf("found %d containers\n", len(containers))
	for _, container := range containers {
		go watchContainerLogs(container, ctx, out)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-out:
			if filter.IsErrorLog(ev.LogLine, ev.SourceStream) {
				log.Printf("detected error in container [%s]: %s", ev.ContainerName, ev.LogLine)
				notifier.Notify(ev.ContainerName, ev.TimeStamp, ev.SourceStream, ev.LogLine)
			}
		// receives new containers to watch, from start and restart events
		case container := <-e:
			go watchContainerLogs(container, ctx, out)
		}
	}
}

// watches for logs from all running containers, spawning a goroutine for each
func watchContainerLogs(container Container, ctx context.Context, out chan<- LogEvent) {
	cli := getDockerClient()
	defer cli.Close()

	logs, err := cli.ContainerLogs(ctx, container.Id, containertypes.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Since:      "0s",
	})
	log.Printf("watching container [%s]\n", container.Names[0])
	if err != nil {
		log.Printf("failed to get logs for container [%s]: %v", container.Id, err)
		return
	}

	name := container.Names[0]

	err = readDemuxedLogs(ctx, logs, container, out)
	if err != nil {
		log.Printf("failed to read logs for container [%s]: %v", name, err)
	}

	log.Printf("finished reading logs for container [%s]", name)

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

			streamType := getLogStreamType(header[0])
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
func getLogStreamType(streamType byte) string {
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
func listActiveContainers() []Container {
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
