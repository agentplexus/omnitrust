//go:build windows

package inspector

import (
	"strings"
	"syscall"
	"unsafe"
)

var (
	credui                    = syscall.NewLazyDLL("credui.dll")
	procCredUIPromptForWindowsCredentialsW = credui.NewProc("CredUIPromptForWindowsCredentialsW")

	// For checking Windows Hello
	userenv = syscall.NewLazyDLL("userenv.dll")
)

// BiometricCapabilities contains detailed biometric capability information
type BiometricCapabilities struct {
	TouchIDAvailable bool   `json:"touch_id_available"`
	TouchIDEnrolled  bool   `json:"touch_id_enrolled"`
	FaceIDAvailable  bool   `json:"face_id_available"`
	FaceIDEnrolled   bool   `json:"face_id_enrolled"`
	BiometryType     string `json:"biometry_type"`
	// Windows-specific fields
	WindowsHelloAvailable  bool   `json:"windows_hello_available,omitempty"`
	WindowsHelloConfigured bool   `json:"windows_hello_configured,omitempty"`
	FingerprintAvailable   bool   `json:"fingerprint_available,omitempty"`
	FingerprintEnrolled    bool   `json:"fingerprint_enrolled,omitempty"`
	FacialRecognition      bool   `json:"facial_recognition,omitempty"`
	PINConfigured          bool   `json:"pin_configured,omitempty"`
	Platform               string `json:"platform"`
}

// GetBiometricCapabilities returns biometric capabilities (Windows)
func GetBiometricCapabilities() (*BiometricCapabilities, error) {
	result := &BiometricCapabilities{
		Platform:     "windows",
		BiometryType: "none",
	}

	// Check Windows Hello availability
	// This is a simplified check - in production you'd use Windows.Security.Credentials.UI
	// or WMI queries for more detailed information

	// Check if Windows Hello is available via registry or system capabilities
	// For now, we'll check for biometric devices

	// Try to detect fingerprint reader via WMI
	fingerprintAvailable := checkFingerprintSensor()
	result.FingerprintAvailable = fingerprintAvailable
	result.TouchIDAvailable = fingerprintAvailable // Map to TouchID equivalent

	// Check for Windows Hello face recognition (IR camera)
	faceAvailable := checkFaceRecognition()
	result.FacialRecognition = faceAvailable
	result.FaceIDAvailable = faceAvailable

	// Check if Windows Hello is configured
	helloConfigured := checkWindowsHelloConfigured()
	result.WindowsHelloConfigured = helloConfigured
	result.WindowsHelloAvailable = fingerprintAvailable || faceAvailable

	// Determine biometry type
	if fingerprintAvailable && faceAvailable {
		result.BiometryType = "fingerprint_and_face"
	} else if fingerprintAvailable {
		result.BiometryType = "fingerprint"
		result.TouchIDEnrolled = helloConfigured
	} else if faceAvailable {
		result.BiometryType = "face"
		result.FaceIDEnrolled = helloConfigured
	}

	return result, nil
}

// checkFingerprintSensor checks for fingerprint sensor availability
func checkFingerprintSensor() bool {
	// Check for biometric devices in the system
	// This is a simplified check - actual implementation would query WMI or use Windows Biometric Framework

	// Try to load the Windows Biometric Framework DLL
	winbio := syscall.NewLazyDLL("winbio.dll")
	if winbio.Load() == nil {
		// DLL loaded successfully, biometric framework is available
		// In a full implementation, you would call WinBioEnumBiometricUnits
		return true
	}
	return false
}

// checkFaceRecognition checks for Windows Hello face recognition availability
func checkFaceRecognition() bool {
	// Check for IR camera / Windows Hello face recognition
	// This would typically query the camera capabilities
	// Simplified implementation
	return false
}

// checkWindowsHelloConfigured checks if Windows Hello is set up for the current user
func checkWindowsHelloConfigured() bool {
	// Check NGC (Next Generation Credential) container
	// This indicates if Windows Hello is configured
	// Simplified check via registry or credential APIs

	// In production, you would check:
	// HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\Authentication\LogonUI\NgcPin
	// or use Windows.Security.Credentials APIs

	return false
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
	sb.WriteString(Info(IconChip + " Windows (Windows Hello)"))
	sb.WriteString("\n\n")

	// Windows Hello status
	sb.WriteString(BoldText("Windows Hello: "))
	if result.WindowsHelloAvailable {
		if result.WindowsHelloConfigured {
			sb.WriteString(Success("Available & Configured"))
		} else {
			sb.WriteString(Warning("Available (Not Configured)"))
		}
	} else {
		sb.WriteString(Muted("Not Available"))
	}
	sb.WriteString("\n\n")

	// Capabilities table
	sb.WriteString(TableTop(20, 14, 14))
	sb.WriteString("\n")
	sb.WriteString(TableRowColored(
		Header(PadRight("Biometric", 20)),
		Header(PadRight("Available", 14)),
		Header(PadRight("Enrolled", 14)),
	))
	sb.WriteString("\n")
	sb.WriteString(TableSeparator(20, 14, 14))
	sb.WriteString("\n")

	// Fingerprint row
	sb.WriteString(TableRowColored(
		PadRight(IconFingerprint+" Fingerprint", 20),
		PadRight(BoolToStatusColored(result.FingerprintAvailable), 14),
		PadRight(BoolToStatusColored(result.TouchIDEnrolled), 14),
	))
	sb.WriteString("\n")

	// Face Recognition row
	sb.WriteString(TableRowColored(
		PadRight(IconFace+" Face Recognition", 20),
		PadRight(BoolToStatusColored(result.FacialRecognition), 14),
		PadRight(BoolToStatusColored(result.FaceIDEnrolled), 14),
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

// IsBiometricsSupported returns true on Windows
func IsBiometricsSupported() bool {
	return true
}

// Suppress unused variable warning
var _ = unsafe.Sizeof(0)
