//go:build linux

package inspector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TPMResult contains TPM status information
type TPMResult struct {
	Present            bool     `json:"present"`
	Enabled            bool     `json:"enabled"`
	Version            string   `json:"version"`
	Manufacturer       string   `json:"manufacturer"`
	Type               string   `json:"type"`
	Platform           string   `json:"platform"`
	Capabilities       []string `json:"capabilities"`
	HardwareKeySupport bool     `json:"hardware_key_support"`
}

// GetTPMStatus returns the TPM status (Linux)
func GetTPMStatus() (*TPMResult, error) {
	// Check for TPM devices in /sys/class/tpm/
	tpmPath := "/sys/class/tpm"

	entries, err := os.ReadDir(tpmPath)
	if err != nil || len(entries) == 0 {
		// No TPM found
		return &TPMResult{
			Present:            false,
			Enabled:            false,
			Version:            "Not detected",
			Manufacturer:       "Unknown",
			Type:               "none",
			Platform:           "linux",
			Capabilities:       []string{},
			HardwareKeySupport: false,
		}, nil
	}

	// Use the first TPM device found (usually tpm0)
	tpmDevice := entries[0].Name()
	devicePath := filepath.Join(tpmPath, tpmDevice)

	// Read TPM version
	version := readSysFile(filepath.Join(devicePath, "tpm_version_major"))
	versionMinor := readSysFile(filepath.Join(devicePath, "tpm_version_minor"))

	tpmType := "tpm_1.2"
	versionStr := "1.2"
	if version == "2" {
		tpmType = "tpm_2.0"
		versionStr = "2.0"
	}
	if versionMinor != "" {
		versionStr = fmt.Sprintf("%s.%s", version, versionMinor)
	}

	// Read manufacturer info from device
	manufacturer := readSysFile(filepath.Join(devicePath, "device/vendor"))
	if manufacturer == "" {
		// Try to get from caps
		manufacturer = "Unknown"
	}

	// Check if device is accessible (enabled)
	_, devErr := os.Stat("/dev/" + tpmDevice)
	enabled := devErr == nil

	capabilities := []string{}
	if enabled {
		capabilities = append(capabilities, "hardware_key_generation")
		capabilities = append(capabilities, "hardware_key_storage")
		capabilities = append(capabilities, "platform_integrity")
		capabilities = append(capabilities, "secure_boot_support")
		if tpmType == "tpm_2.0" {
			capabilities = append(capabilities, "enhanced_authorization")
			capabilities = append(capabilities, "algorithm_agility")
		}
	}

	return &TPMResult{
		Present:            true,
		Enabled:            enabled,
		Version:            fmt.Sprintf("TPM %s", versionStr),
		Manufacturer:       manufacturer,
		Type:               tpmType,
		Platform:           "linux",
		Capabilities:       capabilities,
		HardwareKeySupport: enabled,
	}, nil
}

// readSysFile reads a sysfs file and returns trimmed content
func readSysFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// FormatTPMTable formats TPM status as a colored table
func FormatTPMTable(result *TPMResult) string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(Header(IconShield + " TPM Status"))
	sb.WriteString("\n")
	sb.WriteString(Muted(strings.Repeat("─", 55)))
	sb.WriteString("\n\n")

	// Platform badge
	sb.WriteString(BoldText("Platform: "))
	sb.WriteString(Info(IconChip + " Linux"))
	sb.WriteString("\n\n")

	// Status table
	sb.WriteString(TableTop(28, 22))
	sb.WriteString("\n")
	sb.WriteString(TableRowColored(
		Header(PadRight("Property", 28)),
		Header(PadRight("Value", 22)),
	))
	sb.WriteString("\n")
	sb.WriteString(TableSeparator(28, 22))
	sb.WriteString("\n")

	// Present
	sb.WriteString(TableRowColored(
		PadRight(IconShield+" TPM Present", 28),
		PadRight(BoolToStatusColored(result.Present), 22),
	))
	sb.WriteString("\n")

	// Enabled
	sb.WriteString(TableRowColored(
		PadRight(IconCheck+" Enabled", 28),
		PadRight(BoolToStatusColored(result.Enabled), 22),
	))
	sb.WriteString("\n")

	// Version
	sb.WriteString(TableRowColored(
		PadRight(IconInfo+" Version", 28),
		PadRight(Info(result.Version), 22),
	))
	sb.WriteString("\n")

	// Manufacturer
	sb.WriteString(TableRowColored(
		PadRight(IconDiamond+" Manufacturer", 28),
		PadRight(result.Manufacturer, 22),
	))
	sb.WriteString("\n")

	// Type
	typeDisplay := result.Type
	if result.Type == "tpm_2.0" {
		typeDisplay = Success("TPM 2.0")
	} else if result.Type == "tpm_1.2" {
		typeDisplay = Warning("TPM 1.2")
	}
	sb.WriteString(TableRowColored(
		PadRight(IconChip+" Type", 28),
		PadRight(typeDisplay, 22),
	))
	sb.WriteString("\n")

	// Hardware Key Support
	sb.WriteString(TableRowColored(
		PadRight(IconKey+" Hardware Key Support", 28),
		PadRight(BoolToStatusColored(result.HardwareKeySupport), 22),
	))
	sb.WriteString("\n")

	sb.WriteString(TableBottom(28, 22))
	sb.WriteString("\n\n")

	// Capabilities section
	if len(result.Capabilities) > 0 {
		sb.WriteString(BoldText("Capabilities:"))
		sb.WriteString("\n")
		sb.WriteString(Muted(strings.Repeat("─", 35)))
		sb.WriteString("\n")
		for _, cap := range result.Capabilities {
			sb.WriteString(fmt.Sprintf("  %s %s\n", Success(IconCheck), cap))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// FormatTPM formats TPM status in the specified format
func FormatTPM(result *TPMResult, format string) string {
	return FormatOutput(result, func() string {
		return FormatTPMTable(result)
	}, format)
}

// IsTPMSupported returns true on Linux
func IsTPMSupported() bool {
	return true
}
