#!/usr/bin/env bash

#############################################################################
# Mobile Recon Toolkit - Installation Script
#############################################################################
# This script builds and installs the mobile-recon CLI globally.
#
# Usage:
#   ./scripts/install.sh              # Install the toolkit
#   ./scripts/install.sh --verbose    # Install with verbose output
#   ./scripts/install.sh --help       # Show help
#
# What it does:
#   1. Checks prerequisites (Go installation)
#   2. Cleans up old installation (removes mobile-recon-cli, old aliases)
#   3. Builds the unified mobile-recon CLI (includes all tools)
#   4. Installs it to $GOPATH/bin or $HOME/go/bin
#   5. Verifies PATH configuration
#   6. Creates helpful aliases (optional)
#############################################################################

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
GO_TOOLS_DIR="$PROJECT_ROOT/go-tools"
CLI_DIR="$GO_TOOLS_DIR/mobile-recon-cli"

# Verbose mode
VERBOSE=false

#############################################################################
# Helper Functions
#############################################################################

print_header() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
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

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_step() {
    echo -e "${BLUE}▶${NC} $1"
}

#############################################################################
# Check Prerequisites
#############################################################################

check_prerequisites() {
    print_header "Checking Prerequisites"

    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed!"
        echo ""
        echo "Please install Go from: https://golang.org/dl/"
        echo "Minimum required version: 1.21"
        exit 1
    fi

    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_success "Go is installed (version $GO_VERSION)"

    # Check GOPATH
    if [ -z "$GOPATH" ]; then
        GOPATH="$HOME/go"
        print_info "GOPATH not set, using default: $GOPATH"
    else
        print_success "GOPATH: $GOPATH"
    fi

    # Check if GOPATH/bin is in PATH
    if [[ ":$PATH:" != *":$GOPATH/bin:"* ]]; then
        print_warning "GOPATH/bin is not in your PATH"
        setup_path
    else
        print_success "GOPATH/bin is in PATH"
    fi
}

#############################################################################
# Setup PATH
#############################################################################

setup_path() {
    echo ""
    echo "  Would you like to add Go bin to your PATH automatically?"
    echo ""
    read -p "  Add PATH configuration? (y/N): " -n 1 -r
    echo ""

    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # Detect shell config file
        local shell_config=""
        if [ -f "$HOME/.zshrc" ]; then
            shell_config="$HOME/.zshrc"
        elif [ -f "$HOME/.bashrc" ]; then
            shell_config="$HOME/.bashrc"
        else
            echo "  Could not detect shell config. Please add manually:"
            echo -e "  ${YELLOW}export PATH=\"\$HOME/go/bin:\$PATH\"${NC}"
            echo ""
            return
        fi

        # Check if already configured
        if grep -q 'export PATH=.*go/bin' "$shell_config" 2>/dev/null; then
            print_info "PATH already configured in $shell_config"
            return
        fi

        echo "" >> "$shell_config"
        echo "# Go bin PATH (added by mobile-recon install script)" >> "$shell_config"
        echo 'export PATH="$HOME/go/bin:$PATH"' >> "$shell_config"

        print_success "PATH added to $shell_config"
        echo ""
        echo "  Reload your shell to apply:"
        echo -e "  ${YELLOW}source $shell_config${NC}"
    else
        print_info "Skipping PATH setup"
        echo ""
        echo "  Add this manually to your shell config:"
        echo -e "  ${YELLOW}export PATH=\"\$HOME/go/bin:\$PATH\"${NC}"
    fi
    echo ""
}

#############################################################################
# Cleanup Old Installation
#############################################################################

