package main

import (
	"fmt"
	"os"

	"github.com/agentplexus/omnitrust/inspector"
	"github.com/spf13/cobra"
)

var summaryCmd = &cobra.Command{
	Use:     "summary",
	Aliases: []string{"sum", "status", "security"},
	Short:   "Show unified security summary",
	Long: `Display a unified security posture overview.

Checks all security features and provides:
  - Overall security score (0-100)
  - Status of TPM/Secure Enclave
  - Status of Secure Boot
  - Status of disk encryption
  - Status of biometric authentication
  - Recommendations for improving security

Use --format=table for a colored ASCII table with visual score bar.`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := inspector.GetSecuritySummary()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		output := inspector.FormatSecuritySummary(result, formatFlag)
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(summaryCmd)
}
