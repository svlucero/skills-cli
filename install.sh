#!/usr/bin/env bash
set -e

# Skills CLI Installation Script
# Usage: curl -fsSL https://raw.githubusercontent.com/svlucero/skills-cli/main/install.sh | bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
REPO="svlucero/skills-cli"
BINARY_NAME="skills"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Functions
print_info() {
    echo -e "${CYAN}→${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin)
            echo "Darwin"
            ;;
        Linux)
            echo "Linux"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            echo "Windows"
            ;;
        *)
            print_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
}

# Detect architecture
detect_arch() {
    local arch
    arch="$(uname -m)"

    case "$arch" in
        x86_64|amd64)
            echo "x86_64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
}

# Get latest release version
get_latest_version() {
    local version
    version=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$version" ]; then
        print_error "Failed to get latest version"
        exit 1
    fi

    echo "$version"
}

# Download and install
install_binary() {
    local os="$1"
    local arch="$2"
    local version="$3"

    print_info "Installing Skills CLI $version for $os/$arch..."

    # Remove 'v' prefix from version for filename
    local version_number="${version#v}"

    # Construct download URL
    local base_url="https://github.com/$REPO/releases/download/$version"
    local filename

    if [ "$os" = "Windows" ]; then
        filename="skills_${version_number}_${os}_${arch}.zip"
    else
        filename="skills_${version_number}_${os}_${arch}.tar.gz"
    fi

    local download_url="$base_url/$filename"

    print_info "Downloading from $download_url..."

    # Create temporary directory
    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap 'rm -rf "$tmp_dir"' EXIT

    # Download and extract binary
    if [ "$os" = "Windows" ]; then
        if ! curl -fsSL "$download_url" -o "$tmp_dir/$filename"; then
            print_error "Failed to download binary"
            exit 1
        fi

        # Extract zip
        unzip -q "$tmp_dir/$filename" -d "$tmp_dir"
        local binary_path="$tmp_dir/skills.exe"
    else
        # Download and extract tar.gz
        if ! curl -fsSL "$download_url" | tar -xz -C "$tmp_dir"; then
            print_error "Failed to download or extract binary"
            exit 1
        fi
        local binary_path="$tmp_dir/$BINARY_NAME"
    fi

    print_success "Binary downloaded successfully"

    # Make binary executable
    chmod +x "$binary_path"

    # Check if install directory is writable
    if [ -w "$INSTALL_DIR" ]; then
        # Install binary
        mv "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
        print_success "Installed to $INSTALL_DIR/$BINARY_NAME"
    else
        # Need sudo
        print_warning "Installing to $INSTALL_DIR requires sudo privileges"
        if ! sudo mv "$binary_path" "$INSTALL_DIR/$BINARY_NAME"; then
            print_error "Failed to install binary"
            exit 1
        fi
        print_success "Installed to $INSTALL_DIR/$BINARY_NAME"
    fi
}

# Verify installation
verify_installation() {
    if ! command -v "$BINARY_NAME" &> /dev/null; then
        print_warning "$BINARY_NAME command not found in PATH"
        print_info "Make sure $INSTALL_DIR is in your PATH"
        print_info "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
        echo "    export PATH=\"$INSTALL_DIR:\$PATH\""
        return 1
    fi

    local installed_version
    installed_version=$("$BINARY_NAME" --version 2>&1 | head -n1)

    print_success "Installation verified: $installed_version"
    return 0
}

# Main
main() {
    echo ""
    echo "╔═══════════════════════════════════════╗"
    echo "║   Skills CLI Installation Script     ║"
    echo "╚═══════════════════════════════════════╝"
    echo ""

    # Check dependencies
    if ! command -v curl &> /dev/null; then
        print_error "curl is required but not installed"
        exit 1
    fi

    # Detect system
    local os arch version
    os=$(detect_os)
    arch=$(detect_arch)

    print_info "Detected system: $os/$arch"

    # Get latest version
    print_info "Fetching latest version..."
    version=$(get_latest_version)
    print_success "Latest version: $version"

    # Install
    install_binary "$os" "$arch" "$version"

    echo ""
    print_success "Skills CLI has been installed successfully!"
    echo ""

    # Verify
    if verify_installation; then
        echo ""
        print_info "Get started with:"
        echo "    skills --help"
        echo "    skills repository add myrepo https://github.com/org/skills-repo.git"
        echo ""
    else
        echo ""
        print_info "Restart your shell or run:"
        echo "    export PATH=\"$INSTALL_DIR:\$PATH\""
        echo ""
    fi
}

# Run main function
main "$@"
