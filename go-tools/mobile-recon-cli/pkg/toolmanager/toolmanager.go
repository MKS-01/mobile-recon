// Package toolmanager handles tool discovery and execution
package toolmanager

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Tool represents a reconnaissance tool
type Tool struct {
	Name        string
	DisplayName string
	Path        string
	Binary      string
	Description string
	Category    string
	Available   bool
}

// ToolCategory represents a category of tools
type ToolCategory struct {
	Name        string
	DisplayName string
	Tools       []Tool
}

// ToolManager manages available reconnaissance tools
type ToolManager struct {
	RootPath   string
	Categories []ToolCategory
}

// NewToolManager creates a new tool manager
func NewToolManager() (*ToolManager, error) {
	// Get the root path (go-tools directory)
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %v", err)
	}

	// Navigate up to go-tools directory
	rootPath := filepath.Dir(filepath.Dir(execPath))

	tm := &ToolManager{
		RootPath:   rootPath,
		Categories: []ToolCategory{},
	}

	if err := tm.DiscoverTools(); err != nil {
		return nil, err
	}

	return tm, nil
}

// DiscoverTools discovers all available tools in the go-tools directory
func (tm *ToolManager) DiscoverTools() error {
	categories := make(map[string][]Tool)

	// Define known tools
	knownTools := []struct {
		name        string
		displayName string
		dir         string
		binary      string
		description string
		category    string
	}{
		{
			name:        "adb-toolkit",
			displayName: "ADB Toolkit",
			dir:         "adb-toolkit",
			binary:      "adb-toolkit",
			description: "Android Debug Bridge automation toolkit",
			category:    "Mobile",
		},
		{
			name:        "nmap-toolkit",
			displayName: "Nmap Toolkit",
			dir:         "nmap-toolkit",
			binary:      "nmap-toolkit",
			description: "Network reconnaissance and scanning toolkit",
			category:    "Network",
		},
	}

	for _, kt := range knownTools {
		toolPath := filepath.Join(tm.RootPath, kt.dir)
		binaryPath := filepath.Join(toolPath, kt.binary)

		tool := Tool{
			Name:        kt.name,
			DisplayName: kt.displayName,
			Path:        toolPath,
			Binary:      binaryPath,
			Description: kt.description,
			Category:    kt.category,
			Available:   false,
		}

		// Check if binary exists
		if _, err := os.Stat(binaryPath); err == nil {
			tool.Available = true
		}

		categories[kt.category] = append(categories[kt.category], tool)
	}

	// Convert map to slice of categories
	for categoryName, tools := range categories {
		tm.Categories = append(tm.Categories, ToolCategory{
			Name:        strings.ToLower(categoryName),
			DisplayName: categoryName,
			Tools:       tools,
		})
	}

	return nil
}

// GetTool returns a tool by name
func (tm *ToolManager) GetTool(name string) (*Tool, error) {
	for _, category := range tm.Categories {
		for _, tool := range category.Tools {
			if tool.Name == name || strings.ToLower(tool.DisplayName) == strings.ToLower(name) {
				return &tool, nil
			}
		}
	}
	return nil, fmt.Errorf("tool not found: %s", name)
}

// GetCategory returns a category by name
func (tm *ToolManager) GetCategory(name string) (*ToolCategory, error) {
	for _, category := range tm.Categories {
		if category.Name == strings.ToLower(name) || strings.ToLower(category.DisplayName) == strings.ToLower(name) {
			return &category, nil
		}
	}
	return nil, fmt.Errorf("category not found: %s", name)
}

// ListTools returns all available tools
func (tm *ToolManager) ListTools() []Tool {
	var tools []Tool
	for _, category := range tm.Categories {
		tools = append(tools, category.Tools...)
	}
	return tools
}

// ListAvailableTools returns only built and available tools
func (tm *ToolManager) ListAvailableTools() []Tool {
	var tools []Tool
	for _, category := range tm.Categories {
		for _, tool := range category.Tools {
			if tool.Available {
				tools = append(tools, tool)
			}
		}
	}
	return tools
}

// RunTool executes a tool with the given arguments
func (tm *ToolManager) RunTool(toolName string, args []string) error {
	tool, err := tm.GetTool(toolName)
	if err != nil {
		return err
	}

	if !tool.Available {
		return fmt.Errorf("tool not built: %s (run 'mobile-recon build %s' first)", tool.DisplayName, tool.Name)
	}

	cmd := exec.Command(tool.Binary, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// BuildTool builds a specific tool
func (tm *ToolManager) BuildTool(toolName string) error {
	tool, err := tm.GetTool(toolName)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "build", "-o", tool.Name)
	cmd.Dir = tool.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build %s: %v", tool.DisplayName, err)
	}

	tool.Available = true
	return nil
}

// BuildAllTools builds all tools
func (tm *ToolManager) BuildAllTools() error {
	for _, category := range tm.Categories {
		for _, tool := range category.Tools {
			fmt.Printf("Building %s...\n", tool.DisplayName)
			if err := tm.BuildTool(tool.Name); err != nil {
				return err
			}
		}
	}
	return nil
}

// InstallTool installs a tool globally
func (tm *ToolManager) InstallTool(toolName string) error {
	tool, err := tm.GetTool(toolName)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "install")
	cmd.Dir = tool.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install %s: %v", tool.DisplayName, err)
	}

	return nil
}
