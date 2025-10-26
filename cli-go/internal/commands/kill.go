package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/drasdp/codekiwi-cli/internal/config"
	"github.com/drasdp/codekiwi-cli/internal/docker"
	"github.com/drasdp/codekiwi-cli/internal/platform"
	"github.com/drasdp/codekiwi-cli/internal/state"
	"github.com/spf13/cobra"
)

var (
	killAll bool
	force   bool
)

// KillCmd stops a CodeKiwi project
var KillCmd = &cobra.Command{
	Use:   "kill [path|container]",
	Short: "Stop a CodeKiwi instance",
	Long:  `Stop a running CodeKiwi development environment.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runKill,
}

func init() {
	KillCmd.Flags().BoolVarP(&killAll, "all", "a", false, "Stop all running instances")
	KillCmd.Flags().BoolVarP(&force, "force", "f", false, "Force stop without confirmation")
}

func runKill(cmd *cobra.Command, args []string) error {
	// Kill all instances
	if killAll {
		return killAllInstances()
	}

	// Need a target if not killing all
	if len(args) == 0 {
		return fmt.Errorf("please specify a project path or container name, or use --all flag")
	}

	target := args[0]

	// Try to find instance by path first
	absPath, _ := filepath.Abs(target)
	absPath = config.NormalizePath(absPath)

	instance, err := state.LoadByPath(absPath)
	if err != nil {
		return fmt.Errorf("failed to load instance state: %w", err)
	}

	// If not found by path, try by container name
	if instance == nil {
		// Look for instance by container name
		instances, err := state.List()
		if err != nil {
			return fmt.Errorf("failed to list instances: %w", err)
		}

		for _, inst := range instances {
			if inst.ContainerName == target {
				instance = inst
				break
			}
		}
	}

	// If still not found, check if it's a partial match
	if instance == nil {
		instances, err := state.List()
		if err != nil {
			return fmt.Errorf("failed to list instances: %w", err)
		}

		var matches []*state.Instance
		for _, inst := range instances {
			if strings.Contains(inst.ContainerName, target) ||
				strings.Contains(inst.ProjectPath, target) ||
				strings.Contains(filepath.Base(inst.ProjectPath), target) {
				matches = append(matches, inst)
			}
		}

		if len(matches) == 1 {
			instance = matches[0]
		} else if len(matches) > 1 {
			platform.PrintError(fmt.Sprintf("Multiple instances match '%s':", target))
			for _, inst := range matches {
				fmt.Printf("  - %s (%s)\n", inst.ContainerName, inst.ProjectPath)
			}
			return fmt.Errorf("please be more specific")
		}
	}

	if instance == nil {
		return fmt.Errorf("no instance found for '%s'", target)
	}

	// Check if container is running
	if !docker.IsContainerRunning(instance.ContainerName) {
		platform.PrintWarning(fmt.Sprintf("Container %s is not running", instance.ContainerName))
		// Clean up state file
		instance.Delete()
		return nil
	}

	// Confirm before stopping
	if !force {
		platform.PrintWarning(fmt.Sprintf("Stopping CodeKiwi for %s", instance.ProjectPath))
		fmt.Print("Are you sure? (y/N): ")

		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			platform.PrintInfo("Cancelled")
			return nil
		}
	}

	// Stop the container
	platform.PrintInfo(fmt.Sprintf("Stopping container %s...", instance.ContainerName))

	if err := docker.ComposeDown(instance.ContainerName, instance.ProjectPath); err != nil {
		// Try direct docker stop if compose fails
		platform.PrintWarning("Trying alternative stop method...")
		if err := docker.StopContainer(instance.ContainerName); err != nil {
			return fmt.Errorf("failed to stop container: %w", err)
		}
		if err := docker.RemoveContainer(instance.ContainerName); err != nil {
			platform.PrintWarning("Failed to remove container")
		}
	}

	// Delete state file
	if err := instance.Delete(); err != nil {
		platform.PrintWarning("Failed to clean up state file")
	}

	platform.PrintSuccess(fmt.Sprintf("Stopped CodeKiwi for %s", instance.ProjectPath))
	return nil
}

func killAllInstances() error {
	// Load all instances
	instances, err := state.List()
	if err != nil {
		return fmt.Errorf("failed to load instances: %w", err)
	}

	// Filter running instances
	var runningInstances []*state.Instance
	for _, instance := range instances {
		if docker.IsContainerRunning(instance.ContainerName) {
			runningInstances = append(runningInstances, instance)
		} else {
			// Clean up stale state
			instance.Delete()
		}
	}

	if len(runningInstances) == 0 {
		platform.PrintInfo("No running CodeKiwi instances")
		return nil
	}

	// Confirm before stopping all
	if !force {
		platform.PrintWarning(fmt.Sprintf("This will stop %d running instance(s):", len(runningInstances)))
		for _, inst := range runningInstances {
			fmt.Printf("  - %s (%s)\n", inst.ContainerName, inst.ProjectPath)
		}
		fmt.Print("Are you sure? (y/N): ")

		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			platform.PrintInfo("Cancelled")
			return nil
		}
	}

	// Stop all instances
	successCount := 0
	failCount := 0

	for _, instance := range runningInstances {
		platform.PrintInfo(fmt.Sprintf("Stopping %s...", instance.ContainerName))

		if err := docker.ComposeDown(instance.ContainerName, instance.ProjectPath); err != nil {
			// Try direct docker stop if compose fails
			if err := docker.StopContainer(instance.ContainerName); err != nil {
				platform.PrintError(fmt.Sprintf("Failed to stop %s: %v", instance.ContainerName, err))
				failCount++
				continue
			}
			docker.RemoveContainer(instance.ContainerName)
		}

		// Delete state file
		instance.Delete()
		successCount++
	}

	// Print summary
	if successCount > 0 {
		platform.PrintSuccess(fmt.Sprintf("Stopped %d instance(s)", successCount))
	}
	if failCount > 0 {
		platform.PrintError(fmt.Sprintf("Failed to stop %d instance(s)", failCount))
		return fmt.Errorf("some instances failed to stop")
	}

	return nil
}