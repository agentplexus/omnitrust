package main

import (
	"context"
	"fmt"
	"os"

	"github.com/agentplexus/posture/inspector"
	"github.com/spf13/cobra"
)

var (
	processLimit int
)

var processesCmd = &cobra.Command{
	Use:     "processes",
	Aliases: []string{"ps", "proc"},
	Short:   "List running processes",
	Long: `List running processes with resource usage.

Shows PID, name, CPU usage, memory usage, and status for each process.
Results are sorted by CPU usage in descending order.
Use --limit to restrict the number of processes shown.
Use --format=table for a colored ASCII table.`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := inspector.ListProcesses(context.Background(), processLimit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		output := inspector.FormatProcessList(result, formatFlag)
		fmt.Println(output)
	},
}

func init() {
	processesCmd.Flags().IntVarP(&processLimit, "limit", "n", 0, "Maximum number of processes to show (0 for all)")
	rootCmd.AddCommand(processesCmd)
}
