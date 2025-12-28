package main

import (
	"fmt"
	"os"

	"github.com/agentplexus/omnitrust/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MCP server",
	Long: `Start the Model Context Protocol (MCP) server over stdio.

This command starts the MCP server that exposes system inspection tools
to AI assistants like Claude Desktop. The server communicates over
stdin/stdout using JSON-RPC.

Configure in Claude Desktop's claude_desktop_config.json:
  {
    "mcpServers": {
      "system-inspector": {
        "command": "/path/to/mcp-system-inspector",
        "args": ["serve"]
      }
    }
  }`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := server.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
