package notifier

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// sends the error notification to a Telegram ID
func Notify(container string, timestamp string, logStream string, logStr string) {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatId := os.Getenv("TELEGRAM_CHAT_ID")

	str := `Error in container %s at %s
[%s] %s`

	message := fmt.Sprintf(str, container, timestamp, logStream, logStr)

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s",
		botToken, chatId, url.QueryEscape(message))

	_, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
}
