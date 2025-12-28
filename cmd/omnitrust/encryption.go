package main

import (
	"fmt"
	"os"

	"github.com/agentplexus/omnitrust/inspector"
	"github.com/spf13/cobra"
)

var encryptionCmd = &cobra.Command{
	Use:     "encryption",
	Aliases: []string{"enc", "fv", "filevault", "bitlocker", "luks"},
	Short:   "Show disk encryption status",
	Long: `Display disk encryption status.

On macOS, this shows FileVault status.
On Windows, this shows BitLocker status.
On Linux, this shows LUKS/dm-crypt encryption status.

Shows whether encryption is enabled and lists encrypted volumes.

Use --format=table for a colored ASCII table.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !inspector.IsEncryptionSupported() {
			fmt.Fprintln(os.Stderr, "Error: Encryption status not supported on this platform")
			os.Exit(1)
		}

		result, err := inspector.GetEncryptionStatus()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		output := inspector.FormatEncryption(result, formatFlag)
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(encryptionCmd)
}
