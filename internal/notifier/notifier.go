package notifier

import (
	"bytes"
	"docklog/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"slices"
)

func Notify(container string, timestamp string, logStream string, logStr string) {
	str := `Error in container %s at %s
[%s] %s`
	message := fmt.Sprintf(str, container, timestamp, logStream, logStr)

	notifiers := config.GetNotifiers()
	if slices.Contains(notifiers, "telegram") {
		NotifyTelegram(message)
	}
	if slices.Contains(notifiers, "discord") {
		NotifyDiscord(message)
	}

	if slices.Contains(notifiers, "slack") {
		NotifySlack(message)
	}
}

// sends the error notification to a Telegram ID
func NotifyTelegram(message string) {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatId := os.Getenv("TELEGRAM_CHAT_ID")

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s",
		botToken, chatId, url.QueryEscape(message))

	_, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
}

// sends the error notification to a Discord webhook
func NotifyDiscord(message string) {
	webhookUrl := os.Getenv("DISCORD_WEBHOOK_URL")

	body := map[string]string{
		"content": message,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Print("failed to marshal discord body: ", err)
	}

	_, err = http.Post(webhookUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Print("failed to post to discord webhook: ", err)
	}
}

// sends the error notification to a Slack webhook
func NotifySlack(message string) {
	webhookUrl := os.Getenv("SLACK_WEBHOOK_URL")

	body := map[string]string{
		"text": message,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Print("failed to marshal slack body: ", err)
	}

	_, err = http.Post(webhookUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Print("failed to post to slack webhook: ", err)
	}
}
