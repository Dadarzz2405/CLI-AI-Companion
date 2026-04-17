#!/bin/sh
set -e

REPO="https://github.com/Dadarzz2405/ai-cli"
RELEASES="https://github.com/Dadarzz2405/ai-cli/releases/latest/download"
BIN_NAME="ai"
INSTALL_DIR="/usr/local/bin"

# detect OS
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Darwin)
    case "$ARCH" in
      arm64) BINARY="ai-darwin-arm64" ;;
      x86_64) BINARY="ai-darwin-amd64" ;;
      *) echo "unsupported arch: $ARCH" && exit 1 ;;
    esac
    ;;
  Linux)
    case "$ARCH" in
      x86_64) BINARY="ai-linux-amd64" ;;
      *) echo "unsupported arch: $ARCH" && exit 1 ;;
    esac
    ;;
  *)
    echo "unsupported OS: $OS"
    echo "windows users: run install.ps1 instead"
    exit 1
    ;;
esac

echo "→ detected $OS/$ARCH"
echo "→ downloading $BINARY..."

curl -fsSL "$RELEASES/$BINARY" -o "/tmp/$BIN_NAME"
chmod +x "/tmp/$BIN_NAME"

echo "→ installing to $INSTALL_DIR/$BIN_NAME (may ask for password)"
sudo mv "/tmp/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"

echo ""
echo "done! run: ai"
echo "first launch will walk you through setup."