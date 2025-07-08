# Docklog
A simple Docker log monitor that sends Telegram notifications when errors are detected in your container logs.

## What it does
Docklog watches all your running Docker containers and sends you a Telegram message whenever an error is logged. Perfect for monitoring low-usage containers where you want to be notified of issues without constantly checking logs.

## Features

- 🐳 Monitors all running Docker containers automatically
- 📱 Sends notifications via Telegram
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
**1. Create a Telegram Bot**

1. Message @BotFather on Telegram
2. Create a new bot with `/newbot`
3. Save the bot token

**2. Get your Chat ID**

1. Message your bot
2. Visit `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
3. Look for your chat ID in the response

## Usage
Simply run the binary:
```bash
docklog
```

The program will start monitoring all running Docker containers and send notifications to your Telegram chat when errors are 
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
Environment=TELEGRAM_BOT_TOKEN=your_bot_token_here
Environment=TELEGRAM_CHAT_ID=your_chat_id_here

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
