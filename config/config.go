package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Notifiers struct {
		Telegram struct {
			BotToken string `json:"bot_token"`
			ChatId   string `json:"chat_id"`
		} `json:"telegram"`

		Discord struct {
			WebhookUrl string `json:"webhook_url"`
		} `json:"discord"`
	} `json:"notifiers"`
}

func loadConfig() Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	configFile, err := os.Open(homeDir + "/.config/docklog/config.json")
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}

	config := Config{}

	decoder := json.NewDecoder(configFile)
	if err = decoder.Decode(&config); err != nil {
		log.Fatalf("Error decoding config file: %v", err)
	}

	return config
}

// checks which tokens are set in the config file and sets them as environment variables
func CheckTokens() {
	config := loadConfig()

	// checking if all the tokens are empty
	if config.Notifiers.Telegram.BotToken == "" && config.Notifiers.Telegram.ChatId == "" && config.Notifiers.Discord.WebhookUrl == "" {
		log.Fatal("No tokens found in config file!")
	}

	setTokens()
}

// sets the tokens as environment variables
func setTokens() {
	config := loadConfig()

	if config.Notifiers.Telegram.BotToken != "" && config.Notifiers.Telegram.ChatId != "" {
		os.Setenv("TELEGRAM_BOT_TOKEN", config.Notifiers.Telegram.BotToken)
		os.Setenv("TELEGRAM_CHAT_ID", config.Notifiers.Telegram.ChatId)
	}

	if config.Notifiers.Discord.WebhookUrl != "" {
		os.Setenv("DISCORD_WEBHOOK_URL", config.Notifiers.Discord.WebhookUrl)
	}
}

// sees which notifiers are present and returns their names
func GetNotifiers() []string {
	config := loadConfig()
	present := []string{}

	if config.Notifiers.Telegram.BotToken != "" && config.Notifiers.Telegram.ChatId != "" {
		present = append(present, "telegram")
	}

	if config.Notifiers.Discord.WebhookUrl != "" {
		present = append(present, "discord")
	}

	return present
}
