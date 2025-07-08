package main

import (
	"docklog/internal/watcher"
	"log"
	"os"
)

func main() {
	if os.Getenv("TELEGRAM_BOT_TOKEN") == "" || os.Getenv("TELEGRAM_CHAT_ID") == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID must be set")
	}

	watcher.Start()

	select {}
}
