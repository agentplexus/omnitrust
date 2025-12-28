package main

import (
	"fmt"
	"os"

	"github.com/agentplexus/omnitrust/inspector"
	"github.com/spf13/cobra"
)

var securityChipCmd = &cobra.Command{
	Use:     "security-chip",
	Aliases: []string{"chip", "tpm", "se", "secureenclave"},
	Short:   "Show platform security chip status",
	Long: `Display platform security chip status.

On macOS, this shows Secure Enclave status including hardware-backed
key support and platform type (Apple Silicon or Intel T2).

On Windows and Linux, this shows TPM (Trusted Platform Module) presence,
version, manufacturer, and whether hardware key storage is supported.

Use --format=table for a colored ASCII table.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !inspector.IsTPMSupported() {
			fmt.Fprintln(os.Stderr, "Error: Platform security chip not supported on this platform")
			os.Exit(1)
		}

		result, err := inspector.GetTPMStatus()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		output := inspector.FormatTPM(result, formatFlag)
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(securityChipCmd)
}
