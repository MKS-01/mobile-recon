#!/usr/bin/env bash

#############################################################################
# Mobile Recon Toolkit - New Tool Generator
#############################################################################
# This script generates a new tool from the template with proper naming
# and structure.
#
# Usage:
#   ./scripts/new-tool.sh <tool-name> <category> <description>
#
# Example:
#   ./scripts/new-tool.sh frida-toolkit Mobile "Frida dynamic instrumentation"
#############################################################################

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# Script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TEMPLATE_DIR="$PROJECT_ROOT/.templates/tool-template"
GO_TOOLS_DIR="$PROJECT_ROOT/go-tools"

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
    echo -e "${CYAN}ℹ${NC} $1"
}

print_step() {
    echo -e "${BLUE}▶${NC} $1"
}

show_usage() {
    cat << EOF
${BLUE}Mobile Recon Toolkit - New Tool Generator${NC}

Usage: $0 <tool-name> <category> <description>

Arguments:
  tool-name     Name of the tool (kebab-case, e.g., frida-toolkit)
  category      Tool category: Mobile, Network, Web, or Other
  description   Short description of the tool (quoted)

Examples:
  $0 frida-toolkit Mobile "Frida dynamic instrumentation toolkit"
  $0 burp-toolkit Web "Burp Suite integration for API testing"
  $0 wifi-toolkit Network "WiFi security analysis toolkit"

The script will:
  1. Create a new tool directory from the template
  2. Replace all template variables with your values
  3. Initialize Go modules
  4. Register the tool with mobile-recon-cli
  5. Provide next steps for implementation

EOF
}

validate_tool_name() {
    local name=$1

    # Check if name is kebab-case
    if ! [[ $name =~ ^[a-z][a-z0-9-]*$ ]]; then
        print_error "Tool name must be in kebab-case (lowercase with hyphens)"
        echo "  Examples: frida-toolkit, burp-integration, wifi-analyzer"
        return 1
    fi

    # Check if tool already exists
    if [ -d "$GO_TOOLS_DIR/$name" ]; then
        print_error "Tool '$name' already exists in $GO_TOOLS_DIR"
        return 1
    fi

    return 0
}

validate_category() {
    local category=$1

    case $category in
        Mobile|Network|Web|Other)
            return 0
            ;;
        *)
            print_error "Category must be one of: Mobile, Network, Web, Other"
            return 1
            ;;
    esac
}

to_title_case() {
    local input=$1
    # Convert kebab-case to Title Case (e.g., frida-toolkit -> Frida Toolkit)
    echo "$input" | sed 's/-/ /g' | sed 's/\b\(.\)/\u\1/g'
}

#############################################################################
# Template Processing
#############################################################################

process_template_file() {
    local src=$1
    local dst=$2
    local tool_name=$3
    local tool_name_title=$4
    local tool_short_desc=$5
    local tool_category=$6

    # Read file and replace variables
    sed -e "s/{{TOOL_NAME}}/$tool_name/g" \
        -e "s/{{TOOL_NAME_TITLE}}/$tool_name_title/g" \
        -e "s/{{TOOL_SHORT_DESCRIPTION}}/$tool_short_desc/g" \
        -e "s/{{TOOL_LONG_DESCRIPTION}}/$tool_short_desc/g" \
        -e "s/{{TOOL_CATEGORY}}/$tool_category/g" \
        -e "s/{{VERSION}}/1.0.0/g" \
        -e "s/{{USE_CASE}}/security testing and analysis/g" \
        -e "s/{{USE_CASE_1}}/Basic Analysis/g" \
        -e "s/{{USE_CASE_2}}/Advanced Testing/g" \
        "$src" > "$dst"
}

copy_template() {
    local tool_name=$1
    local tool_name_title=$2
    local tool_short_desc=$3
    local tool_category=$4
    local dest_dir="$GO_TOOLS_DIR/$tool_name"

    print_step "Creating tool directory structure..."

    # Create directory structure
    mkdir -p "$dest_dir/cmd"
    mkdir -p "$dest_dir/pkg/core"

    # Process and copy template files
    print_step "Copying template files..."

    # Main files
    process_template_file "$TEMPLATE_DIR/main.go" "$dest_dir/main.go" \
        "$tool_name" "$tool_name_title" "$tool_short_desc" "$tool_category"

    process_template_file "$TEMPLATE_DIR/go.mod.template" "$dest_dir/go.mod" \
        "$tool_name" "$tool_name_title" "$tool_short_desc" "$tool_category"

    process_template_file "$TEMPLATE_DIR/README.md" "$dest_dir/README.md" \
        "$tool_name" "$tool_name_title" "$tool_short_desc" "$tool_category"

    cp "$TEMPLATE_DIR/.gitignore" "$dest_dir/.gitignore"

    # CMD files
    process_template_file "$TEMPLATE_DIR/cmd/root.go" "$dest_dir/cmd/root.go" \
        "$tool_name" "$tool_name_title" "$tool_short_desc" "$tool_category"

    process_template_file "$TEMPLATE_DIR/cmd/example.go" "$dest_dir/cmd/example.go" \
        "$tool_name" "$tool_name_title" "$tool_short_desc" "$tool_category"

    # PKG files
    process_template_file "$TEMPLATE_DIR/pkg/example/example.go" "$dest_dir/pkg/core/core.go" \
        "$tool_name" "$tool_name_title" "$tool_short_desc" "$tool_category"

    print_success "Template files created"
}

