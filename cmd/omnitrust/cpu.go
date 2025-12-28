package main

import (
	"context"
	"fmt"
	"os"

	"github.com/agentplexus/omnitrust/inspector"
	"github.com/spf13/cobra"
)

var cpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "Show CPU usage",
	Long: `Display current system CPU usage.

Shows overall CPU usage percentage and per-core usage statistics.
Use --format=table for a colored ASCII table with progress bars.`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := inspector.GetCPUUsage(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		output := inspector.FormatCPUUsage(result, formatFlag)
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(cpuCmd)
}
