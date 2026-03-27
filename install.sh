#!/bin/sh
set -e

REPO="IFAKA/coding-typing-tutor"
BINARY="coding-type"
INSTALL_DIR="${CODING_TYPE_INSTALL_DIR:-$HOME/.local/bin}"

# --- Detect OS and architecture ---
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64 | amd64) ARCH="amd64" ;;
  arm64 | aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

case "$OS" in
  linux | darwin) ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

# --- Fetch latest release version ---
echo "Fetching latest version..."
VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' \
  | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
  echo "Could not determine latest version. Is the repo public and does it have releases?"
  exit 1
fi

echo "Installing ${BINARY} ${VERSION} (${OS}/${ARCH})..."

# --- Download and extract ---
ARCHIVE="${BINARY}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"
TMP=$(mktemp -d)

curl -fsSL "$URL" -o "$TMP/$ARCHIVE"
tar -xzf "$TMP/$ARCHIVE" -C "$TMP"
chmod +x "$TMP/$BINARY"

# --- Install ---
mkdir -p "$INSTALL_DIR"
mv "$TMP/$BINARY" "$INSTALL_DIR/$BINARY"
rm -rf "$TMP"

echo ""
echo "Installed to: ${INSTALL_DIR}/${BINARY}"

# --- PATH hint ---
case ":$PATH:" in
  *":${INSTALL_DIR}:"*) ;;
  *)
    echo ""
    echo "Add to your shell config to use from anywhere:"
    echo ""
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
    echo ""
    ;;
esac

echo "Run: ${BINARY}"
echo ""
echo "To uninstall: curl -fsSL https://raw.githubusercontent.com/${REPO}/main/uninstall.sh | sh"
