#!/usr/bin/env bash

#############################################################################
# Mobile Recon Toolkit - Installation Script
#############################################################################
# This script builds and installs all Go-based CLI tools globally.
#
# Usage:
#   ./scripts/install.sh              # Install all tools
#   ./scripts/install.sh --verbose    # Install with verbose output
#   ./scripts/install.sh --help       # Show help
#
# What it does:
#   1. Checks prerequisites (Go installation)
#   2. Builds all Go CLI tools
#   3. Installs them to $GOPATH/bin or $HOME/go/bin
#   4. Verifies PATH configuration
#   5. Creates helpful aliases (optional)
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
        echo ""
        echo "  Add this to your shell configuration (~/.bashrc or ~/.zshrc):"
        echo -e "  ${YELLOW}export PATH=\"\$HOME/go/bin:\$PATH\"${NC}"
        echo ""
    else
        print_success "GOPATH/bin is in PATH"
    fi
}

#############################################################################
# Build and Install Tools
#############################################################################

install_tool() {
    local tool_name=$1
    local tool_dir=$2
    local tool_path="$GO_TOOLS_DIR/$tool_dir"

    if [ ! -d "$tool_path" ]; then
        print_error "Tool directory not found: $tool_path"
        return 1
    fi

    print_step "Building $tool_name..."

    cd "$tool_path"

    # Download dependencies
    if [ "$VERBOSE" = true ]; then
        go mod download
    else
        go mod download > /dev/null 2>&1
    fi

    # Build locally first
    if [ "$VERBOSE" = true ]; then
        go build -o "$tool_dir"
    else
        go build -o "$tool_dir" > /dev/null 2>&1
    fi

    if [ $? -eq 0 ]; then
        print_success "Built $tool_name successfully"
    else
        print_error "Failed to build $tool_name"
        return 1
    fi

    # Install globally
    print_step "Installing $tool_name globally..."

    if [ "$VERBOSE" = true ]; then
        go install
    else
        go install > /dev/null 2>&1
    fi

    if [ $? -eq 0 ]; then
        print_success "Installed $tool_name to $GOPATH/bin/$tool_dir"
    else
        print_error "Failed to install $tool_name"
        return 1
    fi

    echo ""
}

install_all_tools() {
    print_header "Building and Installing Tools"

    # Array of tools: "Display Name:directory-name"
    local tools=(
        "Mobile Recon CLI:mobile-recon-cli"
        "ADB Toolkit:adb-toolkit"
        "Nmap Toolkit:nmap-toolkit"
    )

    local success_count=0
    local fail_count=0

    for tool in "${tools[@]}"; do
        IFS=':' read -r name dir <<< "$tool"

        if install_tool "$name" "$dir"; then
            ((success_count++))
        else
            ((fail_count++))
        fi
    done

    echo ""
    print_header "Installation Summary"

    if [ $success_count -gt 0 ]; then
        print_success "$success_count tool(s) installed successfully"
    fi

    if [ $fail_count -gt 0 ]; then
        print_error "$fail_count tool(s) failed to install"
        return 1
    fi

    return 0
}

#############################################################################
# Verify Installation
#############################################################################

verify_installation() {
    print_header "Verifying Installation"

    local tools=("mobile-recon-cli" "adb-toolkit" "nmap-toolkit")
    local all_found=true

    for tool in "${tools[@]}"; do
        if command -v "$tool" &> /dev/null; then
            local tool_path=$(which "$tool")
            print_success "$tool found at: $tool_path"
        else
            print_warning "$tool not found in PATH"
            echo "  Expected location: $GOPATH/bin/$tool"
            all_found=false
        fi
    done

    echo ""

    if [ "$all_found" = false ]; then
        print_warning "Some tools are not in your PATH"
        echo ""
        echo "  Make sure to add this to your shell configuration:"
        echo -e "  ${YELLOW}export PATH=\"\$HOME/go/bin:\$PATH\"${NC}"
        echo ""
        echo "  Then reload your shell:"
        echo -e "  ${YELLOW}source ~/.bashrc${NC}  # for bash"
        echo -e "  ${YELLOW}source ~/.zshrc${NC}   # for zsh"
        echo ""
    fi
}