initialize_go_module() {
    local tool_name=$1
    local dest_dir="$GO_TOOLS_DIR/$tool_name"

    print_step "Initializing Go module..."

    cd "$dest_dir"
    go mod tidy > /dev/null 2>&1

    print_success "Go module initialized"
}

register_with_cli() {
    local tool_name=$1
    local tool_name_title=$2
    local tool_short_desc=$3
    local tool_category=$4

    print_step "Registering tool with mobile-recon-cli..."

    local cli_manager="$GO_TOOLS_DIR/mobile-recon-cli/pkg/toolmanager/toolmanager.go"

    if [ ! -f "$cli_manager" ]; then
        print_warning "Could not find toolmanager.go - skipping automatic registration"
        print_info "You'll need to manually add the tool to mobile-recon-cli"
        return
    fi

    # Create backup
    cp "$cli_manager" "$cli_manager.bak"

    # Find the knownTools array and add new entry
    local new_entry="		{
			name:        \"$tool_name\",
			displayName: \"$tool_name_title\",
			dir:         \"$tool_name\",
			binary:      \"$tool_name\",
			description: \"$tool_short_desc\",
			category:    \"$tool_category\",
		},"

    # Insert before the closing brace of knownTools
    awk -v entry="$new_entry" '
        /^	}$/ && !done {
            print entry
            done=1
        }
        { print }
    ' "$cli_manager.bak" > "$cli_manager"

    rm "$cli_manager.bak"

    print_success "Tool registered with mobile-recon-cli"
}

show_next_steps() {
    local tool_name=$1
    local dest_dir="$GO_TOOLS_DIR/$tool_name"

    print_header "Tool Created Successfully!"

    echo "📁 Location: $dest_dir"
    echo ""
    echo "📝 Next Steps:"
    echo ""
    echo "1. Navigate to the tool directory:"
    echo -e "   ${YELLOW}cd go-tools/$tool_name${NC}"
    echo ""
    echo "2. Review and customize the files:"
    echo -e "   ${YELLOW}# Edit the main command logic${NC}"
    echo "   vim cmd/root.go"
    echo ""
    echo -e "   ${YELLOW}# Implement your core functionality${NC}"
    echo "   vim pkg/core/core.go"
    echo ""
    echo -e "   ${YELLOW}# Add more commands as needed${NC}"
    echo "   cp cmd/example.go cmd/your-command.go"
    echo ""
    echo "3. Build and test:"
    echo -e "   ${YELLOW}go build -o $tool_name${NC}"
    echo -e "   ${YELLOW}./$tool_name --help${NC}"
    echo ""
    echo "4. Install globally:"
    echo -e "   ${YELLOW}go install${NC}"
    echo ""
    echo "5. Update documentation:"
    echo "   • Edit README.md with actual features and usage"
    echo "   • Add examples and use cases"
    echo "   • Document all commands and flags"
    echo ""
    echo "6. Test with mobile-recon-cli:"
    echo -e "   ${YELLOW}mobile-recon-cli build $tool_name${NC}"
    echo -e "   ${YELLOW}mobile-recon-cli run $tool_name${NC}"
    echo ""
    echo "📚 Resources:"
    echo "   • Template structure: .templates/tool-template/"
    echo "   • Existing tools: go-tools/adb-toolkit/, go-tools/nmap-toolkit/"
    echo "   • Cobra docs: https://github.com/spf13/cobra"
    echo ""
    print_success "Happy coding! 🚀"
    echo ""
}

#############################################################################
# Main
#############################################################################

main() {
    # Check arguments
    if [ $# -lt 3 ]; then
        show_usage
        exit 1
    fi

    if [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
        show_usage
        exit 0
    fi

    local tool_name=$1
    local tool_category=$2
    local tool_short_desc=$3

    # Welcome message
    echo ""
    echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║         Mobile Recon - New Tool Generator                 ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"

    # Validate inputs
    print_header "Validating Input"

    if ! validate_tool_name "$tool_name"; then
        exit 1
    fi
    print_success "Tool name is valid"

    if ! validate_category "$tool_category"; then
        exit 1
    fi
    print_success "Category is valid"

    # Generate title case name
    tool_name_title=$(to_title_case "$tool_name")

    # Confirm details
    echo ""
    print_info "Tool Details:"
    echo "  Name:        $tool_name"
    echo "  Title:       $tool_name_title"
    echo "  Category:    $tool_category"
    echo "  Description: $tool_short_desc"
    echo ""

    read -p "Proceed with tool creation? (y/N): " -n 1 -r
    echo ""

    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_warning "Tool creation cancelled"
        exit 0
    fi

    # Create tool
    print_header "Creating Tool"

    copy_template "$tool_name" "$tool_name_title" "$tool_short_desc" "$tool_category"
    initialize_go_module "$tool_name"
    register_with_cli "$tool_name" "$tool_name_title" "$tool_short_desc" "$tool_category"

    # Show next steps
    show_next_steps "$tool_name"
}

# Run main function
main "$@"
