#!/bin/bash

set -e

# Detect OS and architecture
OS_TYPE=$(uname -o | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $OS_TYPE in
    gnu/linux)
        OS="linux"
        EXT=""
        case $ARCH in
            x86_64) SUFFIX="linux-x64" ;;
            aarch64) SUFFIX="linux-arm64" ;;
            *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
        esac
        ;;
    darwin)
        OS="darwin"
        EXT=""
        case $ARCH in
            x86_64) SUFFIX="darwin-x64" ;;
            arm64) SUFFIX="darwin-arm64" ;;
            *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
        esac
        ;;
    msys)
        OS="windows"
        EXT=".exe"
        case $ARCH in
            x86_64) SUFFIX="windows-x64" ;;
            *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
        esac
        ;;
    *)
        echo "Unsupported OS: $OS_TYPE"
        exit 1
        ;;
esac

BINARY_NAME="goback$EXT"

# Get latest release URL
REPO="NarmadaWeb/goback"
LATEST_URL="https://api.github.com/repos/$REPO/releases/latest"
DOWNLOAD_URL=$(curl -s $LATEST_URL | grep "browser_download_url.*goback-$SUFFIX$EXT" | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Could not find download URL for $SUFFIX"
    exit 1
fi

# Download and install
echo "Downloading goback for $SUFFIX..."
curl -L -o "$BINARY_NAME" "$DOWNLOAD_URL"

chmod +x "$BINARY_NAME"

# Install to /usr/local/bin if possible, else ~/bin
if [ "$OS" = "linux" ] || [ "$OS" = "darwin" ]; then
    if [ -w /usr/local/bin ]; then
        sudo mv "$BINARY_NAME" /usr/local/bin/goback
    else
        mkdir -p ~/bin
        mv "$BINARY_NAME" ~/bin/goback
        echo "Installed to ~/bin. Make sure ~/bin is in your PATH."
    fi
else
    # For Windows, assume Git Bash or similar, install to ~/bin
    mkdir -p ~/bin
    mv "$BINARY_NAME" ~/bin/goback.exe
    echo "Installed to ~/bin. Make sure ~/bin is in your PATH."
fi

echo "Installation complete. Run 'goback' to start."