//go:build !darwin && !windows && !linux

package inspector

import "errors"

// BiometricCapabilities contains detailed biometric capability information
type BiometricCapabilities struct {
	TouchIDAvailable bool   `json:"touch_id_available"`
	TouchIDEnrolled  bool   `json:"touch_id_enrolled"`
	FaceIDAvailable  bool   `json:"face_id_available"`
	FaceIDEnrolled   bool   `json:"face_id_enrolled"`
	BiometryType     string `json:"biometry_type"`
	Platform         string `json:"platform"`
}

// GetBiometricCapabilities returns an error on unsupported platforms
func GetBiometricCapabilities() (*BiometricCapabilities, error) {
	return nil, errors.New("biometric capabilities are not available on this platform")
}

// FormatBiometricCapabilitiesTable is not available on unsupported platforms
func FormatBiometricCapabilitiesTable(result *BiometricCapabilities) string {
	return "Biometric capabilities are not available on this platform"
}

// FormatBiometricCapabilities is not available on unsupported platforms
func FormatBiometricCapabilities(result *BiometricCapabilities, format string) string {
	return "Biometric capabilities are not available on this platform"
}

// IsBiometricsSupported returns false on unsupported platforms
func IsBiometricsSupported() bool {
	return false
}
