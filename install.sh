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

# Install to user local bin directory
INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"
mv "$BINARY_NAME" "$INSTALL_DIR/goback"
echo "Installed to $INSTALL_DIR. Make sure $INSTALL_DIR is in your PATH."

echo "Installation complete. Run 'goback' to start."