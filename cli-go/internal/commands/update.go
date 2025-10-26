package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/drasdp/codekiwi-cli/internal/config"
	"github.com/drasdp/codekiwi-cli/internal/docker"
	"github.com/drasdp/codekiwi-cli/internal/platform"
	"github.com/spf13/cobra"
)

var (
	skipImage bool
	skipCLI   bool
)

// UpdateCmd updates CodeKiwi CLI and Docker image
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update CodeKiwi CLI and Docker image",
	Long:  `Update the CodeKiwi CLI to the latest version and pull the latest Docker image.`,
	RunE:  runUpdate,
}

func init() {
	UpdateCmd.Flags().BoolVar(&skipImage, "skip-image", false, "Skip Docker image update")
	UpdateCmd.Flags().BoolVar(&skipCLI, "skip-cli", false, "Skip CLI update")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	platform.PrintInfo("Checking for updates...")

	cfg := config.Get()
	updated := false

	// Update CLI
	if !skipCLI {
		platform.PrintInfo("Updating CodeKiwi CLI...")
		if err := updateCLI(); err != nil {
			platform.PrintWarning(fmt.Sprintf("Failed to update CLI: %v", err))
		} else {
			updated = true
			platform.PrintSuccess("CLI updated successfully")
		}
	}

	// Update Docker image
	if !skipImage && !cfg.IsDevelopment {
		platform.PrintInfo("Updating Docker image...")
		if err := docker.PullImage(cfg.GetFullImageName()); err != nil {
			platform.PrintWarning(fmt.Sprintf("Failed to update Docker image: %v", err))
		} else {
			updated = true
			platform.PrintSuccess("Docker image updated successfully")
		}
	}

	if updated {
		platform.PrintSuccess("Update completed!")
		platform.PrintInfo("Restart any running instances to use the new version")
	} else {
		platform.PrintInfo("Nothing to update")
	}

	return nil
}

func updateCLI() error {
	// Determine the install script URL based on OS
	var scriptURL string
	var installCmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		scriptURL = "https://raw.githubusercontent.com/drasdp/codekiwi-cli/main/install/install.bat"
		// Download and run batch file
		installCmd = exec.Command("cmd", "/c",
			fmt.Sprintf("curl -fsSL %s -o %%TEMP%%\\install.bat && %%TEMP%%\\install.bat", scriptURL))
	case "darwin", "linux":
		scriptURL = "https://raw.githubusercontent.com/drasdp/codekiwi-cli/main/install/install.sh"
		// Run install script directly with bash
		installCmd = exec.Command("bash", "-c",
			fmt.Sprintf("curl -fsSL %s | bash", scriptURL))
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Set up environment
	installCmd.Env = os.Environ()
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	// Run the install script
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("installation script failed: %w", err)
	}

	return nil
}