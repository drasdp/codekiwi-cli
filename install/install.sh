#!/bin/bash
set -e

# Configuration
GITHUB_REPO="drasdp/codekiwi-cli"
INSTALL_DIR="$HOME/.codekiwi"
BIN_NAME="codekiwi"
COMPOSE_VERSION="docker-compose.yaml"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Detect OS and architecture
detect_system() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    # Map architecture names
    case "$ARCH" in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            ;;
    esac

    # Determine binary name
    case "$OS" in
        darwin)
            BINARY_NAME="${BIN_NAME}-darwin-${ARCH}"
            ;;
        linux)
            BINARY_NAME="${BIN_NAME}-linux-${ARCH}"
            ;;
        *)
            print_error "Unsupported operating system: $OS"
            ;;
    esac

    print_info "Detected system: $OS/$ARCH"
}

# Check Docker installation
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        echo "  macOS: https://www.docker.com/products/docker-desktop"
        echo "  Linux: https://docs.docker.com/engine/install/"
        exit 1
    fi

    if ! docker info &> /dev/null; then
        print_error "Docker daemon is not running. Please start Docker."
        exit 1
    fi

    print_success "Docker is installed and running"
}

# Get latest release URL
get_latest_release() {
    print_info "Fetching latest release..."

    # Get the latest release info from GitHub API
    RELEASE_INFO=$(curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest")

    # Extract the download URL for our binary
    DOWNLOAD_URL=$(echo "$RELEASE_INFO" | grep "browser_download_url.*$BINARY_NAME\"" | cut -d '"' -f 4)

    if [ -z "$DOWNLOAD_URL" ]; then
        print_error "Failed to find release for $BINARY_NAME"
        echo "Please check https://github.com/$GITHUB_REPO/releases"
        exit 1
    fi

    # Extract version
    VERSION=$(echo "$RELEASE_INFO" | grep '"tag_name":' | cut -d '"' -f 4)
    print_info "Latest version: $VERSION"
}

# Download and install binary
install_binary() {
    print_info "Downloading CodeKiwi..."

    # Create install directory
    mkdir -p "$INSTALL_DIR/bin"

    # Download binary
    if ! curl -fsSL "$DOWNLOAD_URL" -o "$INSTALL_DIR/bin/$BIN_NAME"; then
        print_error "Failed to download CodeKiwi"
        exit 1
    fi

    # Make executable
    chmod +x "$INSTALL_DIR/bin/$BIN_NAME"

    print_success "CodeKiwi binary downloaded"
}

# Download configuration files
download_configs() {
    print_info "Downloading configuration files..."

    # Download docker-compose.yaml
    if ! curl -fsSL "https://raw.githubusercontent.com/$GITHUB_REPO/main/$COMPOSE_VERSION" \
        -o "$INSTALL_DIR/$COMPOSE_VERSION"; then
        print_warning "Failed to download $COMPOSE_VERSION"
    fi

    # Download config.env
    if ! curl -fsSL "https://raw.githubusercontent.com/$GITHUB_REPO/main/config.env" \
        -o "$INSTALL_DIR/config.env"; then
        print_warning "Failed to download config.env"
    fi

    print_success "Configuration files downloaded"
}

# Add to PATH
add_to_path() {
    print_info "Configuring PATH..."

    PATH_EXPORT="export PATH=\"\$HOME/.codekiwi/bin:\$PATH\""
    PATH_ADDED=false

    # Function to add to a shell config file
    add_to_shell_config() {
        local config_file="$1"
        if [ -f "$config_file" ]; then
            if ! grep -q ".codekiwi/bin" "$config_file"; then
                echo "" >> "$config_file"
                echo "# CodeKiwi CLI" >> "$config_file"
                echo "$PATH_EXPORT" >> "$config_file"
                print_info "Added to $config_file"
                PATH_ADDED=true
            fi
        fi
    }

    # Add to various shell configs
    add_to_shell_config "$HOME/.bashrc"
    add_to_shell_config "$HOME/.zshrc"
    add_to_shell_config "$HOME/.profile"

    # For macOS, also check bash_profile
    if [ "$OS" = "darwin" ]; then
        add_to_shell_config "$HOME/.bash_profile"
    fi

    if [ "$PATH_ADDED" = true ]; then
        print_success "PATH configuration updated"
    else
        print_info "PATH already configured"
    fi
}

# Create symlink for global access (optional)
create_symlink() {
    # Check if we can write to /usr/local/bin
    if [ -w "/usr/local/bin" ]; then
        if [ ! -e "/usr/local/bin/$BIN_NAME" ]; then
            ln -sf "$INSTALL_DIR/bin/$BIN_NAME" "/usr/local/bin/$BIN_NAME"
            print_success "Created symlink in /usr/local/bin"
        fi
    else
        print_info "Skipping /usr/local/bin symlink (requires sudo)"
    fi
}

# Main installation
main() {
    echo "╔═══════════════════════════════════════╗"
    echo "║     CodeKiwi CLI Installation        ║"
    echo "╚═══════════════════════════════════════╝"
    echo

    # Check Docker first
    check_docker

    # Detect system
    detect_system

    # Get latest release
    get_latest_release

    # Install binary
    install_binary

    # Download configs
    download_configs

    # Configure PATH
    add_to_path

    # Optional: create symlink
    create_symlink

    echo
    print_success "CodeKiwi installed successfully!"
    echo
    echo "To start using CodeKiwi, run one of the following:"
    echo "  1. Restart your terminal"
    echo "  2. Run: source ~/.bashrc (or ~/.zshrc)"
    echo "  3. Run: export PATH=\"\$HOME/.codekiwi/bin:\$PATH\""
    echo
    echo "Then run:"
    echo "  codekiwi --help"
    echo
    echo "To start a project:"
    echo "  codekiwi start [project-directory]"
    echo
}

# Run main function
main