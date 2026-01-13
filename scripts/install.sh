#!/bin/bash
#
# CCM Quick Install Script
#
# Usage: curl -fsSL https://raw.githubusercontent.com/taliove/go-claude-model/main/scripts/install.sh | bash
#

set -e

REPO="taliove/go-claude-model"
BINARY_NAME="ccm"
INSTALL_DIR="${CCM_INSTALL_DIR:-$HOME/.local/bin}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[OK]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }

# Check dependencies
check_dependencies() {
    echo ""
    log_info "Checking dependencies..."

    # Check npm
    if command -v npm &> /dev/null; then
        log_success "npm is installed"
    else
        log_warn "npm is not installed"
        echo "  Install Node.js: https://nodejs.org/"
    fi

    # Check claude
    if command -v claude &> /dev/null; then
        log_success "claude-code is installed"
    else
        log_warn "claude-code is not installed"
        echo "  Install: npm install -g @anthropic-ai/claude-code"
    fi
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case $ARCH in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $ARCH"
            ;;
    esac

    case $OS in
        linux|darwin)
            ;;
        mingw*|msys*|cygwin*)
            OS="windows"
            ;;
        *)
            log_error "Unsupported OS: $OS"
            ;;
    esac

    echo "${OS}_${ARCH}"
}

# Get latest version
get_latest_version() {
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

main() {
    log_info "Installing CCM..."

    PLATFORM=$(detect_platform)
    VERSION=$(get_latest_version)

    if [ -z "$VERSION" ]; then
        log_error "Failed to get latest version"
    fi

    log_info "Detected platform: $PLATFORM"
    log_info "Latest version: $VERSION"

    # Download URL
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}_${VERSION#v}_${PLATFORM}.tar.gz"

    log_info "Downloading from: $DOWNLOAD_URL"

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT

    # Download and extract
    curl -fsSL "$DOWNLOAD_URL" | tar xz -C "$TMP_DIR"

    # Install
    mkdir -p "$INSTALL_DIR"
    mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"

    log_success "CCM installed to $INSTALL_DIR/$BINARY_NAME"

    # Check PATH
    if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
        echo ""
        log_info "Add to your PATH:"
        echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    fi

    echo ""
    echo "=========================================="
    echo "  Getting Started"
    echo "=========================================="
    echo ""
    echo "1. Run initial setup:"
    echo "   ccm init"
    echo ""
    echo "2. Configure a provider:"
    echo "   ccm add doubao --key \"your-api-key\""
    echo ""
    echo "3. Start Claude Code:"
    echo "   ccm run doubao"
    echo ""

    # Check dependencies
    check_dependencies

    echo ""
    log_success "Installation complete!"
}

main
