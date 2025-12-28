package main

import (
	"fmt"
	"os"

	"github.com/agentplexus/omnitrust/inspector"
	"github.com/spf13/cobra"
)

var securebootCmd = &cobra.Command{
	Use:     "secureboot",
	Aliases: []string{"sb", "boot"},
	Short:   "Show Secure Boot status",
	Long: `Display UEFI Secure Boot status.

Shows whether Secure Boot is enabled, the security mode, and any
relevant details about the boot security configuration.

On macOS, this shows Apple Secure Boot status (Full/Reduced/Permissive).
On Windows and Linux, this shows UEFI Secure Boot status.

Use --format=table for a colored ASCII table.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !inspector.IsSecureBootSupported() {
			fmt.Fprintln(os.Stderr, "Error: Secure Boot not supported on this platform")
			os.Exit(1)
		}

		result, err := inspector.GetSecureBootStatus()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		output := inspector.FormatSecureBoot(result, formatFlag)
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(securebootCmd)
}