#############################################################################
# Setup Aliases
#############################################################################

setup_aliases() {
    print_header "Setting Up Aliases (Optional)"

    echo "Would you like to add helpful aliases to your shell configuration?"
    echo ""
    echo "Suggested aliases:"
    echo "  alias mobile-recon='mobile-recon-cli'"
    echo "  alias mr='mobile-recon-cli'"
    echo "  alias mradb='mobile-recon-cli adb'"
    echo "  alias mrnmap='mobile-recon-cli nmap'"
    echo ""

    read -p "Add aliases? (y/N): " -n 1 -r
    echo ""

    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # Detect shell
        local shell_config=""
        if [ -n "$ZSH_VERSION" ]; then
            shell_config="$HOME/.zshrc"
        elif [ -n "$BASH_VERSION" ]; then
            shell_config="$HOME/.bashrc"
        else
            echo "Could not detect shell. Please add aliases manually."
            return
        fi

        echo "" >> "$shell_config"
        echo "# Mobile Recon Toolkit Aliases (added by install script)" >> "$shell_config"
        echo "alias mobile-recon='mobile-recon-cli'" >> "$shell_config"
        echo "alias mr='mobile-recon-cli'" >> "$shell_config"
        echo "alias mradb='mobile-recon-cli adb'" >> "$shell_config"
        echo "alias mrnmap='mobile-recon-cli nmap'" >> "$shell_config"

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
    echo "  2. Builds all Go CLI tools"
    echo "  3. Installs them to \$GOPATH/bin or \$HOME/go/bin"
    echo "  4. Verifies PATH configuration"
    echo "  5. Optionally sets up helpful aliases"
    echo ""
    echo "Tools that will be installed:"
    echo "  • mobile-recon-cli  - Unified CLI manager"
    echo "  • adb-toolkit       - Android Debug Bridge toolkit"
    echo "  • nmap-toolkit      - Network reconnaissance toolkit"
    echo ""
    echo "After installation, you can use:"
    echo "  mobile-recon-cli list"
    echo "  mobile-recon-cli adb device list"
    echo "  adb-toolkit device list"
    echo "  nmap-toolkit scan quick 192.168.1.0/24"
    echo ""
}

show_next_steps() {
    print_header "Next Steps"

    echo "1. Reload your shell or open a new terminal:"
    echo -e "   ${YELLOW}source ~/.zshrc${NC}   # for zsh"
    echo -e "   ${YELLOW}source ~/.bashrc${NC}  # for bash"
    echo ""
    echo "2. Verify installation:"
    echo -e "   ${YELLOW}mobile-recon-cli list${NC}"
    echo ""
    echo "3. Try some commands:"
    echo -e "   ${YELLOW}adb-toolkit device list${NC}"
    echo -e "   ${YELLOW}nmap-toolkit scan quick 192.168.1.0/24${NC}"
    echo -e "   ${YELLOW}mobile-recon-cli interactive${NC}"
    echo ""
    echo "4. Get help:"
    echo -e "   ${YELLOW}mobile-recon-cli --help${NC}"
    echo -e "   ${YELLOW}adb-toolkit --help${NC}"
    echo ""
    echo "📚 Documentation:"
    echo "   • Quick Start: $PROJECT_ROOT/QUICK_START.md"
    echo "   • Full README: $PROJECT_ROOT/README.md"
    echo ""
    print_success "Installation complete! Happy hacking! 🚀"
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
    echo -e "${BLUE}║         Mobile Recon Toolkit - Installation Script        ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"

    # Run installation steps
    check_prerequisites

    if ! install_all_tools; then
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
