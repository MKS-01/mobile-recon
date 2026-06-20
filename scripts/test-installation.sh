#!/usr/bin/env bash

#############################################################################
# Mobile Recon - Installation Test Script
#############################################################################
# Verifies that the unified `mobile-recon` binary is installed, on PATH, and
# that its toolkit subcommands respond.
#############################################################################

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo ""
echo -e "${BLUE}Mobile Recon - Installation Test${NC}"
echo ""

# Ensure PATH includes Go bin
export PATH="$HOME/go/bin:$PATH"

fail_count=0

# 1. The binary must be on PATH
echo -n "Testing mobile-recon... "
if command -v mobile-recon &> /dev/null; then
    echo -e "${GREEN}✓ Found${NC}"
else
    echo -e "${RED}✗ Not found${NC}"
    fail_count=$((fail_count + 1))
fi

# 2. Each toolkit subcommand should respond to --help
if command -v mobile-recon &> /dev/null; then
    for tool in adb apk ios nmap; do
        echo -n "Testing 'mobile-recon $tool'... "
        if mobile-recon "$tool" --help &> /dev/null; then
            echo -e "${GREEN}✓ OK${NC}"
        else
            echo -e "${RED}✗ Failed${NC}"
            fail_count=$((fail_count + 1))
        fi
    done
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

if [ $fail_count -eq 0 ]; then
    echo -e "${GREEN}✓ mobile-recon is installed and all toolkits respond.${NC}"
    echo ""
    echo "Try these commands:"
    echo -e "  ${YELLOW}mobile-recon list${NC}"
    echo -e "  ${YELLOW}mobile-recon adb device list${NC}"
    echo -e "  ${YELLOW}mobile-recon nmap scan quick 192.168.1.0/24${NC}"
else
    echo -e "${RED}✗ mobile-recon is not installed or not on PATH${NC}"
    echo ""
    echo "To fix this:"
    echo "1. Run: ./scripts/install.sh"
    echo "2. Reload shell: source ~/.zshrc"
fi

echo ""
