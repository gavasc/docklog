package main

import (
	"docklog/config"
	"docklog/internal/watcher"
)

func main() {
	// checks if tokens are set in the config file
	config.CheckTokens()

	watcher.Start()
}
