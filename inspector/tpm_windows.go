//go:build windows

package inspector

import (
	"fmt"
	"strings"

	"github.com/yusufpapurcu/wmi"
)

// Win32_Tpm represents WMI TPM class
type Win32_Tpm struct {
	IsActivated_InitialValue   bool
	IsEnabled_InitialValue     bool
	IsOwned_InitialValue       bool
	ManufacturerId             uint32
	ManufacturerIdTxt          string
	ManufacturerVersion        string
	ManufacturerVersionFull20  string
	ManufacturerVersionInfo    string
	PhysicalPresenceVersionInfo string
	SpecVersion                string
}

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

// GetTPMStatus returns the TPM status (Windows)
func GetTPMStatus() (*TPMResult, error) {
	var tpmInfo []Win32_Tpm

	// Query WMI for TPM information
	// Note: Requires running as administrator for full access
	query := "SELECT * FROM Win32_Tpm"
	err := wmi.QueryNamespace(query, &tpmInfo, `root\cimv2\Security\MicrosoftTpm`)

	if err != nil || len(tpmInfo) == 0 {
		// TPM not found or not accessible
		return &TPMResult{
			Present:            false,
			Enabled:            false,
			Version:            "Not detected",
			Manufacturer:       "Unknown",
			Type:               "none",
			Platform:           "windows",
			Capabilities:       []string{},
			HardwareKeySupport: false,
		}, nil
	}

	tpm := tpmInfo[0]

	// Determine TPM version from spec version
	version := tpm.SpecVersion
	tpmType := "tpm_1.2"
	if strings.Contains(version, "2.0") {
		tpmType = "tpm_2.0"
	}

	capabilities := []string{}
	if tpm.IsEnabled_InitialValue {
		capabilities = append(capabilities, "hardware_key_generation")
		capabilities = append(capabilities, "hardware_key_storage")
		capabilities = append(capabilities, "platform_integrity")
		capabilities = append(capabilities, "secure_boot_support")
		if tpmType == "tpm_2.0" {
			capabilities = append(capabilities, "enhanced_authorization")
			capabilities = append(capabilities, "algorithm_agility")
		}
	}

	manufacturer := tpm.ManufacturerIdTxt
	if manufacturer == "" {
		manufacturer = fmt.Sprintf("ID: %d", tpm.ManufacturerId)
	}

	return &TPMResult{
		Present:            true,
		Enabled:            tpm.IsEnabled_InitialValue,
		Version:            version,
		Manufacturer:       manufacturer,
		Type:               tpmType,
		Platform:           "windows",
		Capabilities:       capabilities,
		HardwareKeySupport: tpm.IsEnabled_InitialValue && tpm.IsActivated_InitialValue,
	}, nil
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
	sb.WriteString(Info(IconChip + " Windows"))
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

// IsTPMSupported returns true on Windows
func IsTPMSupported() bool {
	return true
}
