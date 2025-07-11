# Docklog
A simple Docker log monitor that sends Telegram notifications when errors are detected in your container logs.

## What it does
Docklog watches all your running Docker containers and sends you a Telegram message whenever an error is logged. Perfect for monitoring low-usage containers where you want to be notified of issues without constantly checking logs.

## Features

- 🐳 Monitors all running Docker containers automatically
- 📱 Sends notifications via Telegram, Discord, Slack or all of them
- 🔍 Detects errors in both stdout and stderr streams
- ⬆️ Detects when containers are started, restarted or when they stop and acts accordingly
- 🚀 Lightweight and easy to set up
- 📦 Single binary, no dependencies

## Installation
### Option 1: One-line install (Recommended)
```bash
curl -fsSL https://raw.githubusercontent.com/gavasc/docklog/main/install.sh | bash
```

### Option 2: Download the binary
Go to the releases page and download the binary for your operating system.
Extract and install:
```bash
# Extract the archive
tar -xzf docklog_Linux_x86_64.tar.gz

# Make it executable and move to your PATH
chmod +x docklog
sudo mv docklog /usr/local/bin/
```

## Setup
Docklog supports sending notifications through Telegram and Discord. You can choose one of them or both.

The Discord notifier uses a simple webhook, you can acquire one by going to the channel you want the messages sent to, click on the configurations cog -> Integrations -> Webhooks

The slack notifier also uses a simple webhook, you can see [here](https://api.slack.com/messaging/webhooks) how to create one for your channel

When you install Docklog a config file is created at $HOME/.config/docklog/config.json, use the following pattern to configure your notifiers:

```json
{
  "notifiers": {
        "telegram": {
            "bot_token": "",
            "chat_id": ""
        },
        "discord": {
            "webhook_url": ""
        },
        "slack": {
            "webhook_url": ""
        }
  }
}
```

## Usage
Simply run the binary:
```bash
docklog
```

The program will start monitoring all running Docker containers and send notifications to your prefered chats when errors are 
detected.

## Running as a service
To run docklog as a background service, you can use your system's service manager:
### Using systemd (Linux)
Create a service file at `/etc/systemd/system/docklog.service`:
```ini
[Unit]
Description=Docklog - Docker Log Monitor
After=docker.service
Requires=docker.service

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/docklog
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Then enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable docklog
sudo systemctl start docklog
```

## Error detection
Docklog considers a log line an error if:

- It comes from stderr
- It contains words like "error", "fail", "exception", or "fatal"

## Requirements

- Docker running on your system
- Network access to send Telegram messages
- Telegram bot token and chat ID
