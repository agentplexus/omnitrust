package inspector

import (
	"fmt"
	"runtime"
	"strings"
)

// SecuritySummary contains a unified security posture overview
type SecuritySummary struct {
	Platform        string       `json:"platform"`
	OverallScore    int          `json:"overall_score"`
	OverallStatus   string       `json:"overall_status"`
	TPM             *TPMSummary  `json:"tpm"`
	SecureBoot      *BootSummary `json:"secure_boot"`
	Encryption      *EncSummary  `json:"encryption"`
	Biometrics      *BioSummary  `json:"biometrics"`
	Recommendations []string     `json:"recommendations,omitempty"`
}

// TPMSummary contains TPM summary info
type TPMSummary struct {
	Present bool   `json:"present"`
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
}

// BootSummary contains Secure Boot summary info
type BootSummary struct {
	Enabled bool   `json:"enabled"`
	Mode    string `json:"mode"`
}

// EncSummary contains encryption summary info
type EncSummary struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
	Status  string `json:"status"`
}

// BioSummary contains biometrics summary info
type BioSummary struct {
	Available  bool   `json:"available"`
	Configured bool   `json:"configured"`
	Type       string `json:"type"`
}

// GetSecuritySummary returns a unified security posture overview
func GetSecuritySummary() (*SecuritySummary, error) {
	summary := &SecuritySummary{
		Platform: runtime.GOOS,
	}

	var score int
	var recommendations []string

	// Get TPM status
	if IsTPMSupported() {
		tpmResult, err := GetTPMStatus()
		if err == nil {
			summary.TPM = &TPMSummary{
				Present: tpmResult.Present,
				Enabled: tpmResult.Enabled,
				Type:    tpmResult.Type,
			}
			if tpmResult.Present && tpmResult.Enabled {
				score += 25
			} else if !tpmResult.Present {
				recommendations = append(recommendations, "Hardware security module (TPM/Secure Enclave) not detected")
			}
		}
	}

	// Get Secure Boot status
	if IsSecureBootSupported() {
		bootResult, err := GetSecureBootStatus()
		if err == nil {
			summary.SecureBoot = &BootSummary{
				Enabled: bootResult.Enabled,
				Mode:    bootResult.Mode,
			}
			if bootResult.Enabled {
				score += 25
			} else {
				recommendations = append(recommendations, "Enable Secure Boot for enhanced boot security")
			}
		}
	}

	// Get Encryption status
	if IsEncryptionSupported() {
		encResult, err := GetEncryptionStatus()
		if err == nil {
			summary.Encryption = &EncSummary{
				Enabled: encResult.Enabled,
				Type:    encResult.Type,
				Status:  encResult.Status,
			}
			if encResult.Enabled {
				score += 25
			} else {
				encType := "disk encryption"
				switch runtime.GOOS {
				case "darwin":
					encType = "FileVault"
				case "windows":
					encType = "BitLocker"
				case "linux":
					encType = "LUKS"
				}
				recommendations = append(recommendations, fmt.Sprintf("Enable %s to protect data at rest", encType))
			}
		}
	}

	// Get Biometrics status
	if IsBiometricsSupported() {
		bioResult, err := GetBiometricCapabilities()
		if err == nil {
			available := bioResult.TouchIDAvailable || bioResult.FaceIDAvailable
			configured := bioResult.TouchIDEnrolled || bioResult.FaceIDEnrolled
			summary.Biometrics = &BioSummary{
				Available:  available,
				Configured: configured,
				Type:       bioResult.BiometryType,
			}
			if configured {
				score += 25
			} else if available {
				recommendations = append(recommendations, "Configure biometric authentication for enhanced security")
			}
		}
	}

	summary.OverallScore = score
	summary.Recommendations = recommendations

	// Determine overall status
	switch {
	case score >= 100:
		summary.OverallStatus = "excellent"
	case score >= 75:
		summary.OverallStatus = "good"
	case score >= 50:
		summary.OverallStatus = "fair"
	case score >= 25:
		summary.OverallStatus = "needs_improvement"
	default:
		summary.OverallStatus = "critical"
	}

	return summary, nil
}

