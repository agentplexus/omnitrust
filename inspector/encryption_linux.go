//go:build linux

package inspector

import (
	"os"
	"os/exec"
	"path/filepath"
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

// GetEncryptionStatus returns the disk encryption status (Linux - LUKS)
func GetEncryptionStatus() (*EncryptionResult, error) {
	result := &EncryptionResult{
		Platform: "linux",
		Type:     "luks",
	}

	var encryptedVolumes []EncryptedVolume

	// Check for dm-crypt/LUKS encrypted volumes
	// Look in /dev/mapper for crypt devices
	dmMapperPath := "/dev/mapper"
	entries, err := os.ReadDir(dmMapperPath)
	if err == nil {
		for _, entry := range entries {
			if entry.Name() == "control" {
				continue
			}

			// Check if this is a crypt device
			devicePath := filepath.Join(dmMapperPath, entry.Name())
			dmPath := filepath.Join("/sys/block", "dm-*", "dm/name")

			// Use dmsetup to check if it's a crypt target
			// #nosec G204 -- entry.Name() comes from trusted /dev/mapper directory listing
			out, err := exec.Command("dmsetup", "table", entry.Name()).Output()
			if err == nil && strings.Contains(string(out), "crypt") {
				vol := EncryptedVolume{
					Name:      entry.Name(),
					Encrypted: true,
					Status:    "encrypted_active",
				}

				// Try to find mount point
				mountOut, err := exec.Command("findmnt", "-n", "-o", "TARGET", devicePath).Output()
				if err == nil {
					vol.MountPoint = strings.TrimSpace(string(mountOut))
				}

				encryptedVolumes = append(encryptedVolumes, vol)
			}
			_ = dmPath // suppress unused warning
		}
	}

	// Also check /etc/crypttab for configured encrypted volumes
	crypttabData, err := os.ReadFile("/etc/crypttab")
	if err == nil {
		lines := strings.Split(string(crypttabData), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) >= 2 {
				name := fields[0]

				// Check if already in our list
				found := false
				for _, vol := range encryptedVolumes {
					if vol.Name == name {
						found = true
						break
					}
				}

				if !found {
					// Check if device exists in /dev/mapper
					_, err := os.Stat(filepath.Join(dmMapperPath, name))
					status := "configured_inactive"
					if err == nil {
						status = "configured_active"
					}

					encryptedVolumes = append(encryptedVolumes, EncryptedVolume{
						Name:      name,
						Encrypted: true,
						Status:    status,
					})
				}
			}
		}
	}

	// Check for LUKS headers on block devices
	blockDevices, _ := filepath.Glob("/dev/sd*")
	blockDevices2, _ := filepath.Glob("/dev/nvme*")
	blockDevices = append(blockDevices, blockDevices2...)

	for _, dev := range blockDevices {
		// Skip if it's a partition number > 9 to avoid too many checks
		if strings.HasSuffix(dev, "0") {
			continue
		}

		out, err := exec.Command("cryptsetup", "isLuks", dev).Output()
		_ = out
		if err == nil {
			// This is a LUKS device
			name := filepath.Base(dev)

			// Check if already mapped
			found := false
			for _, vol := range encryptedVolumes {
				if strings.Contains(vol.Name, name) {
					found = true
					break
				}
			}

			if !found {
				encryptedVolumes = append(encryptedVolumes, EncryptedVolume{
					Name:      name + " (LUKS)",
					Encrypted: true,
					Status:    "luks_device",
				})
			}
		}
	}

	result.EncryptedVolumes = encryptedVolumes

	if len(encryptedVolumes) > 0 {
		result.Enabled = true
		result.Status = "enabled"
		result.Details = "LUKS/dm-crypt encryption detected"
	} else {
		result.Enabled = false
		result.Status = "disabled"
		result.Details = "No LUKS/dm-crypt encrypted volumes detected"
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
	sb.WriteString(Info(IconChip + " Linux (LUKS/dm-crypt)"))
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
		PadRight(IconLock+" LUKS Encryption", 24),
		PadRight(BoolToStatusColored(result.Enabled), 26),
	))
	sb.WriteString("\n")

	// Status
	var statusDisplay string
	switch result.Status {
	case "enabled":
		statusDisplay = Success("Enabled")
	case "disabled":
		statusDisplay = Warning("Not Detected")
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
		sb.WriteString(BoldText("Encrypted Volumes:"))
		sb.WriteString("\n")
		sb.WriteString(Muted(strings.Repeat("─", 50)))
		sb.WriteString("\n")

		for _, vol := range result.EncryptedVolumes {
			statusStr := vol.Status
			switch vol.Status {
			case "encrypted_active", "configured_active":
				statusStr = Success("Active")
			case "configured_inactive":
				statusStr = Warning("Inactive")
			case "luks_device":
				statusStr = Info("LUKS Device")
			}

			sb.WriteString("  " + BoolToCheckbox(vol.Encrypted) + " ")
			sb.WriteString(vol.Name)
			if vol.MountPoint != "" {
				sb.WriteString(Muted(" -> " + vol.MountPoint))
			}
			sb.WriteString(" [" + statusStr + "]")
			sb.WriteString("\n")
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

// IsEncryptionSupported returns true on Linux
func IsEncryptionSupported() bool {
	return true
}
