//go:build linux

package inspector

import (
	"os"
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

// GetSecureBootStatus returns the Secure Boot status (Linux)
func GetSecureBootStatus() (*SecureBootResult, error) {
	result := &SecureBootResult{
		Platform: "linux",
	}

	// Check if booted in UEFI mode
	_, err := os.Stat("/sys/firmware/efi")
	if os.IsNotExist(err) {
		result.Enabled = false
		result.Mode = "legacy_bios"
		result.SecureBootType = "none"
		result.Details = "System booted in Legacy BIOS mode"
		return result, nil
	}

	result.SecureBootType = "uefi_secure_boot"

	// Check Secure Boot status from efivars
	// The SecureBoot variable is at:
	// /sys/firmware/efi/efivars/SecureBoot-8be4df61-93ca-11d2-aa0d-00e098032b8c
	secureBootPath := "/sys/firmware/efi/efivars/SecureBoot-8be4df61-93ca-11d2-aa0d-00e098032b8c"

	data, err := os.ReadFile(secureBootPath)
	if err != nil {
		// Try alternative path or mokutil
		result.Mode = "unknown"
		result.Details = "Unable to read Secure Boot variable (may require root)"

		// Check if secureboot directory exists as fallback
		if _, err := os.Stat("/sys/firmware/efi/efivars"); err == nil {
			// UEFI is present but can't read secure boot status
			result.Enabled = false
		}
		return result, nil
	}

	// The efivars format: first 4 bytes are attributes, then the value
	// SecureBoot value: 0 = disabled, 1 = enabled
	if len(data) >= 5 {
		secureBootValue := data[4]
		if secureBootValue == 1 {
			result.Enabled = true
			result.Mode = "enabled"
			result.Details = "UEFI Secure Boot is enabled"
		} else {
			result.Enabled = false
			result.Mode = "disabled"
			result.Details = "UEFI Secure Boot is disabled"
		}
	} else {
		result.Mode = "unknown"
		result.Details = "Unable to parse Secure Boot variable"
	}

	// Check SetupMode (indicates if keys can be modified)
	setupModePath := "/sys/firmware/efi/efivars/SetupMode-8be4df61-93ca-11d2-aa0d-00e098032b8c"
	if data, err := os.ReadFile(setupModePath); err == nil && len(data) >= 5 {
		if data[4] == 1 {
			result.Details += " (Setup Mode active - keys can be modified)"
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
	sb.WriteString(Info(IconChip + " Linux"))
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
	typeDisplay := result.SecureBootType
	if result.SecureBootType == "uefi_secure_boot" {
		typeDisplay = "UEFI Secure Boot"
	} else if result.SecureBootType == "none" {
		typeDisplay = Muted("Not Available")
	}
	sb.WriteString(TableRowColored(
		PadRight(IconShield+" Type", 24),
		PadRight(typeDisplay, 26),
	))
	sb.WriteString("\n")

	// Mode
	var modeDisplay string
	switch result.Mode {
	case "enabled":
		modeDisplay = Success("Enabled")
	case "disabled":
		modeDisplay = Warning("Disabled")
	case "legacy_bios":
		modeDisplay = Danger("Legacy BIOS")
	default:
		modeDisplay = Muted(result.Mode)
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

// IsSecureBootSupported returns true on Linux
func IsSecureBootSupported() bool {
	return true
}
