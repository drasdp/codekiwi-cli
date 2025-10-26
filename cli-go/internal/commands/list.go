package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/drasdp/codekiwi-cli/internal/docker"
	"github.com/drasdp/codekiwi-cli/internal/platform"
	"github.com/drasdp/codekiwi-cli/internal/state"
	"github.com/spf13/cobra"
)

var (
	showAll bool
	quiet   bool
)

// ListCmd lists all running CodeKiwi instances
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List running CodeKiwi instances",
	Long:  `List all running CodeKiwi development environments.`,
	Aliases: []string{"ls"},
	RunE:  runList,
}

func init() {
	ListCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all instances including stopped ones")
	ListCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Only display container names")
}

func runList(cmd *cobra.Command, args []string) error {
	// Load all instances
	instances, err := state.List()
	if err != nil {
		return fmt.Errorf("failed to load instances: %w", err)
	}

	if len(instances) == 0 {
		platform.PrintInfo("No CodeKiwi instances found")
		return nil
	}

	// Filter running instances
	var runningInstances []*state.Instance
	var stoppedInstances []*state.Instance

	for _, instance := range instances {
		if docker.IsContainerRunning(instance.ContainerName) {
			runningInstances = append(runningInstances, instance)
		} else {
			stoppedInstances = append(stoppedInstances, instance)
			// Clean up stale state file
			if !showAll {
				instance.Delete()
			}
		}
	}

	// Quiet mode - just print container names
	if quiet {
		for _, instance := range runningInstances {
			fmt.Println(instance.ContainerName)
		}
		if showAll {
			for _, instance := range stoppedInstances {
				fmt.Println(instance.ContainerName)
			}
		}
		return nil
	}

	// Print running instances
	if len(runningInstances) > 0 {
		platform.PrintInfo(fmt.Sprintf("Running CodeKiwi instances (%d):", len(runningInstances)))
		fmt.Println()

		// Create tabwriter for aligned output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		// Print header
		fmt.Fprintln(w, "CONTAINER\tPROJECT\tPORT\tUPTIME\tPATH")
		fmt.Fprintln(w, "─────────\t───────\t────\t──────\t────")

		// Print instances
		for _, instance := range runningInstances {
			projectName := filepath.Base(instance.ProjectPath)
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n",
				instance.ContainerName,
				projectName,
				instance.WebPort,
				instance.GetUptime(),
				instance.ProjectPath,
			)
		}

		w.Flush()
		fmt.Println()
	} else {
		platform.PrintInfo("No running CodeKiwi instances")
	}

	// Print stopped instances if --all flag is used
	if showAll && len(stoppedInstances) > 0 {
		fmt.Println()
		platform.PrintWarning(fmt.Sprintf("Stopped instances (%d):", len(stoppedInstances)))
		fmt.Println()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "CONTAINER\tPROJECT\tSTOPPED\tPATH")
		fmt.Fprintln(w, "─────────\t───────\t───────\t────")

		for _, instance := range stoppedInstances {
			projectName := filepath.Base(instance.ProjectPath)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				instance.ContainerName,
				projectName,
				instance.StartedAt.Format("2006-01-02 15:04"),
				instance.ProjectPath,
			)
		}

		w.Flush()
		fmt.Println()
	}

	// Print usage hints
	if len(runningInstances) > 0 {
		fmt.Println("To stop an instance, run:")
		fmt.Println("  codekiwi kill <project-path>")
		fmt.Println()
		fmt.Println("To view logs, run:")
		fmt.Println("  codekiwi logs <project-path>")
		fmt.Println()
	}

	return nil
}