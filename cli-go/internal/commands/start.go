package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/drasdp/codekiwi-cli/internal/config"
	"github.com/drasdp/codekiwi-cli/internal/docker"
	"github.com/drasdp/codekiwi-cli/internal/platform"
	"github.com/drasdp/codekiwi-cli/internal/state"
	"github.com/spf13/cobra"
)

var (
	webPort    int
	devPort    int
	noOpen     bool
	followLogs bool
	template   string
)

// StartCmd starts a CodeKiwi project
var StartCmd = &cobra.Command{
	Use:   "start [path]",
	Short: "Start CodeKiwi for a project",
	Long:  `Start a CodeKiwi development environment for the specified project directory.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runStart,
}

func init() {
	StartCmd.Flags().IntVarP(&webPort, "web-port", "p", 0, "Web port (default: auto)")
	StartCmd.Flags().IntVar(&devPort, "dev-port", 0, "Dev server port (default: auto)")
	StartCmd.Flags().BoolVarP(&noOpen, "no-open", "n", false, "Don't open browser automatically")
	StartCmd.Flags().BoolVarP(&followLogs, "follow", "f", false, "Follow container logs")
	StartCmd.Flags().StringVarP(&template, "template", "t", "", "Project template to use")
}

func runStart(cmd *cobra.Command, args []string) error {
	// Get project path
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}

	// Resolve absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}
	projectPath = absPath

	// Normalize path for Windows
	projectPath = config.NormalizePath(projectPath)

	// Get project name
	projectName := filepath.Base(projectPath)
	if projectName == "." || projectName == "/" || projectName == "\\" {
		// If root directory, use current directory name
		cwd, _ := os.Getwd()
		projectName = filepath.Base(cwd)
	}

	platform.PrintStarting(projectName)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check Docker
	platform.PrintInfo("Checking Docker...")
	if err := docker.CheckDocker(); err != nil {
		return err
	}
	platform.PrintSuccess("Docker is ready")

	// Check if project is already running
	instance, err := state.LoadByPath(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	if instance != nil && docker.IsContainerRunning(instance.ContainerName) {
		platform.PrintWarning(fmt.Sprintf("Project already running at %s", projectPath))
		platform.PrintURL("Web URL", instance.GetURL())
		platform.PrintInfo(fmt.Sprintf("Container: %s", instance.ContainerName))
		platform.PrintInfo(fmt.Sprintf("Uptime: %s", instance.GetUptime()))

		if !noOpen && platform.CanOpenBrowser() {
			platform.OpenBrowser(instance.GetURL())
		}

		if followLogs {
			platform.PrintInfo("Following container logs...")
			return docker.FollowContainerLogs(instance.ContainerName)
		}

		return nil
	}

	// Generate container name
	projectHash := state.GenerateHash(projectPath)
	containerName := fmt.Sprintf("%s-%s", cfg.ContainerPrefix, projectHash)

	// Find available ports
	if webPort == 0 {
		webPort, err = platform.FindAvailablePort(cfg.WebPortDefault, 100)
		if err != nil {
			return fmt.Errorf("failed to find available web port: %w", err)
		}
	}

	if devPort == 0 {
		devPort, err = platform.FindAvailablePort(cfg.DevPortDefault, 100)
		if err != nil {
			return fmt.Errorf("failed to find available dev port: %w", err)
		}
	}

	platform.PrintInfo(fmt.Sprintf("Using ports - Web: %d, Dev: %d", webPort, devPort))

	// Create project directory if it doesn't exist
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		platform.PrintInfo(fmt.Sprintf("Creating directory: %s", projectPath))
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Set template environment variable if specified
	if template != "" {
		os.Setenv("INSTALL_TEMPLATE", template)
		platform.PrintInfo(fmt.Sprintf("Using template: %s", template))
	} else {
		// Check if directory is empty
		entries, err := os.ReadDir(projectPath)
		if err == nil && len(entries) == 0 {
			os.Setenv("INSTALL_TEMPLATE", "yes")
			platform.PrintInfo("Empty directory detected, will install default template")
		} else {
			os.Setenv("INSTALL_TEMPLATE", "no")
		}
	}

	// Pull image if not in development mode
	if !cfg.IsDevelopment {
		if err := docker.PullImage(cfg.GetFullImageName()); err != nil {
			platform.PrintWarning("Failed to pull latest image, using local if available")
		}
	}

	// Start container
	if err := docker.ComposeUp(projectPath, containerName, webPort, devPort); err != nil {
		return err
	}

	// Save instance state
	instance = state.Create(projectPath, containerName, webPort, devPort)
	if err := instance.Save(); err != nil {
		platform.PrintWarning("Failed to save instance state")
	}

	// Print success message
	platform.PrintSuccess("CodeKiwi started successfully!")
	platform.PrintDirectory("Project", projectPath)
	platform.PrintURL("Web URL", instance.GetURL())
	platform.PrintInfo(fmt.Sprintf("Container: %s", containerName))

	// Wait a bit for services to fully start
	platform.PrintInfo("Waiting for services to start...")
	time.Sleep(3 * time.Second)

	// Open browser
	if !noOpen && platform.CanOpenBrowser() {
		platform.PrintInfo("Opening browser...")
		if err := platform.OpenBrowser(instance.GetURL()); err != nil {
			platform.PrintWarning(fmt.Sprintf("Failed to open browser: %v", err))
		}
	}

	// Print instructions
	fmt.Println()
	fmt.Println("To stop CodeKiwi, run:")
	fmt.Printf("  codekiwi kill %s\n", projectPath)
	fmt.Println()

	// Follow logs if requested
	if followLogs {
		platform.PrintInfo("Following container logs...")
		return docker.FollowContainerLogs(containerName)
	}

	return nil
}

// Helper function to check if directory is empty
func isEmptyDir(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	// Ignore hidden files like .git
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			return false
		}
	}

	return true
}