cleanup_old_installation() {
    print_header "Cleaning Up Old Installation"

    local cleaned=false

    # Remove old mobile-recon-cli binary
    if [ -f "$GOPATH/bin/mobile-recon-cli" ]; then
        print_step "Removing old mobile-recon-cli binary..."
        rm -f "$GOPATH/bin/mobile-recon-cli"
        print_success "Removed mobile-recon-cli"
        cleaned=true
    fi

    # Remove old separate toolkit binaries
    local old_binaries=("adb-toolkit" "nmap-toolkit" "apk-analyzer" "ios-toolkit")
    for binary in "${old_binaries[@]}"; do
        if [ -f "$GOPATH/bin/$binary" ]; then
            print_step "Removing old $binary binary..."
            rm -f "$GOPATH/bin/$binary"
            print_success "Removed $binary"
            cleaned=true
        fi
    done

    # Remove old alias from shell config
    local shell_configs=("$HOME/.zshrc" "$HOME/.bashrc")
    for config in "${shell_configs[@]}"; do
        if [ -f "$config" ]; then
            # Check for old alias pattern
            if grep -q "alias mobile-recon='mobile-recon-cli'" "$config" 2>/dev/null; then
                print_step "Removing old alias from $config..."
                # Remove the old alias line
                sed -i '' "/alias mobile-recon='mobile-recon-cli'/d" "$config" 2>/dev/null || \
                sed -i "/alias mobile-recon='mobile-recon-cli'/d" "$config" 2>/dev/null
                print_success "Removed old alias from $config"
                cleaned=true
            fi
        fi
    done

    if [ "$cleaned" = true ]; then
        print_info "Old installation cleaned up"
    else
        print_info "No old installation found"
    fi

    echo ""
}

#############################################################################
# Build and Install
#############################################################################

install_mobile_recon() {
    print_header "Building and Installing Mobile Recon"

    local is_upgrade=false
    local install_path="$GOPATH/bin/mobile-recon"

    # Check if already installed
    if [ -f "$install_path" ]; then
        is_upgrade=true
        print_info "mobile-recon is already installed at: $install_path"
        print_step "Replacing with new build..."
    else
        print_step "Building mobile-recon..."
    fi

    cd "$CLI_DIR"

    # Download dependencies
    print_step "Downloading dependencies..."
    if [ "$VERBOSE" = true ]; then
        go mod tidy
    else
        go mod tidy > /dev/null 2>&1
    fi

    # Build
    print_step "Compiling..."
    if [ "$VERBOSE" = true ]; then
        go build -o mobile-recon
    else
        go build -o mobile-recon > /dev/null 2>&1
    fi

    if [ $? -eq 0 ]; then
        print_success "Built mobile-recon successfully"
    else
        print_error "Failed to build mobile-recon"
        return 1
    fi

    # Remove existing binary if upgrading
    if [ "$is_upgrade" = true ]; then
        rm -f "$install_path"
    fi

    # Install
    print_step "Installing to $GOPATH/bin..."
    mv mobile-recon "$install_path"

    if [ $? -eq 0 ]; then
        if [ "$is_upgrade" = true ]; then
            print_success "Replaced mobile-recon at $install_path"
        else
            print_success "Installed mobile-recon to $install_path"
        fi
    else
        print_error "Failed to install mobile-recon"
        return 1
    fi

    echo ""
    return 0
}

#############################################################################
# Verify Installation
#############################################################################

verify_installation() {
    print_header "Verifying Installation"

    if command -v mobile-recon &> /dev/null; then
        local tool_path=$(which mobile-recon)
        print_success "mobile-recon found at: $tool_path"

        echo ""
        print_info "Installed tools (all built into mobile-recon):"
        echo "  • mobile-recon adb   - Android Debug Bridge toolkit"
        echo "  • mobile-recon nmap  - Network reconnaissance toolkit"
        echo "  • mobile-recon apk   - Android APK static analysis"
        echo "  • mobile-recon ios   - iOS Simulator toolkit"
    else
        print_warning "mobile-recon not found in PATH"
        echo "  Expected location: $GOPATH/bin/mobile-recon"
        echo ""
        echo "  Make sure to add this to your shell configuration:"
        echo -e "  ${YELLOW}export PATH=\"\$HOME/go/bin:\$PATH\"${NC}"
        echo ""
        echo "  Then reload your shell:"
        echo -e "  ${YELLOW}source ~/.bashrc${NC}  # for bash"
        echo -e "  ${YELLOW}source ~/.zshrc${NC}   # for zsh"
    fi

    echo ""
}

#############################################################################
# Setup Aliases
#############################################################################

