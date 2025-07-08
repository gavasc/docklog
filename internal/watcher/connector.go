package watcher

import (
	"context"

	"github.com/docker/docker/client"
)

// spawns a goroutine for each watcher
func Start() {
	events := make(chan Container)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go eventsWatcher(events, ctx)
	go logWatcher(events, ctx)

	select {}
}

func getDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return cli
}
