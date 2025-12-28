package server

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/agentplexus/posture/inspector"
)

// Tool argument types - System metrics
type GetCPUUsageArgs struct {
	Format string `json:"format,omitempty" jsonschema:"description=Output format: 'json' (default) or 'table'"`
}

type GetMemoryArgs struct {
	Format string `json:"format,omitempty" jsonschema:"description=Output format: 'json' (default) or 'table'"`
}

type ListProcessesArgs struct {
	Limit  int    `json:"limit,omitempty" jsonschema:"description=Maximum number of processes to return (0 for all)"`
	Format string `json:"format,omitempty" jsonschema:"description=Output format: 'json' (default) or 'table'"`
}

// Tool argument types - Security tools
type GetPlatformSecurityChipArgs struct {
	Format string `json:"format,omitempty" jsonschema:"description=Output format: 'json' (default) or 'table'"`
}

type GetSecureBootStatusArgs struct {
	Format string `json:"format,omitempty" jsonschema:"description=Output format: 'json' (default) or 'table'"`
}

type GetEncryptionStatusArgs struct {
	Format string `json:"format,omitempty" jsonschema:"description=Output format: 'json' (default) or 'table'"`
}

type GetBiometricCapabilitiesArgs struct {
	Format string `json:"format,omitempty" jsonschema:"description=Output format: 'json' (default) or 'table'"`
}

type GetSecuritySummaryArgs struct {
	Format string `json:"format,omitempty" jsonschema:"description=Output format: 'json' (default) or 'table'"`
}

// System metric handlers

func handleGetCPUUsage(ctx context.Context, req *mcp.CallToolRequest, args GetCPUUsageArgs) (*mcp.CallToolResult, any, error) {
	result, err := inspector.GetCPUUsage(ctx)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
			IsError: true,
		}, nil, nil
	}

	output := inspector.FormatCPUUsage(result, args.Format)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output},
		},
	}, nil, nil
}

func handleGetMemory(ctx context.Context, req *mcp.CallToolRequest, args GetMemoryArgs) (*mcp.CallToolResult, any, error) {
	result, err := inspector.GetMemory(ctx)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
			IsError: true,
		}, nil, nil
	}

	output := inspector.FormatMemory(result, args.Format)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output},
		},
	}, nil, nil
}

func handleListProcesses(ctx context.Context, req *mcp.CallToolRequest, args ListProcessesArgs) (*mcp.CallToolResult, any, error) {
	result, err := inspector.ListProcesses(ctx, args.Limit)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
			IsError: true,
		}, nil, nil
	}

	output := inspector.FormatProcessList(result, args.Format)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output},
		},
	}, nil, nil
}

// Security tool handlers

func handleGetPlatformSecurityChip(_ context.Context, req *mcp.CallToolRequest, args GetPlatformSecurityChipArgs) (*mcp.CallToolResult, any, error) {
	result, err := inspector.GetTPMStatus()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
			IsError: true,
		}, nil, nil
	}

	output := inspector.FormatTPM(result, args.Format)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output},
		},
	}, nil, nil
}

func handleGetSecureBootStatus(_ context.Context, req *mcp.CallToolRequest, args GetSecureBootStatusArgs) (*mcp.CallToolResult, any, error) {
	result, err := inspector.GetSecureBootStatus()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
			IsError: true,
		}, nil, nil
	}

	output := inspector.FormatSecureBoot(result, args.Format)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output},
		},
	}, nil, nil
}

func handleGetEncryptionStatus(_ context.Context, req *mcp.CallToolRequest, args GetEncryptionStatusArgs) (*mcp.CallToolResult, any, error) {
	result, err := inspector.GetEncryptionStatus()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
			IsError: true,
		}, nil, nil
	}

	output := inspector.FormatEncryption(result, args.Format)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output},
		},
	}, nil, nil
}

func handleGetBiometricCapabilities(_ context.Context, req *mcp.CallToolRequest, args GetBiometricCapabilitiesArgs) (*mcp.CallToolResult, any, error) {
	result, err := inspector.GetBiometricCapabilities()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
			IsError: true,
		}, nil, nil
	}

	output := inspector.FormatBiometricCapabilities(result, args.Format)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output},
		},
	}, nil, nil
}

func handleGetSecuritySummary(_ context.Context, req *mcp.CallToolRequest, args GetSecuritySummaryArgs) (*mcp.CallToolResult, any, error) {
	result, err := inspector.GetSecuritySummary()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
			IsError: true,
		}, nil, nil
	}

	output := inspector.FormatSecuritySummary(result, args.Format)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output},
		},
	}, nil, nil
}

// NewMCPServer creates and configures a new MCP server
func NewMCPServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "posture",
		Version: "1.0.0",
	}, nil)

	// ============================================
	// Security Tools (Primary Focus)
	// ============================================

	// Platform Security Chip status (TPM on Windows/Linux, Secure Enclave on macOS)
	if inspector.IsTPMSupported() {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "get_platform_security_chip",
			Description: "Returns platform security chip status: Secure Enclave on macOS, TPM (Trusted Platform Module) on Windows/Linux. Includes presence, version, manufacturer, and hardware key support capabilities. Use format='table' for colored ASCII table output.",
		}, handleGetPlatformSecurityChip)
	}

	// Secure Boot status (all platforms)
	if inspector.IsSecureBootSupported() {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "get_secure_boot_status",
			Description: "Returns UEFI Secure Boot status including whether it's enabled, the security mode, and boot policy. Use format='table' for colored ASCII table output.",
		}, handleGetSecureBootStatus)
	}

	// Disk Encryption status (all platforms)
	if inspector.IsEncryptionSupported() {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "get_encryption_status",
			Description: "Returns disk encryption status (FileVault on macOS, BitLocker on Windows, LUKS on Linux) including whether encryption is enabled and which volumes are encrypted. Use format='table' for colored ASCII table output.",
		}, handleGetEncryptionStatus)
	}

	// Biometric capabilities (all platforms)
	if inspector.IsBiometricsSupported() {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "get_biometric_capabilities",
			Description: "Returns biometric authentication capabilities including Touch ID/fingerprint, Face ID/facial recognition availability and enrollment status. On Windows this includes Windows Hello status. Use format='table' for colored ASCII table output.",
		}, handleGetBiometricCapabilities)
	}

	// Security Summary (all platforms)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_security_summary",
		Description: "Returns a unified security posture overview including platform security chip (Secure Enclave/TPM), Secure Boot, disk encryption, and biometric status with an overall security score and recommendations. Use format='table' for colored ASCII table output.",
	}, handleGetSecuritySummary)

	// ============================================
	// System Metrics Tools (Bonus utilities)
	// ============================================

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_cpu_usage",
		Description: "Returns current system CPU usage percentage, both overall and per-core. Use format='table' for colored ASCII table output with progress bars.",
	}, handleGetCPUUsage)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_memory",
		Description: "Returns current system memory usage including total, used, free, and available memory. Use format='table' for colored ASCII table output with progress bars.",
	}, handleGetMemory)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_processes",
		Description: "Lists running processes with their PID, name, CPU usage, memory usage, and status. Results are sorted by CPU usage. Use format='table' for colored ASCII table output.",
	}, handleListProcesses)

	return server
}

// Run starts the MCP server on stdio
func Run() error {
	server := NewMCPServer()
	return server.Run(context.Background(), &mcp.StdioTransport{})
}
