#!/bin/bash

set -e

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $OS in
    linux)
        case $ARCH in
            x86_64) SUFFIX="linux-x64" ;;
            aarch64) SUFFIX="linux-arm64" ;;
            *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
        esac
        ;;
    darwin)
        case $ARCH in
            x86_64) SUFFIX="darwin-x64" ;;
            arm64) SUFFIX="darwin-arm64" ;;
            *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
        esac
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

# Get latest release URL
REPO="NarmadaWeb/goback"
LATEST_URL="https://api.github.com/repos/$REPO/releases/latest"
DOWNLOAD_URL=$(curl -s $LATEST_URL | grep "browser_download_url.*goback-$SUFFIX" | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Could not find download URL for $SUFFIX"
    exit 1
fi

# Download and install
echo "Downloading goback for $SUFFIX..."
curl -L -o goback "$DOWNLOAD_URL"

chmod +x goback

# Install to /usr/local/bin if possible, else ~/bin
if [ -w /usr/local/bin ]; then
    sudo mv goback /usr/local/bin/
else
    mkdir -p ~/bin
    mv goback ~/bin/
    echo "Installed to ~/bin. Make sure ~/bin is in your PATH."
fi

echo "Installation complete. Run 'goback' to start."