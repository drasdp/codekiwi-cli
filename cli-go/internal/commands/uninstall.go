package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/drasdp/codekiwi-cli/internal/config"
	"github.com/drasdp/codekiwi-cli/internal/docker"
	"github.com/drasdp/codekiwi-cli/internal/platform"
	"github.com/drasdp/codekiwi-cli/internal/state"
	"github.com/spf13/cobra"
)

var (
	keepData   bool
	keepImages bool
	forceUninstall bool
)

// UninstallCmd uninstalls CodeKiwi
var UninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall CodeKiwi",
	Long:  `Completely remove CodeKiwi, including all containers, images, and configuration.`,
	RunE:  runUninstall,
}

func init() {
	UninstallCmd.Flags().BoolVar(&keepData, "keep-data", false, "Keep configuration and state files")
	UninstallCmd.Flags().BoolVar(&keepImages, "keep-images", false, "Keep Docker images")
	UninstallCmd.Flags().BoolVarP(&forceUninstall, "force", "f", false, "Force uninstall without confirmation")
}

func runUninstall(cmd *cobra.Command, args []string) error {
	cfg := config.Get()

	platform.PrintWarning("This will uninstall CodeKiwi completely")

	if !forceUninstall {
		fmt.Println("\nThis will:")
		fmt.Println("  1. Stop all running CodeKiwi containers")
		if !keepImages {
			fmt.Println("  2. Remove CodeKiwi Docker images")
		}
		if !keepData {
			fmt.Printf("  3. Remove configuration from %s\n", cfg.InstallDir)
		}
		fmt.Println("  4. Remove the codekiwi command")
		fmt.Println()
		fmt.Print("Are you sure? (y/N): ")

		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			platform.PrintInfo("Uninstall cancelled")
			return nil
		}
	}

	// Step 1: Stop all running containers
	platform.PrintInfo("Stopping all CodeKiwi containers...")
	instances, err := state.List()
	if err == nil {
		for _, instance := range instances {
			if docker.IsContainerRunning(instance.ContainerName) {
				platform.PrintInfo(fmt.Sprintf("Stopping %s...", instance.ContainerName))
				docker.StopContainer(instance.ContainerName)
				docker.RemoveContainer(instance.ContainerName)
			}
		}
	}

	// Step 2: Remove Docker images
	if !keepImages {
		platform.PrintInfo("Removing Docker images...")
		removeDockerImages(cfg)
	}

	// Step 3: Remove configuration and data
	if !keepData {
		platform.PrintInfo(fmt.Sprintf("Removing configuration from %s...", cfg.InstallDir))
		if err := os.RemoveAll(cfg.InstallDir); err != nil {
			platform.PrintWarning(fmt.Sprintf("Failed to remove %s: %v", cfg.InstallDir, err))
		}
	}

	// Step 4: Remove binary from PATH
	platform.PrintInfo("Removing codekiwi from PATH...")
	if err := removeBinary(); err != nil {
		platform.PrintWarning(fmt.Sprintf("Failed to remove binary: %v", err))
		fmt.Println("\nPlease manually remove the codekiwi binary from your PATH")
	}

	platform.PrintSuccess("CodeKiwi has been uninstalled")

	if keepData {
		platform.PrintInfo(fmt.Sprintf("Configuration preserved in %s", cfg.InstallDir))
		fmt.Println("To completely remove, delete this directory manually")
	}

	fmt.Println("\nThank you for using CodeKiwi!")
	fmt.Println("To reinstall, visit: https://github.com/drasdp/codekiwi-cli")

	return nil
}

func removeDockerImages(cfg *config.Config) {
	// List of images to remove
	images := []string{
		cfg.GetFullImageName(),
		fmt.Sprintf("%s/%s:latest", cfg.ImageRegistry, cfg.ImageName),
		fmt.Sprintf("%s/%s", cfg.ImageRegistry, cfg.ImageName),
	}

	for _, image := range images {
		cmd := exec.Command("docker", "rmi", image)
		if err := cmd.Run(); err == nil {
			platform.PrintInfo(fmt.Sprintf("Removed image: %s", image))
		}
	}
}

func removeBinary() error {
	// Get the current executable path
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// Resolve symlinks
	realPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		realPath = exePath
	}

	switch runtime.GOOS {
	case "windows":
		// On Windows, the binary is typically in .codekiwi\bin
		// We can't delete ourselves while running, so create a batch file to do it
		batchFile := filepath.Join(os.TempDir(), "uninstall_codekiwi.bat")
		batchContent := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak > nul
del /f "%s"
echo CodeKiwi binary removed
del "%s"
`, realPath, batchFile)

		if err := os.WriteFile(batchFile, []byte(batchContent), 0755); err != nil {
			return err
		}

		cmd := exec.Command("cmd", "/c", "start", "/b", batchFile)
		return cmd.Start()

	case "darwin", "linux":
		// Check common installation locations
		locations := []string{
			"/usr/local/bin/codekiwi",
			"/usr/bin/codekiwi",
			filepath.Join(os.Getenv("HOME"), ".local", "bin", "codekiwi"),
			filepath.Join(os.Getenv("HOME"), ".codekiwi", "bin", "codekiwi"),
		}

		// Try to remove from each location
		removed := false
		for _, loc := range locations {
			if _, err := os.Stat(loc); err == nil {
				// Check if it's a symlink
				if info, err := os.Lstat(loc); err == nil && info.Mode()&os.ModeSymlink != 0 {
					// Remove symlink
					if err := os.Remove(loc); err == nil {
						platform.PrintInfo(fmt.Sprintf("Removed symlink: %s", loc))
						removed = true
					}
				} else if strings.Contains(realPath, ".codekiwi") {
					// If it's the actual binary in .codekiwi, just remove the symlink
					if loc != realPath {
						os.Remove(loc)
						removed = true
					}
				}
			}
		}

		if !removed {
			return fmt.Errorf("could not find codekiwi in standard locations")
		}
	}

	return nil
}