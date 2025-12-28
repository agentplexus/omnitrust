package main

import (
	"github.com/spf13/cobra"
)

var (
	formatFlag string
)

var rootCmd = &cobra.Command{
	Use:   "omnitrust",
	Short: "Cross-platform security posture assessment with MCP server support",
	Long: `OmniTrust provides unified security posture assessment tools across macOS, Windows,
and Linux. It can run as a Model Context Protocol (MCP) server for AI assistants,
or as standalone CLI commands.

Security Features:
  - TPM / Secure Enclave status (hardware security module)
  - Secure Boot verification (UEFI boot security)
  - Disk Encryption status (FileVault/BitLocker/LUKS)
  - Biometric capabilities (Touch ID/Face ID/Windows Hello)
  - Unified security summary with score and recommendations

System Metrics:
  - CPU usage monitoring (overall and per-core)
  - Memory usage statistics
  - Process listing with resource usage

Output formats:
  - JSON (default): Structured data for programmatic use
  - Table: Rich ASCII tables with ANSI colors and UTF-8 icons`,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&formatFlag, "format", "f", "json", "Output format: 'json' (default) or 'table'")
}
