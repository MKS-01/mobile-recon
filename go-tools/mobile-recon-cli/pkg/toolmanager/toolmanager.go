// Package toolmanager handles tool discovery and execution.
package toolmanager

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// Version is set via ldflags at build time
var Version = "dev"

// SourcePath is set via ldflags at build time to the go-tools directory
var SourcePath = ""

//go:embed tools.yaml
var embeddedConfig embed.FS

type Tool struct {
	Name        string `yaml:"name"`
	DisplayName string `yaml:"display_name"`
	Path        string `yaml:"-"`
	Binary      string `yaml:"binary"`
	Description string `yaml:"description"`
	Category    string `yaml:"category"`
	Available   bool   `yaml:"-"`
}

type ToolCategory struct {
	Name        string
	DisplayName string
	Tools       []Tool
}

type ToolManager struct {
	RootPath   string
	Categories []ToolCategory
	Verbose    bool
}

// BuildOptions configures the build process
type BuildOptions struct {
	Verbose  bool
	Version  string
	Parallel bool
}

type toolConfig struct {
	Tools []struct {
		Name        string `yaml:"name"`
		DisplayName string `yaml:"display_name"`
		Dir         string `yaml:"dir"`
		Binary      string `yaml:"binary"`
		Description string `yaml:"description"`
		Category    string `yaml:"category"`
	} `yaml:"tools"`
}

func NewToolManager() (*ToolManager, error) {
	var rootPath string

	// Use embedded source path if set (from ldflags during build)
	if SourcePath != "" {
		rootPath = SourcePath
	} else {
		// Fall back to executable-relative path for development
		execPath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("failed to get executable path: %v", err)
		}
		rootPath = filepath.Dir(filepath.Dir(execPath))
	}

	tm := &ToolManager{
		RootPath:   rootPath,
		Categories: []ToolCategory{},
	}

	if err := tm.DiscoverTools(); err != nil {
		return nil, err
	}

	return tm, nil
}

func New(baseDir string) *ToolManager {
	tm := &ToolManager{
		RootPath:   baseDir,
		Categories: []ToolCategory{},
	}
	tm.DiscoverTools()
	return tm
}

func (tm *ToolManager) loadConfig() (*toolConfig, error) {
	// Try embedded config first
	data, err := embeddedConfig.ReadFile("tools.yaml")
	if err != nil {
		// Fall back to file in root path
		configPath := filepath.Join(tm.RootPath, "tools.yaml")
		data, err = os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load tools config: %v", err)
		}
	}

	var config toolConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse tools config: %v", err)
	}

	return &config, nil
}

func (tm *ToolManager) DiscoverTools() error {
	config, err := tm.loadConfig()
	if err != nil {
		return err
	}

	categories := make(map[string][]Tool)

	for _, kt := range config.Tools {
		toolPath := filepath.Join(tm.RootPath, kt.Dir)
		binaryPath := filepath.Join(toolPath, kt.Binary)

		tool := Tool{
			Name:        kt.Name,
			DisplayName: kt.DisplayName,
			Path:        toolPath,
			Binary:      binaryPath,
			Description: kt.Description,
			Category:    kt.Category,
			Available:   false,
		}

		if _, err := os.Stat(binaryPath); err == nil {
			tool.Available = true
		}

		categories[kt.Category] = append(categories[kt.Category], tool)
	}

	for categoryName, tools := range categories {
		tm.Categories = append(tm.Categories, ToolCategory{
			Name:        strings.ToLower(categoryName),
			DisplayName: categoryName,
			Tools:       tools,
		})
	}

	return nil
}

func (tm *ToolManager) GetTool(name string) (*Tool, error) {
	for _, category := range tm.Categories {
		for i := range category.Tools {
			tool := &category.Tools[i]
			if tool.Name == name || strings.EqualFold(tool.DisplayName, name) {
				return tool, nil
			}
		}
	}
	return nil, fmt.Errorf("tool not found: %s", name)
}

func (tm *ToolManager) ListTools() []Tool {
	var tools []Tool
	for _, category := range tm.Categories {
		tools = append(tools, category.Tools...)
	}
	return tools
}

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

func (tm *ToolManager) BuildTool(toolName string) error {
	return tm.BuildToolWithOptions(toolName, BuildOptions{Verbose: tm.Verbose})
}

func (tm *ToolManager) BuildToolWithOptions(toolName string, opts BuildOptions) error {
	tool, err := tm.GetTool(toolName)
	if err != nil {
		return err
	}

	version := opts.Version
	if version == "" {
		version = Version
	}

	ldflags := fmt.Sprintf("-s -w -X main.Version=%s", version)

	args := []string{"build", "-trimpath", "-ldflags", ldflags, "-o", tool.Name}
	cmd := exec.Command("go", args...)
	cmd.Dir = tool.Path

	if opts.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build %s: %v", tool.DisplayName, err)
	}

	tool.Available = true
	return nil
}

func (tm *ToolManager) BuildAllTools() error {
	return tm.BuildAllToolsWithOptions(BuildOptions{Verbose: tm.Verbose})
}

func (tm *ToolManager) BuildAllToolsWithOptions(opts BuildOptions) error {
	tools := tm.ListTools()

	if opts.Parallel && len(tools) > 1 {
		return tm.buildToolsParallel(tools, opts)
	}

	for _, tool := range tools {
		fmt.Printf("Building %s...\n", tool.DisplayName)
		if err := tm.BuildToolWithOptions(tool.Name, opts); err != nil {
			return err
		}
	}
	return nil
}

func (tm *ToolManager) buildToolsParallel(tools []Tool, opts BuildOptions) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(tools))

	for _, tool := range tools {
		wg.Add(1)
		go func(t Tool) {
			defer wg.Done()
			fmt.Printf("Building %s...\n", t.DisplayName)
			if err := tm.BuildToolWithOptions(t.Name, opts); err != nil {
				errChan <- err
			}
		}(tool)
	}

	wg.Wait()
	close(errChan)

	var errs []string
	for err := range errChan {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return fmt.Errorf("build errors:\n  %s", strings.Join(errs, "\n  "))
	}

	return nil
}

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
