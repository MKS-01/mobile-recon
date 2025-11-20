// Package main is the entry point for the Mobile Recon unified CLI.
// It provides a unified interface to access all mobile reconnaissance tools.
package main

import (
	"github.com/MKS-01/mobile-recon/go-tools/mobile-recon-cli/cmd"
)

func main() {
	cmd.Execute()
}
