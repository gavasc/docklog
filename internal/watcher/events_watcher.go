package watcher

import (
	"context"
	"docklog/internal/filter"
	"log"

	"github.com/docker/docker/api/types/events"
)

// watches for events from Docker, filters out container start and restart events and sends them to the out channel
func eventsWatcher(e chan<- Container, ctx context.Context) {
	cli := getDockerClient()
	defer cli.Close()
	log.Print("starting events watcher")

	eventCh, errCh := cli.Events(ctx, events.ListOptions{})

	for {
		select {
		case err := <-errCh:
			log.Fatalf("error watching events: %v", err)
			return
		case event := <-eventCh:
			if filter.IsContainerAction(event.Type, event.Action) {
				log.Printf("container event detected: [%v] %ved\n", event.Actor.Attributes["name"], event.Action)
				e <- Container{
					Id: event.Actor.ID,
					Names: []string{
						event.Actor.Attributes["name"],
					},
				}
			}
		}
	}
}
