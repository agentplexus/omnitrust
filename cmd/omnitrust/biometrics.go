package main

import (
	"fmt"
	"os"

	"github.com/agentplexus/omnitrust/inspector"
	"github.com/spf13/cobra"
)

var biometricsCmd = &cobra.Command{
	Use:     "biometrics",
	Aliases: []string{"bio"},
	Short:   "Show biometric capabilities (macOS only)",
	Long: `Display macOS biometric capabilities.

Shows Touch ID and Face ID availability and enrollment status.
This command is only available on macOS.
Use --format=table for a colored ASCII table.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !inspector.IsBiometricsSupported() {
			fmt.Fprintln(os.Stderr, "Error: Biometrics are only available on macOS")
			os.Exit(1)
		}

		result, err := inspector.GetBiometricCapabilities()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		output := inspector.FormatBiometricCapabilities(result, formatFlag)
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(biometricsCmd)
}
