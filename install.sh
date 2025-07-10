#!/bin/bash

# Docklog installer script

set -e

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Convert architecture names to match GoReleaser naming
case $ARCH in
    x86_64) ARCH="x86_64" ;;
    amd64) ARCH="x86_64" ;;
    arm64) ARCH="arm64" ;;
    aarch64) ARCH="arm64" ;;
    i386) ARCH="i386" ;;
    i686) ARCH="i386" ;;
    *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

# Convert OS names to match GoReleaser naming
case $OS in
    linux) OS="Linux" ;;
    darwin) OS="Darwin" ;;
    windows) OS="Windows" ;;
    *) echo "Unsupported OS: $OS" && exit 1 ;;
esac

REPO="gavasc/docklog"
BINARY_NAME="docklog"

# Get latest release
echo "Fetching latest release information..."
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "Failed to fetch latest release"
    exit 1
fi

echo "Latest release: $LATEST_RELEASE"

# Construct download URL
if [ "$OS" = "Windows" ]; then
    ARCHIVE_NAME="${BINARY_NAME}_${OS}_${ARCH}.zip"
else
    ARCHIVE_NAME="${BINARY_NAME}_${OS}_${ARCH}.tar.gz"
fi

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/$ARCHIVE_NAME"

TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

echo "Downloading $ARCHIVE_NAME..."
curl -L -o "$ARCHIVE_NAME" "$DOWNLOAD_URL"

echo "Extracting..."
if [ "$OS" = "Windows" ]; then
    unzip "$ARCHIVE_NAME"
else
    tar -xzf "$ARCHIVE_NAME"
fi

INSTALL_DIR="/usr/local/bin"
if [ "$OS" = "Windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
    INSTALL_DIR="/usr/bin"  # For WSL/Git Bash
fi

echo "Installing $BINARY_NAME to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY_NAME" "$INSTALL_DIR/"
else
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
fi

# Make executable (not needed on Windows)
if [ "$OS" != "Windows" ]; then
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

# Create config file
echo "Creating config file..."
mkdir -p "$HOME/.config/docklog"
if [ ! -f "$HOME/.config/docklog/config.json" ]; then
    touch "$HOME/.config/docklog/config.json"
    echo "Config file created in $HOME/.config/docklog/config.json"
else
    echo "Config file already exists in $HOME/.config/docklog/config.json"
fi

# Cleanup
cd /
rm -rf "$TMP_DIR"

echo "Docklog installed successfully!"
echo "Run 'docklog' to start monitoring your Docker containers."