// FormatSecuritySummaryTable formats security summary as a colored table
func FormatSecuritySummaryTable(result *SecuritySummary) string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(Header(IconShield + " Security Summary"))
	sb.WriteString("\n")
	sb.WriteString(Muted(strings.Repeat("─", 60)))
	sb.WriteString("\n\n")

	// Platform
	platformName := result.Platform
	platformIcon := IconChip
	switch result.Platform {
	case "darwin":
		platformName = "macOS"
		platformIcon = IconApple
	case "windows":
		platformName = "Windows"
	case "linux":
		platformName = "Linux"
	}
	sb.WriteString(BoldText("Platform: "))
	sb.WriteString(Info(platformIcon + " " + platformName))
	sb.WriteString("\n\n")

	// Overall Score with visual bar
	sb.WriteString(BoldText("Security Score: "))
	scoreColor := UsageColor(float64(100 - result.OverallScore)) // Invert for security (higher is better)
	sb.WriteString(Colorize(scoreColor+Bold, fmt.Sprintf("%d/100", result.OverallScore)))
	sb.WriteString("\n")
	sb.WriteString(securityScoreBar(result.OverallScore, 40))
	sb.WriteString("\n\n")

	// Overall Status Badge
	sb.WriteString(BoldText("Status: "))
	switch result.OverallStatus {
	case "excellent":
		sb.WriteString(Success(IconCheck + " Excellent"))
	case "good":
		sb.WriteString(Success(IconCheck + " Good"))
	case "fair":
		sb.WriteString(Warning(IconWarning + " Fair"))
	case "needs_improvement":
		sb.WriteString(Warning(IconWarning + " Needs Improvement"))
	case "critical":
		sb.WriteString(Danger(IconCross + " Critical"))
	}
	sb.WriteString("\n\n")

	// Security Features Table
	sb.WriteString(BoldText("Security Features:"))
	sb.WriteString("\n")
	sb.WriteString(TableTop(24, 12, 18))
	sb.WriteString("\n")
	sb.WriteString(TableRowColored(
		Header(PadRight("Feature", 24)),
		Header(PadRight("Status", 12)),
		Header(PadRight("Details", 18)),
	))
	sb.WriteString("\n")
	sb.WriteString(TableSeparator(24, 12, 18))
	sb.WriteString("\n")

	// TPM / Secure Enclave
	var tpmName string
	switch result.Platform {
	case "darwin":
		tpmName = "Secure Enclave"
	default:
		tpmName = "TPM"
	}
	if result.TPM != nil {
		sb.WriteString(TableRowColored(
			PadRight(IconShield+" "+tpmName, 24),
			PadRight(featureStatus(result.TPM.Present && result.TPM.Enabled), 12),
			PadRight(result.TPM.Type, 18),
		))
	} else {
		sb.WriteString(TableRowColored(
			PadRight(IconShield+" "+tpmName, 24),
			PadRight(Muted("N/A"), 12),
			PadRight(Muted("-"), 18),
		))
	}
	sb.WriteString("\n")

	// Secure Boot
	if result.SecureBoot != nil {
		sb.WriteString(TableRowColored(
			PadRight(IconLock+" Secure Boot", 24),
			PadRight(featureStatus(result.SecureBoot.Enabled), 12),
			PadRight(result.SecureBoot.Mode, 18),
		))
	} else {
		sb.WriteString(TableRowColored(
			PadRight(IconLock+" Secure Boot", 24),
			PadRight(Muted("N/A"), 12),
			PadRight(Muted("-"), 18),
		))
	}
	sb.WriteString("\n")

	// Disk Encryption
	var encName string
	switch result.Platform {
	case "darwin":
		encName = "FileVault"
	case "windows":
		encName = "BitLocker"
	case "linux":
		encName = "LUKS"
	default:
		encName = "Disk Encryption"
	}
	if result.Encryption != nil {
		sb.WriteString(TableRowColored(
			PadRight(IconLock+" "+encName, 24),
			PadRight(featureStatus(result.Encryption.Enabled), 12),
			PadRight(result.Encryption.Status, 18),
		))
	} else {
		sb.WriteString(TableRowColored(
			PadRight(IconLock+" "+encName, 24),
			PadRight(Muted("N/A"), 12),
			PadRight(Muted("-"), 18),
		))
	}
	sb.WriteString("\n")

	// Biometrics
	if result.Biometrics != nil {
		sb.WriteString(TableRowColored(
			PadRight(IconFingerprint+" Biometrics", 24),
			PadRight(featureStatus(result.Biometrics.Configured), 12),
			PadRight(result.Biometrics.Type, 18),
		))
	} else {
		sb.WriteString(TableRowColored(
			PadRight(IconFingerprint+" Biometrics", 24),
			PadRight(Muted("N/A"), 12),
			PadRight(Muted("-"), 18),
		))
	}
	sb.WriteString("\n")

	sb.WriteString(TableBottom(24, 12, 18))
	sb.WriteString("\n")

	// Recommendations
	if len(result.Recommendations) > 0 {
		sb.WriteString("\n")
		sb.WriteString(BoldText(IconWarning + " Recommendations:"))
		sb.WriteString("\n")
		sb.WriteString(Muted(strings.Repeat("─", 50)))
		sb.WriteString("\n")
		for i, rec := range result.Recommendations {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, Warning(rec)))
		}
	}
	sb.WriteString("\n")

	return sb.String()
}

// securityScoreBar creates a security score progress bar (green = good)
func securityScoreBar(score int, width int) string {
	filled := score * width / 100
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	var color string
	switch {
	case score >= 75:
		color = Green
	case score >= 50:
		color = Yellow
	default:
		color = Red
	}

	bar := color + strings.Repeat(IconBar, filled) + Reset
	bar += Muted(strings.Repeat(IconBarLight, width-filled))
	return bar
}

// featureStatus returns a colored status indicator
func featureStatus(enabled bool) string {
	if enabled {
		return Success(IconCheck + " Enabled")
	}
	return Danger(IconCross + " Disabled")
}

// FormatSecuritySummary formats security summary in the specified format
func FormatSecuritySummary(result *SecuritySummary, format string) string {
	return FormatOutput(result, func() string {
		return FormatSecuritySummaryTable(result)
	}, format)
}
