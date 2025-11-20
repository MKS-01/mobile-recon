#!/usr/bin/env bash

#############################################################################
# Mobile Recon Toolkit - Installation Test Script
#############################################################################
# This script tests that all tools are properly installed and accessible.
#############################################################################

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║      Mobile Recon Toolkit - Installation Test             ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Ensure PATH includes Go bin
export PATH="$HOME/go/bin:$PATH"

# Test each tool
echo -e "${BLUE}Testing installed tools...${NC}"
echo ""

success_count=0
fail_count=0

# Test mobile-recon-cli
echo -n "Testing mobile-recon-cli... "
if command -v mobile-recon-cli &> /dev/null; then
    echo -e "${GREEN}✓ Found${NC}"
    ((success_count++))
else
    echo -e "${RED}✗ Not found${NC}"
    ((fail_count++))
fi

# Test adb-toolkit
echo -n "Testing adb-toolkit...       "
if command -v adb-toolkit &> /dev/null; then
    echo -e "${GREEN}✓ Found${NC}"
    ((success_count++))
else
    echo -e "${RED}✗ Not found${NC}"
    ((fail_count++))
fi

# Test nmap-toolkit
echo -n "Testing nmap-toolkit...      "
if command -v nmap-toolkit &> /dev/null; then
    echo -e "${GREEN}✓ Found${NC}"
    ((success_count++))
else
    echo -e "${RED}✗ Not found${NC}"
    ((fail_count++))
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

if [ $fail_count -eq 0 ]; then
    echo -e "${GREEN}✓ All tools installed successfully!${NC}"
    echo ""
    echo "Try these commands:"
    echo -e "  ${YELLOW}mobile-recon-cli list${NC}"
    echo -e "  ${YELLOW}adb-toolkit device list${NC}"
    echo -e "  ${YELLOW}nmap-toolkit scan quick 192.168.1.0/24${NC}"
else
    echo -e "${RED}✗ Some tools are not installed or not in PATH${NC}"
    echo ""
    echo "To fix this:"
    echo "1. Run: ./scripts/install.sh"
    echo "2. Reload shell: source ~/.zshrc"
fi

echo ""
