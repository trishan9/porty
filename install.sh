#!/usr/bin/env bash

set -e

CYAN="\033[36m"
GREEN="\033[32m"
YELLOW="\033[33m"
RED="\033[31m"
MAGENTA="\033[35m"
RESET="\033[0m"

REPO="trishan9/porty"
BIN_NAME="porty"

banner() {
cat << "EOF"

██████╗  ██████╗ ██████╗ ████████╗██╗   ██╗
██╔══██╗██╔═══██╗██╔══██╗╚══██╔══╝╚██╗ ██╔╝
██████╔╝██║   ██║██████╔╝   ██║    ╚████╔╝ 
██╔     ██║   ██║██╔══██╗   ██║     ╚██╔╝  
██╔      ██████╔╝██║  ██║   ██║      ██║   
╚═╝      ╚═════╝ ╚═╝  ╚═╝   ╚═╝      ╚═╝    
     A modern, port manager by @trishan9     
EOF
}

banner
echo -e "${CYAN}→ Starting installation...${RESET}"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')

case "$OS" in
  linux)  OS="linux" ;;
  darwin) OS="darwin" ;;
  *)
    echo -e "${RED}✘ Unsupported OS: $OS${RESET}"
    exit 1
    ;;
esac

ARCH=$(uname -m)

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo -e "${RED}✘ Unsupported architecture: $ARCH${RESET}"
    exit 1
    ;;
esac

echo -e "${GREEN}✓ Detected:${RESET} $OS-$ARCH"

echo -e "${CYAN}→ Checking latest version...${RESET}"

LATEST=$(curl -fsSL https://api.github.com/repos/$REPO/releases/latest \
  | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [[ -z "$LATEST" ]]; then
    echo -e "${RED}✘ Failed to fetch latest release information${RESET}"
    exit 1
fi

echo -e "${GREEN}✓ Latest version:${RESET} $LATEST"

FILE="${BIN_NAME}-${OS}-${ARCH}"
URL="https://github.com/$REPO/releases/download/$LATEST/$FILE"

echo -e "${CYAN}→ Downloading binary:${RESET} $FILE"


TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

curl -fsSL "$URL" -o "$TEMP_DIR/$BIN_NAME"
chmod +x "$TEMP_DIR/$BIN_NAME"

if command -v sudo &> /dev/null; then
  BIN_PATH="/usr/local/bin/$BIN_NAME"
  echo -e "${CYAN}→ Installing to $BIN_PATH (sudo required)...${RESET}"
  sudo mv "$TEMP_DIR/$BIN_NAME" "$BIN_PATH"

  # Check PATH
  if ! echo "$PATH" | grep -q "/usr/local/bin"; then
    echo -e "${YELLOW}⚠ Your PATH does not include /usr/local/bin${RESET}"
    echo -e "${CYAN}→ Automatically adding it to your shell config...${RESET}"

    if [ -f "$HOME/.zshrc" ]; then
      echo 'export PATH="/usr/local/bin:$PATH"' >> "$HOME/.zshrc"
      echo -e "${GREEN}✓ Added to ~/.zshrc${RESET}"
    fi

    if [ -f "$HOME/.bashrc" ]; then
      echo 'export PATH="/usr/local/bin:$PATH"' >> "$HOME/.bashrc"
      echo -e "${GREEN}✓ Added to ~/.bashrc${RESET}"
    fi

    if command -v fish >/dev/null; then
      fish -c 'set -U fish_user_paths /usr/local/bin $fish_user_paths'
      echo -e "${GREEN}✓ Added to fish PATH${RESET}"
    fi

    echo -e "${MAGENTA}Restart your terminal to apply changes.${RESET}"
  fi
else
    INSTALL_PATH="$HOME/.local/bin/$BIN_NAME"
    mkdir -p "$HOME/.local/bin"
    mv "$TEMP_DIR/$BIN_NAME" "$INSTALL_PATH"
    echo -e "${YELLOW}⚠ ~/.local/bin added locally. Ensure it's in your PATH.${RESET}"
fi


echo
echo -e "${GREEN}✓ Installation complete!${RESET}"
echo -e "${MAGENTA}Installed at:${RESET} $INSTALL_PATH"
echo
echo -e "${CYAN}→ Testing porty...${RESET}"

if ! "$INSTALL_PATH" --help >/dev/null 2>&1; then
    echo -e "${RED}✘ Something went wrong. Please check permissions.${RESET}"
    exit 1
fi

echo -e "${GREEN}✓ Porty is ready!${RESET}"
echo
echo -e "${MAGENTA}Run it now:${RESET}   porty list"
echo -e "${MAGENTA}Docs:${RESET}        https://github.com/$REPO"
echo

