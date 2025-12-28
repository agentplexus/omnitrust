---
marp: true
theme: agentplexus
paginate: true
color: #eaeaea
style: |
  section {
    font-family: 'SF Pro Display', -apple-system, BlinkMacSystemFont, sans-serif;
  }
  h1 {
    color: #00d4aa;
  }
  h2 {
    color: #00b894;
  }
  code {
    background-color: #16213e;
    border-radius: 4px;
  }
  pre {
    background-color: #16213e;
    border-radius: 8px;
  }
  table {
    font-size: 0.85em;
  }
  th {
    background-color: #00d4aa;
    color: #1a1a2e;
  }
  .columns {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }
  blockquote {
    border-left: 4px solid #00d4aa;
    padding-left: 1rem;
    font-style: italic;
    color: #b0bec5;
  }
  .green { color: #4caf50; }
  .yellow { color: #ffeb3b; }
  .red { color: #f44336; }
---

<!-- _paginate: false -->

# Posture 🛡️

## Cross-Platform Security Posture Assessment

💻 CLI | 🤖 MCP Server | 📦 Go Module

---

# What is Posture? 🔍

A **unified security inspection tool** for macOS, Windows, and Linux.

### 🔐 Security Features
- 🔑 **Platform Security Chip** - Secure Enclave (macOS) / TPM (Windows/Linux)
- 🛡️ **Secure Boot** - UEFI and Apple Secure Boot verification
- 🔒 **Disk Encryption** - FileVault / BitLocker / LUKS
- 👆 **Biometrics** - Touch ID / Face ID / Windows Hello / fprintd

### ➕ Plus
- 📊 Security scoring with actionable recommendations
- 📈 System metrics (CPU, memory, processes)

---

# Three Ways to Use Posture 🔀

| Method | Use Case | Best For |
|--------|----------|----------|
| 💻 **CLI** | Interactive terminal | DevOps, security audits |
| 🤖 **MCP Server** | AI assistants | Claude Desktop, automation |
| 📦 **Go Module** | Programmatic access | Custom applications |

**Same data, three interfaces.**

---

# CLI Usage 💻

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

---

# CLI Output: Security Summary 📊

```
Security Score: 75/100
Status: Good

Security Features:
+--------------------------+--------------+--------------------+
| Feature                  | Status       | Details            |
+--------------------------+--------------+--------------------+
| Secure Enclave           | Enabled      | secure_enclave     |
| Secure Boot              | Enabled      | full               |
| FileVault                | Disabled     | disabled           |
| Biometrics               | Enabled      | touch_id           |
+--------------------------+--------------+--------------------+

Recommendations:
  1. Enable FileVault to protect data at rest
```

---

# MCP Server for AI Assistants 🤖

Configure Claude Desktop:

**macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`

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

Start: `posture serve`

---

# MCP Tools 🛠️

| Tool | Description |
|------|-------------|
| `get_platform_security_chip` | Secure Enclave / TPM status |
| `get_secure_boot_status` | UEFI Secure Boot verification |
| `get_encryption_status` | Disk encryption status |
| `get_biometric_capabilities` | Biometric authentication status |
| `get_security_summary` | Unified posture with score |
| `get_cpu_usage` | CPU usage statistics |
| `get_memory` | Memory usage statistics |
| `list_processes` | Running process list |

---

# Go Module Usage 📦

### 📥 Installation

```bash
go get github.com/agentplexus/posture
```

### Import

```go
import "github.com/agentplexus/posture/inspector"
```

---

# Go Module: Security Summary 📦

```go
package main

import (
    "fmt"
    "github.com/agentplexus/posture/inspector"
)

func main() {
    summary, err := inspector.GetSecuritySummary()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Security Score: %d/100\n", summary.OverallScore)
    fmt.Printf("Status: %s\n", summary.OverallStatus)

    // Built-in table formatting
    fmt.Println(inspector.FormatSecuritySummaryTable(summary))
}
```

---

# Go Module: Individual Checks 📦

```go
// Platform Security Chip (Secure Enclave / TPM)
if inspector.IsTPMSupported() {
    tpm, _ := inspector.GetTPMStatus()
    fmt.Printf("Security Chip: %s (enabled: %v)\n", tpm.Type, tpm.Enabled)
}

// Disk Encryption
if inspector.IsEncryptionSupported() {
    enc, _ := inspector.GetEncryptionStatus()
    fmt.Printf("Encryption: %s (%s)\n", enc.Type, enc.Status)
}

// System Metrics
cpu, _ := inspector.GetCPUUsage(ctx)
fmt.Printf("CPU Usage: %.1f%%\n", cpu.OverallPercent)

mem, _ := inspector.GetMemory(ctx)
fmt.Printf("Memory: %.1f%% used\n", mem.UsedPercent)
```

---

# Go Module: Available Functions 📦

| Function | Description |
|----------|-------------|
| `GetSecuritySummary()` | Unified security posture with score |
| `GetTPMStatus()` | Platform security chip status |
| `GetSecureBootStatus()` | Secure Boot configuration |
| `GetEncryptionStatus()` | Disk encryption status |
| `GetBiometricCapabilities()` | Biometric authentication status |
| `GetCPUUsage(ctx)` | CPU usage statistics |
| `GetMemory(ctx)` | Memory usage statistics |
| `ListProcesses(ctx, limit)` | Running process list |

Each has a corresponding `IsXXXSupported()` function.

---

# Platform Support 🖥️

| Feature | macOS | Windows | Linux |
|---------|-------|---------|-------|
| Platform Security Chip | Secure Enclave | TPM 1.2/2.0 | TPM 2.0 |
| Secure Boot | Apple Secure Boot | UEFI Secure Boot | UEFI Secure Boot |
| Disk Encryption | FileVault | BitLocker | LUKS/dm-crypt |
| Biometrics | Touch ID/Face ID | Windows Hello | fprintd/Howdy |
| CPU/Memory/Processes | Yes | Yes | Yes |

---

# Architecture 🏗️

```
+------------------+     +---------------------------+
|   CLI / MCP      |     |      Go Applications      |
+------------------+     +---------------------------+
         |                           |
         v                           v
+--------------------------------------------------------+
|              inspector package                         |
|  +----------+ +----------+ +----------+ +-----------+  |
|  |  darwin  | | windows  | |  linux   | | common    |  |
|  |  (cgo)   | |  (WMI)   | | (sysfs)  | |(gopsutil) |  |
|  +----------+ +----------+ +----------+ +-----------+  |
+--------------------------------------------------------+
```

---

# Security Score Calculation 📊

Each feature contributes **25 points** to the total score:

| Score | Status | Description |
|-------|--------|-------------|
| 100 | Excellent | All security features enabled |
| 75 | Good | Most features enabled |
| 50 | Fair | Some features missing |
| 25 | Needs Improvement | Critical features missing |
| 0 | Critical | No security features enabled |

---

# JSON Output 📋

```json
{
  "platform": "darwin",
  "overall_score": 75,
  "overall_status": "good",
  "tpm": {
    "present": true,
    "enabled": true,
    "type": "secure_enclave"
  },
  "secure_boot": { "enabled": true, "mode": "full" },
  "encryption": { "enabled": false, "type": "filevault" },
  "biometrics": { "available": true, "type": "touch_id" },
  "recommendations": ["Enable FileVault to protect data at rest"]
}
```

---

# Security by Design 🔐

### ✅ What Posture Does

- 👁️ Read-only system inspection
- 🔍 Hardware capability verification
- 📊 Security posture assessment
- 📋 Process enumeration

### ❌ What Posture Does NOT Do

- 🚫 Access keychain or secrets
- 🚫 Extract cryptographic keys
- 🚫 Modify system settings
- 🚫 Execute arbitrary commands
- 🚫 Make network requests

---

# Rich Terminal Output 🎨

| Feature | Description |
|---------|-------------|
| **ANSI Colors** | Green/Yellow/Red status indicators |
| **Progress Bars** | Visual usage display |
| **Box Drawing** | Unicode table borders |
| **UTF-8 Icons** | Visual feature indicators |

### Color Coding

| Color | Meaning |
|-------|---------|
| <span class="green">Green</span> | Good / Enabled |
| <span class="yellow">Yellow</span> | Warning (70-90%) |
| <span class="red">Red</span> | Critical / Disabled |

---

# Installation 📥

### Pre-built Binaries

[GitHub Releases](https://github.com/agentplexus/posture/releases)

### Build from Source

```bash
git clone https://github.com/agentplexus/posture.git
cd posture
go build -o posture ./cmd/posture/
```

### Go Module

```bash
go get github.com/agentplexus/posture
```

---

# Cross-Compilation 🔧

```bash
# macOS (Apple Silicon & Intel)
GOOS=darwin GOARCH=arm64 go build -o posture-darwin-arm64 ./cmd/posture/
GOOS=darwin GOARCH=amd64 go build -o posture-darwin-amd64 ./cmd/posture/

# Linux
GOOS=linux GOARCH=amd64 go build -o posture-linux-amd64 ./cmd/posture/
GOOS=linux GOARCH=arm64 go build -o posture-linux-arm64 ./cmd/posture/

# Windows
GOOS=windows GOARCH=amd64 go build -o posture-windows.exe ./cmd/posture/
```

*Note: macOS Secure Enclave requires native compilation (cgo)*

---

# Dependencies 📦

| Package | Purpose |
|---------|---------|
| [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) | Official MCP Go SDK |
| [shirou/gopsutil/v4](https://github.com/shirou/gopsutil) | Cross-platform system metrics |
| [spf13/cobra](https://github.com/spf13/cobra) | CLI framework |
| [mattn/go-runewidth](https://github.com/mattn/go-runewidth) | Unicode width calculation |

---

# Use Cases 💼

<div class="columns">
<div>

### 🔒 Security Audits
- ✅ Compliance verification
- 🔍 Endpoint security assessment
- 🚦 Pre-deployment checks

### ⚙️ DevOps
- 🏗️ Infrastructure validation
- 🚧 CI/CD security gates
- 📡 Fleet monitoring

</div>
<div>

### 🤖 AI Assistants
- 💬 Real-time system queries
- 💡 Security recommendations
- 🔄 Automated remediation

### 🖥️ Applications
- 📊 Security dashboards
- 📱 MDM integrations
- 📜 Policy enforcement

</div>
</div>

---

# Demo: CLI 🎬

```bash
# Check security posture
$ posture summary -f table

Security Score: 75/100
Status: Good

Security Features:
| Secure Enclave  | Enabled  | secure_enclave |
| Secure Boot     | Enabled  | full           |
| FileVault       | Disabled | disabled       |
| Biometrics      | Enabled  | touch_id       |

Recommendations:
  1. Enable FileVault to protect data at rest
```

---

# Demo: MCP with Claude 🎬

> "What's the security status of this machine?"

Claude calls `get_security_summary` and responds:

> "Your machine has a security score of 75/100. Secure Enclave and Secure Boot are enabled, and Touch ID is configured. However, FileVault disk encryption is currently disabled. I recommend enabling FileVault to protect your data at rest."

**No hallucination. Real system data.**

---

# Demo: Go Module 🎬

```go
summary, _ := inspector.GetSecuritySummary()

if summary.OverallScore < 50 {
    alert.Send("Security posture critical!")
}

for _, rec := range summary.Recommendations {
    log.Printf("Action needed: %s", rec)
}
```

**Integrate security checks into your applications.**

---

# Key Takeaways 🎯

🔀 **Three interfaces, one tool** - CLI, MCP, and Go Module

🖥️ **Cross-platform** - macOS, Windows, Linux with native APIs

🔐 **Security-focused** - Read-only, no secrets exposed

🤖 **AI-ready** - MCP integration for Claude Desktop

👩‍💻 **Developer-friendly** - Go module for programmatic access

---

# Resources 📚

- **GitHub**: [github.com/agentplexus/posture](https://github.com/agentplexus/posture)
- **MCP Specification**: [modelcontextprotocol.io](https://modelcontextprotocol.io)
- **MCP Go SDK**: [github.com/modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk)
- **Claude Desktop**: [claude.ai/download](https://claude.ai/download)

---

# Get Started 🚀

```bash
# Install
go install github.com/agentplexus/posture/cmd/posture@latest

# Check your security posture
posture summary -f table

# Start MCP server for Claude
posture serve
```

---

# Thank You 🙏

## Posture 🛡️

**Cross-Platform Security Posture Assessment**

💻 CLI | 🤖 MCP Server | 📦 Go Module

[github.com/agentplexus/posture](https://github.com/agentplexus/posture)
