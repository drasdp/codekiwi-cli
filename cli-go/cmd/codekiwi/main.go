package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/drasdp/codekiwi-cli/internal/commands"
	"github.com/spf13/cobra"
)

var (
	// Version information (set during build)
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "codekiwi",
	Short: "CodeKiwi - Docker-based Cloud IDE",
	Long: `CodeKiwi provides a complete development environment in your browser.
It runs your project in a Docker container with web-based code editor and live preview.`,
	Version: Version,
}

func init() {
	// Set version template
	rootCmd.SetVersionTemplate(fmt.Sprintf("CodeKiwi %s\nCommit: %s\nBuilt: %s\nOS/Arch: %s/%s\n",
		Version, Commit, Date, runtime.GOOS, runtime.GOARCH))

	// Initialize platform-specific settings
	initPlatform()

	// Add commands
	rootCmd.AddCommand(commands.StartCmd)
	rootCmd.AddCommand(commands.ListCmd)
	rootCmd.AddCommand(commands.KillCmd)
	rootCmd.AddCommand(commands.UpdateCmd)
	rootCmd.AddCommand(commands.UninstallCmd)
}

func initPlatform() {
	if runtime.GOOS == "windows" {
		// Enable UTF-8 on Windows
		// This helps with emoji and Unicode characters in cmd
		os.Setenv("CHCP", "65001")
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}