package main

import (
	"os"
	"os/exec"

	"docklog/config"
	"docklog/internal/watcher"
)

func main() {
	if len(os.Args) < 1 {
		// checks if tokens are set in the config file
		config.CheckTokens()

		watcher.Start()
	} else if os.Args[1] == "config" {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi" // fallback default
		}

		cmd := exec.Command(editor, "$HOME/.config/docklog/config.json")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}

		os.Exit(0)
	}
}