setup_aliases() {
    print_header "Setting Up Aliases (Optional)"

    echo "Would you like to add helpful aliases to your shell configuration?"
    echo ""
    echo "Suggested aliases:"
    echo "  alias mr='mobile-recon'"
    echo "  alias mradb='mobile-recon adb'"
    echo "  alias mrnmap='mobile-recon nmap'"
    echo "  alias mrapk='mobile-recon apk'"
    echo "  alias mrios='mobile-recon ios'"
    echo ""

    read -p "Add aliases? (y/N): " -n 1 -r
    echo ""

    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # Detect shell config file
        local shell_config=""
        if [ -f "$HOME/.zshrc" ]; then
            shell_config="$HOME/.zshrc"
        elif [ -f "$HOME/.bashrc" ]; then
            shell_config="$HOME/.bashrc"
        else
            echo "Could not detect shell config. Please add aliases manually."
            return
        fi

        # Check if aliases already exist
        if grep -q "alias mr='mobile-recon'" "$shell_config" 2>/dev/null; then
            print_info "Aliases already exist in $shell_config"
            return
        fi

        echo "" >> "$shell_config"
        echo "# Mobile Recon Toolkit Aliases (added by install script)" >> "$shell_config"
        echo "alias mr='mobile-recon'" >> "$shell_config"
        echo "alias mradb='mobile-recon adb'" >> "$shell_config"
        echo "alias mrnmap='mobile-recon nmap'" >> "$shell_config"
        echo "alias mrapk='mobile-recon apk'" >> "$shell_config"
        echo "alias mrios='mobile-recon ios'" >> "$shell_config"

        print_success "Aliases added to $shell_config"
        echo ""
        echo "  Reload your shell to use them:"
        echo -e "  ${YELLOW}source $shell_config${NC}"
    else
        print_info "Skipping alias setup"
    fi
}

#############################################################################
# Usage Information
#############################################################################

show_usage() {
    print_header "Mobile Recon Toolkit - Installation Script"

    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -v, --verbose    Show detailed build output"
    echo "  -h, --help       Show this help message"
    echo ""
    echo "What this script does:"
    echo "  1. Checks prerequisites (Go installation)"
    echo "  2. Cleans up old installation (removes mobile-recon-cli, old aliases)"
    echo "  3. Builds the unified mobile-recon CLI"
    echo "  4. Installs it to \$GOPATH/bin or \$HOME/go/bin"
    echo "  5. Verifies PATH configuration"
    echo "  6. Optionally sets up helpful aliases"
    echo ""
    echo "After installation, all tools are available via mobile-recon:"
    echo "  mobile-recon adb   - Android Debug Bridge toolkit"
    echo "  mobile-recon nmap  - Network reconnaissance toolkit"
    echo "  mobile-recon apk   - Android APK static analysis"
    echo "  mobile-recon ios   - iOS Simulator toolkit"
    echo ""
    echo "Example commands:"
    echo "  mobile-recon list"
    echo "  mobile-recon adb device list"
    echo "  mobile-recon nmap scan quick 192.168.1.0/24"
    echo "  mobile-recon apk info app.apk"
    echo "  mobile-recon ios device list"
    echo ""
}

show_next_steps() {
    print_header "Next Steps"

    echo "1. Reload your shell or open a new terminal:"
    echo -e "   ${YELLOW}source ~/.zshrc${NC}   # for zsh"
    echo -e "   ${YELLOW}source ~/.bashrc${NC}  # for bash"
    echo ""
    echo "2. Verify installation:"
    echo -e "   ${YELLOW}mobile-recon --help${NC}"
    echo ""
    echo "3. Try some commands:"
    echo -e "   ${YELLOW}mobile-recon adb device list${NC}"
    echo -e "   ${YELLOW}mobile-recon nmap scan quick 192.168.1.0/24${NC}"
    echo -e "   ${YELLOW}mobile-recon apk info app.apk${NC}"
    echo -e "   ${YELLOW}mobile-recon ios device list${NC}"
    echo ""
    echo "4. Get help:"
    echo -e "   ${YELLOW}mobile-recon --help${NC}"
    echo -e "   ${YELLOW}mobile-recon adb --help${NC}"
    echo -e "   ${YELLOW}mobile-recon nmap --help${NC}"
    echo ""
    echo "Documentation:"
    echo "   Quick Start: $PROJECT_ROOT/QUICK_START.md"
    echo "   Full README: $PROJECT_ROOT/README.md"
    echo ""
    print_success "Installation complete! Happy hacking!"
}

#############################################################################
# Main
#############################################################################

main() {
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done

    # Welcome message
    echo ""
    echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║         Mobile Recon Toolkit - Installation Script         ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"

    # Run installation steps
    check_prerequisites
    cleanup_old_installation

    if ! install_mobile_recon; then
        print_error "Installation failed!"
        exit 1
    fi

    verify_installation
    setup_aliases
    show_next_steps

    echo ""
}

# Run main function
main "$@"
