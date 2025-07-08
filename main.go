package main

import (
	"docklog/internal/watcher"
	"log"
	"os"
)

func main() {
	// checking if TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID are set
	if os.Getenv("TELEGRAM_BOT_TOKEN") == "" || os.Getenv("TELEGRAM_CHAT_ID") == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID must be set")
	}
	log.Println("Env variables set.\nStarting docklog...")

	watcher.Start()
}
