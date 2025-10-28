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
	// Use the exact same install commands as in README.md
	var installCmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Use the exact command from README.md for Windows
		installCmd = exec.Command("cmd", "/c",
			"curl -fsSL https://raw.githubusercontent.com/drasdp/codekiwi-cli/main/install/install.bat -o %TEMP%\\codekiwi-install.bat && %TEMP%\\codekiwi-install.bat")
	case "darwin", "linux":
		// Use the exact command from README.md for Unix systems
		installCmd = exec.Command("bash", "-c",
			"curl -fsSL https://raw.githubusercontent.com/drasdp/codekiwi-cli/main/install/install.sh | bash")
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