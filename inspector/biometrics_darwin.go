//go:build darwin

package inspector

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework LocalAuthentication

#import <Foundation/Foundation.h>
#import <LocalAuthentication/LocalAuthentication.h>

// Biometric info structure
typedef struct {
    int touchIDAvailable;
    int touchIDEnrolled;
    int faceIDAvailable;
    int faceIDEnrolled;
    int biometryType;
} BiometricInfo;

BiometricInfo getBiometricInfo() {
    BiometricInfo info = {0, 0, 0, 0, 0};

    LAContext *context = [[LAContext alloc] init];
    NSError *error = nil;

    BOOL canEvaluate = [context canEvaluatePolicy:LAPolicyDeviceOwnerAuthenticationWithBiometrics error:&error];

    if (canEvaluate) {
        info.biometryType = (int)context.biometryType;
        if (context.biometryType == LABiometryTypeTouchID) {
            info.touchIDAvailable = 1;
            info.touchIDEnrolled = 1;
        } else if (context.biometryType == LABiometryTypeFaceID) {
            info.faceIDAvailable = 1;
            info.faceIDEnrolled = 1;
        }
    } else if (error) {
        info.biometryType = (int)context.biometryType;
        if (error.code == LAErrorBiometryNotEnrolled) {
            if (context.biometryType == LABiometryTypeTouchID) {
                info.touchIDAvailable = 1;
            } else if (context.biometryType == LABiometryTypeFaceID) {
                info.faceIDAvailable = 1;
            }
        }
    }

    return info;
}
*/
import "C"
import (
	"strings"
)

// BiometricCapabilities contains detailed biometric capability information
type BiometricCapabilities struct {
	TouchIDAvailable bool   `json:"touch_id_available"`
	TouchIDEnrolled  bool   `json:"touch_id_enrolled"`
	FaceIDAvailable  bool   `json:"face_id_available"`
	FaceIDEnrolled   bool   `json:"face_id_enrolled"`
	BiometryType     string `json:"biometry_type"`
}

// GetBiometricCapabilities returns detailed biometric capabilities (macOS only)
func GetBiometricCapabilities() (*BiometricCapabilities, error) {
	bioInfo := C.getBiometricInfo()

	biometryType := "none"
	switch bioInfo.biometryType {
	case 1:
		biometryType = "touch_id"
	case 2:
		biometryType = "face_id"
	}

	return &BiometricCapabilities{
		TouchIDAvailable: bioInfo.touchIDAvailable == 1,
		TouchIDEnrolled:  bioInfo.touchIDEnrolled == 1,
		FaceIDAvailable:  bioInfo.faceIDAvailable == 1,
		FaceIDEnrolled:   bioInfo.faceIDEnrolled == 1,
		BiometryType:     biometryType,
	}, nil
}

// FormatBiometricCapabilitiesTable formats biometric capabilities as a colored table
func FormatBiometricCapabilitiesTable(result *BiometricCapabilities) string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(Header(IconFingerprint + " Biometric Capabilities"))
	sb.WriteString("\n")
	sb.WriteString(Muted(strings.Repeat("─", 50)))
	sb.WriteString("\n\n")

	// Active biometry type
	sb.WriteString(BoldText("Active Biometry: "))
	switch result.BiometryType {
	case "touch_id":
		sb.WriteString(Success(IconFingerprint + " Touch ID"))
	case "face_id":
		sb.WriteString(Success(IconFace + " Face ID"))
	default:
		sb.WriteString(Muted("None"))
	}
	sb.WriteString("\n\n")

	// Capabilities table
	sb.WriteString(TableTop(14, 14, 14))
	sb.WriteString("\n")
	sb.WriteString(TableRowColored(
		Header(PadRight("Biometric", 14)),
		Header(PadRight("Available", 14)),
		Header(PadRight("Enrolled", 14)),
	))
	sb.WriteString("\n")
	sb.WriteString(TableSeparator(14, 14, 14))
	sb.WriteString("\n")

	// Touch ID row
	sb.WriteString(TableRowColored(
		PadRight(IconFingerprint+" Touch ID", 14),
		PadRight(BoolToStatusColored(result.TouchIDAvailable), 14),
		PadRight(BoolToStatusColored(result.TouchIDEnrolled), 14),
	))
	sb.WriteString("\n")

	// Face ID row
	sb.WriteString(TableRowColored(
		PadRight(IconFace+" Face ID", 14),
		PadRight(BoolToStatusColored(result.FaceIDAvailable), 14),
		PadRight(BoolToStatusColored(result.FaceIDEnrolled), 14),
	))
	sb.WriteString("\n")

	sb.WriteString(TableBottom(14, 14, 14))
	sb.WriteString("\n")

	return sb.String()
}

// FormatBiometricCapabilities formats biometric capabilities in the specified format
func FormatBiometricCapabilities(result *BiometricCapabilities, format string) string {
	return FormatOutput(result, func() string {
		return FormatBiometricCapabilitiesTable(result)
	}, format)
}

// IsBiometricsSupported returns true on macOS
func IsBiometricsSupported() bool {
	return true
}
