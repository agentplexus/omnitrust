//go:build windows

package inspector

import (
	"fmt"
	"strings"

	"github.com/yusufpapurcu/wmi"
)

// Win32_EncryptableVolume represents WMI BitLocker class
type Win32_EncryptableVolume struct {
	DeviceID           string
	DriveLetter        string
	ProtectionStatus   uint32
	ConversionStatus   uint32
	EncryptionMethod   uint32
	VolumeType         uint32
}

// EncryptionResult contains disk encryption status information
type EncryptionResult struct {
	Enabled          bool              `json:"enabled"`
	Platform         string            `json:"platform"`
	Type             string            `json:"type"`
	Status           string            `json:"status"`
	EncryptedVolumes []EncryptedVolume `json:"encrypted_volumes,omitempty"`
	Details          string            `json:"details,omitempty"`
}

// EncryptedVolume represents an encrypted volume
type EncryptedVolume struct {
	Name       string `json:"name"`
	MountPoint string `json:"mount_point,omitempty"`
	Encrypted  bool   `json:"encrypted"`
	Status     string `json:"status"`
}

// GetEncryptionStatus returns the disk encryption status (Windows - BitLocker)
func GetEncryptionStatus() (*EncryptionResult, error) {
	result := &EncryptionResult{
		Platform: "windows",
		Type:     "bitlocker",
	}

	var volumes []Win32_EncryptableVolume

	// Query WMI for BitLocker status
	// Note: Requires running as administrator
	query := "SELECT * FROM Win32_EncryptableVolume"
	err := wmi.QueryNamespace(query, &volumes, `root\cimv2\Security\MicrosoftVolumeEncryption`)

	if err != nil || len(volumes) == 0 {
		// BitLocker not found or not accessible
		result.Status = "unknown"
		result.Details = "Unable to query BitLocker status (may require admin privileges)"
		return result, nil
	}

	var encryptedVolumes []EncryptedVolume
	anyEnabled := false

	for _, vol := range volumes {
		ev := EncryptedVolume{
			MountPoint: vol.DriveLetter,
			Name:       fmt.Sprintf("Volume %s", vol.DriveLetter),
		}

		// ProtectionStatus: 0 = OFF, 1 = ON, 2 = UNKNOWN
		if vol.ProtectionStatus == 1 {
			ev.Encrypted = true
			anyEnabled = true

			// ConversionStatus: 0 = FullyDecrypted, 1 = FullyEncrypted, 2 = EncryptionInProgress, etc.
			switch vol.ConversionStatus {
			case 1:
				ev.Status = "encrypted"
			case 2:
				ev.Status = "encrypting"
			case 3:
				ev.Status = "decrypting"
			case 4:
				ev.Status = "encryption_paused"
			case 5:
				ev.Status = "decryption_paused"
			default:
				ev.Status = "protected"
			}
		} else {
			ev.Encrypted = false
			ev.Status = "not_encrypted"
		}

		encryptedVolumes = append(encryptedVolumes, ev)
	}

	result.EncryptedVolumes = encryptedVolumes
	result.Enabled = anyEnabled

	if anyEnabled {
		result.Status = "enabled"
		result.Details = "BitLocker disk encryption is enabled on one or more volumes"
	} else {
		result.Status = "disabled"
		result.Details = "BitLocker disk encryption is not enabled on any volume"
	}

	return result, nil
}

// FormatEncryptionTable formats encryption status as a colored table
func FormatEncryptionTable(result *EncryptionResult) string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(Header(IconLock + " Disk Encryption Status"))
	sb.WriteString("\n")
	sb.WriteString(Muted(strings.Repeat("─", 55)))
	sb.WriteString("\n\n")

	// Platform badge
	sb.WriteString(BoldText("Platform: "))
	sb.WriteString(Info(IconChip + " Windows (BitLocker)"))
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
		PadRight(IconLock+" BitLocker Enabled", 24),
		PadRight(BoolToStatusColored(result.Enabled), 26),
	))
	sb.WriteString("\n")

	// Status
	statusDisplay := result.Status
	switch result.Status {
	case "enabled":
		statusDisplay = Success("Enabled")
	case "disabled":
		statusDisplay = Danger("Disabled")
	default:
		statusDisplay = Muted(result.Status)
	}
	sb.WriteString(TableRowColored(
		PadRight(IconStatus+" Status", 24),
		PadRight(statusDisplay, 26),
	))
	sb.WriteString("\n")

	sb.WriteString(TableBottom(24, 26))
	sb.WriteString("\n")

	// Encrypted volumes
	if len(result.EncryptedVolumes) > 0 {
		sb.WriteString("\n")
		sb.WriteString(BoldText("Volumes:"))
		sb.WriteString("\n")
		sb.WriteString(Muted(strings.Repeat("─", 50)))
		sb.WriteString("\n")

		sb.WriteString(TableTop(10, 18, 18))
		sb.WriteString("\n")
		sb.WriteString(TableRowColored(
			Header(PadRight("Drive", 10)),
			Header(PadRight("Encrypted", 18)),
			Header(PadRight("Status", 18)),
		))
		sb.WriteString("\n")
		sb.WriteString(TableSeparator(10, 18, 18))
		sb.WriteString("\n")

		for _, vol := range result.EncryptedVolumes {
			statusStr := vol.Status
			switch vol.Status {
			case "encrypted", "protected":
				statusStr = Success("Encrypted")
			case "encrypting":
				statusStr = Warning("Encrypting...")
			case "decrypting":
				statusStr = Warning("Decrypting...")
			case "not_encrypted":
				statusStr = Danger("Not Encrypted")
			}

			sb.WriteString(TableRowColored(
				PadRight(vol.MountPoint, 10),
				PadRight(BoolToStatusColored(vol.Encrypted), 18),
				PadRight(statusStr, 18),
			))
			sb.WriteString("\n")
		}
		sb.WriteString(TableBottom(10, 18, 18))
		sb.WriteString("\n")
	}

	// Details if available
	if result.Details != "" {
		sb.WriteString("\n")
		sb.WriteString(Muted("Details: " + result.Details))
		sb.WriteString("\n")
	}

	return sb.String()
}

// FormatEncryption formats encryption status in the specified format
func FormatEncryption(result *EncryptionResult, format string) string {
	return FormatOutput(result, func() string {
		return FormatEncryptionTable(result)
	}, format)
}

// IsEncryptionSupported returns true on Windows
func IsEncryptionSupported() bool {
	return true
}
