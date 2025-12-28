//go:build linux

package inspector

import (
	"os"
	"os/exec"
	"strings"
)

// BiometricCapabilities contains detailed biometric capability information
type BiometricCapabilities struct {
	TouchIDAvailable bool   `json:"touch_id_available"`
	TouchIDEnrolled  bool   `json:"touch_id_enrolled"`
	FaceIDAvailable  bool   `json:"face_id_available"`
	FaceIDEnrolled   bool   `json:"face_id_enrolled"`
	BiometryType     string `json:"biometry_type"`
	// Linux-specific fields
	FprintdAvailable bool   `json:"fprintd_available,omitempty"`
	FprintdEnrolled  bool   `json:"fprintd_enrolled,omitempty"`
	HowdyAvailable   bool   `json:"howdy_available,omitempty"`
	HowdyConfigured  bool   `json:"howdy_configured,omitempty"`
	Platform         string `json:"platform"`
}

// GetBiometricCapabilities returns biometric capabilities (Linux)
func GetBiometricCapabilities() (*BiometricCapabilities, error) {
	result := &BiometricCapabilities{
		Platform:     "linux",
		BiometryType: "none",
	}

	// Check for fprintd (fingerprint daemon)
	if _, err := exec.LookPath("fprintd-list"); err == nil {
		result.FprintdAvailable = true
		result.TouchIDAvailable = true

		// Check if fingerprints are enrolled
		// #nosec G204 -- USER env var is trusted system input for current user
		out, err := exec.Command("fprintd-list", os.Getenv("USER")).Output()
		if err == nil && strings.Contains(string(out), "fingerprint") {
			result.FprintdEnrolled = true
			result.TouchIDEnrolled = true
		}
	}

	// Check for Howdy (face recognition for Linux)
	if _, err := exec.LookPath("howdy"); err == nil {
		result.HowdyAvailable = true
		result.FaceIDAvailable = true

		// Check if face is configured
		out, err := exec.Command("howdy", "list").Output()
		if err == nil && !strings.Contains(string(out), "No face models") {
			result.HowdyConfigured = true
			result.FaceIDEnrolled = true
		}
	}

	// Determine biometry type
	if result.FprintdAvailable && result.HowdyAvailable {
		result.BiometryType = "fingerprint_and_face"
	} else if result.FprintdAvailable {
		result.BiometryType = "fingerprint"
	} else if result.HowdyAvailable {
		result.BiometryType = "face"
	}

	return result, nil
}

// FormatBiometricCapabilitiesTable formats biometric capabilities as a colored table
func FormatBiometricCapabilitiesTable(result *BiometricCapabilities) string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(Header(IconFingerprint + " Biometric Capabilities"))
	sb.WriteString("\n")
	sb.WriteString(Muted(strings.Repeat("─", 55)))
	sb.WriteString("\n\n")

	// Platform badge
	sb.WriteString(BoldText("Platform: "))
	sb.WriteString(Info(IconChip + " Linux"))
	sb.WriteString("\n\n")

	// Check if any biometrics available
	if !result.FprintdAvailable && !result.HowdyAvailable {
		sb.WriteString(Muted("No biometric authentication services detected."))
		sb.WriteString("\n")
		sb.WriteString(Muted("Consider installing:"))
		sb.WriteString("\n")
		sb.WriteString(Muted("  - fprintd: for fingerprint authentication"))
		sb.WriteString("\n")
		sb.WriteString(Muted("  - howdy: for facial recognition"))
		sb.WriteString("\n")
		return sb.String()
	}

	// Capabilities table
	sb.WriteString(TableTop(20, 14, 14))
	sb.WriteString("\n")
	sb.WriteString(TableRowColored(
		Header(PadRight("Service", 20)),
		Header(PadRight("Available", 14)),
		Header(PadRight("Configured", 14)),
	))
	sb.WriteString("\n")
	sb.WriteString(TableSeparator(20, 14, 14))
	sb.WriteString("\n")

	// fprintd row
	sb.WriteString(TableRowColored(
		PadRight(IconFingerprint+" fprintd", 20),
		PadRight(BoolToStatusColored(result.FprintdAvailable), 14),
		PadRight(BoolToStatusColored(result.FprintdEnrolled), 14),
	))
	sb.WriteString("\n")

	// Howdy row
	sb.WriteString(TableRowColored(
		PadRight(IconFace+" Howdy", 20),
		PadRight(BoolToStatusColored(result.HowdyAvailable), 14),
		PadRight(BoolToStatusColored(result.HowdyConfigured), 14),
	))
	sb.WriteString("\n")

	sb.WriteString(TableBottom(20, 14, 14))
	sb.WriteString("\n")

	return sb.String()
}

// FormatBiometricCapabilities formats biometric capabilities in the specified format
func FormatBiometricCapabilities(result *BiometricCapabilities, format string) string {
	return FormatOutput(result, func() string {
		return FormatBiometricCapabilitiesTable(result)
	}, format)
}

// IsBiometricsSupported returns true on Linux (fprintd/howdy may be available)
func IsBiometricsSupported() bool {
	return true
}
