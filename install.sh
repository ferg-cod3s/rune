#!/bin/bash

# Rune CLI Installation Script
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="johnferguson/rune"
BINARY_NAME="rune"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="x86_64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo -e "${RED}Unsupported architecture: $ARCH${NC}"; exit 1 ;;
esac

case $OS in
    linux) OS="Linux" ;;
    darwin) OS="Darwin" ;;
    *) echo -e "${RED}Unsupported OS: $OS${NC}"; exit 1 ;;
esac

echo -e "${GREEN}Installing Rune CLI...${NC}"
echo "OS: $OS"
echo "Architecture: $ARCH"

# Get latest release
echo -e "${YELLOW}Fetching latest release...${NC}"
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo -e "${RED}Failed to fetch latest release${NC}"
    exit 1
fi

echo "Latest version: $LATEST_RELEASE"

# Construct download URL
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/${BINARY_NAME}_${OS}_${ARCH}.tar.gz"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download and extract
echo -e "${YELLOW}Downloading $DOWNLOAD_URL...${NC}"
curl -sL "$DOWNLOAD_URL" | tar xz

# Check if binary exists
if [ ! -f "$BINARY_NAME" ]; then
    echo -e "${RED}Binary not found in archive${NC}"
    exit 1
fi

# Make binary executable
chmod +x "$BINARY_NAME"

# Install binary
echo -e "${YELLOW}Installing to $INSTALL_DIR...${NC}"
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY_NAME" "$INSTALL_DIR/"
else
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
fi

# Install shell completions if available
if [ -d "completions" ]; then
    echo -e "${YELLOW}Installing shell completions...${NC}"
    
    # Bash completion
    if [ -f "completions/rune.bash" ]; then
        if [ -d "/usr/share/bash-completion/completions" ]; then
            sudo cp "completions/rune.bash" "/usr/share/bash-completion/completions/rune" 2>/dev/null || true
        elif [ -d "/etc/bash_completion.d" ]; then
            sudo cp "completions/rune.bash" "/etc/bash_completion.d/rune" 2>/dev/null || true
        fi
    fi
    
    # Zsh completion
    if [ -f "completions/rune.zsh" ]; then
        if [ -d "/usr/share/zsh/site-functions" ]; then
            sudo cp "completions/rune.zsh" "/usr/share/zsh/site-functions/_rune" 2>/dev/null || true
        fi
    fi
    
    # Fish completion
    if [ -f "completions/rune.fish" ]; then
        if [ -d "/usr/share/fish/completions" ]; then
            sudo cp "completions/rune.fish" "/usr/share/fish/completions/" 2>/dev/null || true
        fi
    fi
fi

# Cleanup
cd /
rm -rf "$TMP_DIR"

# Verify installation
if command -v "$BINARY_NAME" >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Rune CLI installed successfully!${NC}"
    echo ""
    echo "Run 'rune --help' to get started."
    echo "Run 'rune init' to create your first configuration."
else
    echo -e "${RED}Installation failed. Please check that $INSTALL_DIR is in your PATH.${NC}"
    exit 1
fi