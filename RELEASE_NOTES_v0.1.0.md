# Posture v0.1.0 Release Notes

**Release Date:** December 2024

Posture is a cross-platform security posture assessment tool with Model Context Protocol (MCP) server support. This initial release provides unified security inspection across macOS, Windows, and Linux, enabling AI assistants like Claude to query hardware security modules, boot security, disk encryption, and biometric capabilities.

## Highlights

- **Cross-platform security assessment** - Unified API for macOS, Windows, and Linux
- **MCP server integration** - Works with Claude Desktop and other MCP-compatible AI assistants
- **Three interfaces** - CLI, MCP server, and Go module for programmatic access
- **Security scoring** - Automated security posture scoring with actionable recommendations

## Features

### Security Assessment Tools

| Feature | macOS | Windows | Linux |
|---------|-------|---------|-------|
| Platform Security Chip | Secure Enclave | TPM 1.2/2.0 | TPM 2.0 |
| Secure Boot | Apple Secure Boot | UEFI Secure Boot | UEFI Secure Boot |
| Disk Encryption | FileVault | BitLocker | LUKS/dm-crypt |
| Biometrics | Touch ID / Face ID | Windows Hello | fprintd / Howdy |
| Security Summary | Full support | Full support | Full support |

### System Metrics

- **CPU Usage** - Overall and per-core monitoring with visual progress bars
- **Memory Usage** - Total, used, free, and available memory statistics
- **Process List** - Running processes with CPU/memory usage, sorted by resource consumption

### Output Formats

- **JSON** (default) - Structured data for programmatic use and AI consumption
- **Table** - Rich ASCII tables with ANSI colors and UTF-8 icons for human readability

## Usage

### CLI

```bash
# Security summary with score
posture summary -f table

# Individual security checks
posture security-chip -f table
posture secureboot -f table
posture encryption -f table
posture biometrics -f table

# System metrics
posture cpu -f table
posture memory -f table
posture processes -n 10 -f table
```

### MCP Server (Claude Desktop)

Add to your Claude Desktop configuration:

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

**Available MCP Tools:**

| Tool | Description |
|------|-------------|
| `get_platform_security_chip` | Secure Enclave (macOS) / TPM (Windows/Linux) status |
| `get_secure_boot_status` | UEFI Secure Boot verification |
| `get_encryption_status` | Disk encryption (FileVault/BitLocker/LUKS) |
| `get_biometric_capabilities` | Biometric authentication status |
| `get_security_summary` | Unified security posture with score |
| `get_cpu_usage` | CPU usage statistics |
| `get_memory` | Memory usage statistics |
| `list_processes` | Running process list |

### Go Module

```bash
go get github.com/agentplexus/posture
```

```go
import "github.com/agentplexus/posture/inspector"

summary, err := inspector.GetSecuritySummary()
fmt.Printf("Security Score: %d/100\n", summary.OverallScore)
```

## Installation

### Pre-built Binaries

Download from the [Releases](https://github.com/agentplexus/posture/releases) page.

### Build from Source

Requires Go 1.23 or later.

```bash
git clone https://github.com/agentplexus/posture.git
cd posture
go build -o posture ./cmd/posture/
```

### Cross-compilation

```bash
# macOS
GOOS=darwin GOARCH=arm64 go build -o posture-darwin-arm64 ./cmd/posture/
GOOS=darwin GOARCH=amd64 go build -o posture-darwin-amd64 ./cmd/posture/

# Linux
GOOS=linux GOARCH=amd64 go build -o posture-linux-amd64 ./cmd/posture/
GOOS=linux GOARCH=arm64 go build -o posture-linux-arm64 ./cmd/posture/

# Windows
GOOS=windows GOARCH=amd64 go build -o posture-windows-amd64.exe ./cmd/posture/
```

> **Note:** Cross-compiling for macOS from other platforms will not include Secure Enclave support due to cgo dependencies.

## Security Considerations

Posture is designed with security in mind:

- **Read-only operations** - No system modifications are possible
- **No secrets exposed** - Does not access keychain, passwords, or private keys
- **Non-invasive checks** - Only tests capability, never extracts keys
- **Process listing is informational** - Cannot terminate or modify processes

## Dependencies

- [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) v1.2.0 - Official MCP Go SDK
- [shirou/gopsutil/v4](https://github.com/shirou/gopsutil) v4.25.11 - Cross-platform system metrics
- [spf13/cobra](https://github.com/spf13/cobra) v1.10.2 - CLI framework

## Known Limitations

- macOS Secure Enclave detection requires native compilation (cgo)
- Some security checks may require elevated privileges (admin/root)
- Windows TPM detection uses WMI which may have performance overhead on some systems

## What's Next

Planned for future releases:

- Additional security checks (firewall status, antivirus, OS updates)
- Remote attestation support
- Configuration file support
- Prometheus metrics export

## Contributors

- Initial development by AgentPlexus team

## License

MIT License - see [LICENSE](LICENSE) for details.
