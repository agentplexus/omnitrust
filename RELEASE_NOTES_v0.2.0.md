# Posture v0.2.0 Release Notes

**Release Date:** December 2024

This release introduces a standalone MCP server binary, aligning with the AgentPlexus convention for MCP server naming across projects.

## Breaking Changes

### MCP Server Moved to Standalone Binary

The MCP server has been moved from a subcommand (`posture serve`) to a dedicated binary (`mcp-posture`).

**Before (v0.1.0):**
```json
{
  "mcpServers": {
    "posture": {
      "command": "/path/to/posture",
      "args": ["serve"]
    }
  }
}
```

**After (v0.2.0):**
```json
{
  "mcpServers": {
    "posture": {
      "command": "/path/to/mcp-posture"
    }
  }
}
```

This change:
- Aligns with AgentPlexus `mcp-*` naming convention for MCP servers
- Simplifies Claude Desktop configuration (no args needed)
- Provides clearer separation between CLI tool and MCP server

## Bug Fixes

- Fixed MCP SDK jsonschema tag format causing server panic on startup

## Installation

### Pre-built Binaries

Download from the [Releases](https://github.com/agentplexus/posture/releases) page:
- `posture` - CLI tool for manual security assessment
- `mcp-posture` - MCP server for AI assistants

### Build from Source

```bash
git clone https://github.com/agentplexus/posture.git
cd posture

# Build CLI tool
go build -o posture ./cmd/posture/

# Build MCP server
go build -o mcp-posture ./cmd/mcp-posture/
```

## Migration Guide

1. Download or build the new `mcp-posture` binary
2. Update your Claude Desktop configuration:
   - Change `"command"` to point to `mcp-posture`
   - Remove the `"args": ["serve"]` line
3. Restart Claude Desktop

## Full Changelog

- `refactor: move MCP server to standalone cmd/mcp-posture binary`
- `fix: correct jsonschema tag format in MCP tool argument structs`
