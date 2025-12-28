package main

import (
	"context"
	"fmt"
	"os"

	"github.com/agentplexus/omnitrust/inspector"
	"github.com/spf13/cobra"
)

var memoryCmd = &cobra.Command{
	Use:     "memory",
	Aliases: []string{"mem"},
	Short:   "Show memory usage",
	Long: `Display current system memory usage.

Shows total, used, free, and available memory with human-readable sizes.
Use --format=table for a colored ASCII table with progress bars.`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := inspector.GetMemory(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		output := inspector.FormatMemory(result, formatFlag)
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(memoryCmd)
}
