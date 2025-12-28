//go:build darwin

package inspector

import (
	"os/exec"
	"strings"
)

// SecureBootResult contains Secure Boot status information
type SecureBootResult struct {
	Enabled        bool   `json:"enabled"`
	Platform       string `json:"platform"`
	Mode           string `json:"mode"`
	PolicyVersion  string `json:"policy_version,omitempty"`
	SecureBootType string `json:"secure_boot_type"`
	Details        string `json:"details,omitempty"`
}

// GetSecureBootStatus returns the Secure Boot status (macOS)
func GetSecureBootStatus() (*SecureBootResult, error) {
	result := &SecureBootResult{
		Platform: "darwin",
	}

	// Check boot policy using bputil (Apple Silicon) or csrutil/nvram (Intel)
	// First, check if we're on Apple Silicon
	out, err := exec.Command("sysctl", "-n", "hw.optional.arm64").Output()
	isAppleSilicon := err == nil && strings.TrimSpace(string(out)) == "1"

	if isAppleSilicon {
		// Apple Silicon - use bputil
		result.SecureBootType = "apple_secure_boot"

		// Try to get security mode
		out, err := exec.Command("bputil", "-d").Output()
		if err == nil {
			output := string(out)
			if strings.Contains(output, "Full Security") {
				result.Enabled = true
				result.Mode = "full"
				result.Details = "Full Security Mode"
			} else if strings.Contains(output, "Reduced Security") {
				result.Enabled = true
				result.Mode = "reduced"
				result.Details = "Reduced Security Mode"
			} else if strings.Contains(output, "Permissive Security") {
				result.Enabled = false
				result.Mode = "permissive"
				result.Details = "Permissive Security Mode"
			} else {
				// Default to enabled on Apple Silicon
				result.Enabled = true
				result.Mode = "unknown"
			}
		} else {
			// bputil requires admin privileges, assume enabled by default on Apple Silicon
			result.Enabled = true
			result.Mode = "assumed_full"
			result.Details = "Apple Silicon default (verification requires admin)"
		}
	} else {
		// Intel Mac - check for T2 secure boot
		result.SecureBootType = "t2_secure_boot"

		// Try nvram to check secure boot
		out, err := exec.Command("nvram", "94b73556-2197-4702-82a8-3e1337dafbfb:AppleSecureBootPolicy").Output()
		if err == nil {
			output := strings.TrimSpace(string(out))
			if strings.Contains(output, "%02") || strings.Contains(output, "2") {
				result.Enabled = true
				result.Mode = "full"
				result.Details = "Full Security"
			} else if strings.Contains(output, "%01") || strings.Contains(output, "1") {
				result.Enabled = true
				result.Mode = "medium"
				result.Details = "Medium Security"
			} else {
				result.Enabled = false
				result.Mode = "none"
				result.Details = "No Security"
			}
		} else {
			// Check if T2 is present (indicates secure boot capability)
			out, err := exec.Command("system_profiler", "SPiBridgeDataType").Output()
			if err == nil && strings.Contains(string(out), "T2") {
				result.Enabled = true
				result.Mode = "assumed"
				result.Details = "T2 chip detected (verification requires admin)"
			} else {
				// No T2, no secure boot on Intel
				result.Enabled = false
				result.Mode = "unavailable"
				result.Details = "No T2 chip - Secure Boot not available"
				result.SecureBootType = "none"
			}
		}
	}

	return result, nil
}

// FormatSecureBootTable formats Secure Boot status as a colored table
func FormatSecureBootTable(result *SecureBootResult) string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(Header(IconLock + " Secure Boot Status"))
	sb.WriteString("\n")
	sb.WriteString(Muted(strings.Repeat("─", 55)))
	sb.WriteString("\n\n")

	// Platform badge
	sb.WriteString(BoldText("Platform: "))
	sb.WriteString(Info(IconApple + " macOS"))
	sb.WriteString("\n\n")

	// Status table
	sb.WriteString(TableTop(24, 26))
	sb.WriteString("\n")
	sb.WriteString(TableRowColored(
		Header(PadRight("Property", 24)),
		Header(PadRight("Value", 26)),
	))
	sb.WriteString("\n")
	sb.WriteString(TableSeparator(24, 26))
	sb.WriteString("\n")

	// Enabled
	sb.WriteString(TableRowColored(
		PadRight(IconLock+" Secure Boot Enabled", 24),
		PadRight(BoolToStatusColored(result.Enabled), 26),
	))
	sb.WriteString("\n")

	// Type
	var typeDisplay string
	switch result.SecureBootType {
	case "apple_secure_boot":
		typeDisplay = "Apple Secure Boot"
	case "t2_secure_boot":
		typeDisplay = "T2 Secure Boot"
	default:
		typeDisplay = result.SecureBootType
	}
	sb.WriteString(TableRowColored(
		PadRight(IconShield+" Type", 24),
		PadRight(typeDisplay, 26),
	))
	sb.WriteString("\n")

	// Mode
	modeDisplay := result.Mode
	switch result.Mode {
	case "full":
		modeDisplay = Success("Full Security")
	case "reduced":
		modeDisplay = Warning("Reduced Security")
	case "permissive", "none":
		modeDisplay = Danger("Permissive/None")
	case "medium":
		modeDisplay = Warning("Medium Security")
	}
	sb.WriteString(TableRowColored(
		PadRight(IconStatus+" Mode", 24),
		PadRight(modeDisplay, 26),
	))
	sb.WriteString("\n")

	sb.WriteString(TableBottom(24, 26))
	sb.WriteString("\n")

	// Details if available
	if result.Details != "" {
		sb.WriteString("\n")
		sb.WriteString(Muted("Details: " + result.Details))
		sb.WriteString("\n")
	}

	return sb.String()
}

// FormatSecureBoot formats Secure Boot status in the specified format
func FormatSecureBoot(result *SecureBootResult, format string) string {
	return FormatOutput(result, func() string {
		return FormatSecureBootTable(result)
	}, format)
}

// IsSecureBootSupported returns true on macOS
func IsSecureBootSupported() bool {
	return true
}
