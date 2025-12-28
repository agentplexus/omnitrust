//go:build darwin

package inspector

import (
	"os/exec"
	"strings"
)

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

// GetEncryptionStatus returns the disk encryption status (macOS - FileVault)
func GetEncryptionStatus() (*EncryptionResult, error) {
	result := &EncryptionResult{
		Platform: "darwin",
		Type:     "filevault",
	}

	// Check FileVault status using fdesetup
	out, err := exec.Command("fdesetup", "status").Output()
	if err != nil {
		// fdesetup might require admin privileges
		result.Status = "unknown"
		result.Details = "Unable to determine FileVault status (may require admin)"
		return result, nil
	}

	output := strings.TrimSpace(string(out))

	if strings.Contains(output, "FileVault is On") {
		result.Enabled = true
		result.Status = "enabled"
		result.Details = "FileVault disk encryption is enabled"

		// Check for encryption in progress
		if strings.Contains(output, "Encryption in progress") {
			result.Status = "encrypting"
			result.Details = "FileVault encryption in progress"
		} else if strings.Contains(output, "Decryption in progress") {
			result.Status = "decrypting"
			result.Details = "FileVault decryption in progress"
		}
	} else if strings.Contains(output, "FileVault is Off") {
		result.Enabled = false
		result.Status = "disabled"
		result.Details = "FileVault disk encryption is disabled"
	} else {
		result.Status = "unknown"
		result.Details = output
	}

	// Get list of encrypted volumes using diskutil
	volumes := getEncryptedVolumes()
	result.EncryptedVolumes = volumes

	return result, nil
}

// getEncryptedVolumes returns a list of APFS encrypted volumes
func getEncryptedVolumes() []EncryptedVolume {
	var volumes []EncryptedVolume

	// Use diskutil to list APFS containers and check encryption
	out, err := exec.Command("diskutil", "apfs", "list", "-plist").Output()
	if err != nil {
		// Fallback: check just the root volume
		out, err := exec.Command("diskutil", "info", "/").Output()
		if err == nil {
			output := string(out)
			vol := EncryptedVolume{
				Name:       "Macintosh HD",
				MountPoint: "/",
			}
			if strings.Contains(output, "FileVault:") && strings.Contains(output, "Yes") {
				vol.Encrypted = true
				vol.Status = "encrypted"
			} else if strings.Contains(output, "Encrypted:") && strings.Contains(output, "Yes") {
				vol.Encrypted = true
				vol.Status = "encrypted"
			} else {
				vol.Encrypted = false
				vol.Status = "not_encrypted"
			}
			volumes = append(volumes, vol)
		}
		return volumes
	}

	// Parse diskutil output to find encrypted volumes
	// For simplicity, we'll check common volumes
	output := string(out)
	if strings.Contains(output, "Encryption") || strings.Contains(output, "FileVault") {
		// Check root volume
		rootOut, err := exec.Command("diskutil", "info", "/").Output()
		if err == nil {
			rootOutput := string(rootOut)
			vol := EncryptedVolume{
				MountPoint: "/",
			}

			// Get volume name
			for _, line := range strings.Split(rootOutput, "\n") {
				if strings.Contains(line, "Volume Name:") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) == 2 {
						vol.Name = strings.TrimSpace(parts[1])
					}
				}
			}
			if vol.Name == "" {
				vol.Name = "System Volume"
			}

			if strings.Contains(rootOutput, "Yes (Unlocked)") ||
				(strings.Contains(rootOutput, "Encrypted:") && strings.Contains(rootOutput, "Yes")) {
				vol.Encrypted = true
				vol.Status = "encrypted_unlocked"
			} else {
				vol.Encrypted = false
				vol.Status = "not_encrypted"
			}
			volumes = append(volumes, vol)
		}
	}

	return volumes
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
	sb.WriteString(Info(IconApple + " macOS (FileVault)"))
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
		PadRight(IconLock+" FileVault Enabled", 24),
		PadRight(BoolToStatusColored(result.Enabled), 26),
	))
	sb.WriteString("\n")

	// Status
	var statusDisplay string
	switch result.Status {
	case "enabled":
		statusDisplay = Success("Enabled")
	case "disabled":
		statusDisplay = Danger("Disabled")
	case "encrypting":
		statusDisplay = Warning("Encrypting...")
	case "decrypting":
		statusDisplay = Warning("Decrypting...")
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
		sb.WriteString(Muted(strings.Repeat("─", 40)))
		sb.WriteString("\n")
		for _, vol := range result.EncryptedVolumes {
			icon := IconCross
			if vol.Encrypted {
				icon = IconCheck
			}
			statusStr := Danger("Not Encrypted")
			if vol.Encrypted {
				statusStr = Success("Encrypted")
			}
			sb.WriteString("  " + BoolToCheckbox(vol.Encrypted) + " ")
			sb.WriteString(vol.Name)
			if vol.MountPoint != "" {
				sb.WriteString(Muted(" (" + vol.MountPoint + ")"))
			}
			sb.WriteString(" - " + statusStr)
			sb.WriteString("\n")
			_ = icon // suppress unused warning
		}
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

// IsEncryptionSupported returns true on macOS
func IsEncryptionSupported() bool {
	return true
}
