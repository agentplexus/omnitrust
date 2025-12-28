//go:build darwin

package inspector

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Security -framework IOKit

#import <Foundation/Foundation.h>
#import <Security/Security.h>
#include <sys/sysctl.h>

// Check if running on Apple Silicon
int tpm_isAppleSilicon() {
    int ret = 0;
    size_t size = sizeof(ret);
    if (sysctlbyname("hw.optional.arm64", &ret, &size, NULL, 0) == 0) {
        return ret;
    }
    return 0;
}

// Test if Secure Enclave is available by attempting to create an SE key
int tpm_testSecureEnclaveAvailable() {
    // On Apple Silicon, SE is always available
    if (tpm_isAppleSilicon()) {
        return 1;
    }

    // For Intel Macs, try to create an SE key to test availability
    SecAccessControlRef access = SecAccessControlCreateWithFlags(
        kCFAllocatorDefault,
        kSecAttrAccessibleWhenUnlockedThisDeviceOnly,
        kSecAccessControlPrivateKeyUsage,
        NULL
    );

    if (access == NULL) {
        return 0;
    }

    NSDictionary *attributes = @{
        (id)kSecAttrKeyType: (id)kSecAttrKeyTypeECSECPrimeRandom,
        (id)kSecAttrKeySizeInBits: @256,
        (id)kSecAttrTokenID: (id)kSecAttrTokenIDSecureEnclave,
        (id)kSecPrivateKeyAttrs: @{
            (id)kSecAttrIsPermanent: @NO,
            (id)kSecAttrAccessControl: (__bridge id)access,
        },
    };

    CFErrorRef error = NULL;
    SecKeyRef privateKey = SecKeyCreateRandomKey((__bridge CFDictionaryRef)attributes, &error);

    CFRelease(access);

    if (privateKey != NULL) {
        CFRelease(privateKey);
        return 1;
    }

    if (error != NULL) {
        CFIndex code = CFErrorGetCode(error);
        CFRelease(error);
        // errSecUnimplemented (-4) means SE not available
        if (code == -4) {
            return 0;
        }
        // Other errors (auth required, etc.) suggest SE is present
        return 1;
    }

    return 0;
}

// Get platform type string
const char* tpm_getPlatformType() {
    if (tpm_isAppleSilicon()) {
        return "apple_silicon";
    }
    return "intel";
}
*/
import "C"
import (
	"fmt"
	"strings"
)

// TPMResult contains TPM/Secure Enclave status information
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

// GetTPMStatus returns the TPM/Secure Enclave status (macOS)
func GetTPMStatus() (*TPMResult, error) {
	seAvailable := C.tpm_testSecureEnclaveAvailable() == 1
	platform := C.GoString(C.tpm_getPlatformType())

	var version, tpmType string
	var capabilities []string

	if platform == "apple_silicon" {
		version = "Secure Enclave (Apple Silicon)"
		tpmType = "secure_enclave"
		capabilities = []string{
			"hardware_key_generation",
			"hardware_key_storage",
			"biometric_authentication",
			"secure_boot",
			"encrypted_memory",
		}
	} else {
		version = "Secure Enclave (T2)"
		tpmType = "secure_enclave_t2"
		capabilities = []string{
			"hardware_key_generation",
			"hardware_key_storage",
			"biometric_authentication",
			"secure_boot",
		}
	}

	return &TPMResult{
		Present:            seAvailable,
		Enabled:            seAvailable,
		Version:            version,
		Manufacturer:       "Apple",
		Type:               tpmType,
		Platform:           platform,
		Capabilities:       capabilities,
		HardwareKeySupport: seAvailable,
	}, nil
}

// FormatTPMTable formats TPM status as a colored table
func FormatTPMTable(result *TPMResult) string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(Header(IconShield + " TPM / Secure Enclave Status"))
	sb.WriteString("\n")
	sb.WriteString(Muted(strings.Repeat("─", 55)))
	sb.WriteString("\n\n")

	// Platform badge
	var platformIcon string
	if result.Platform == "apple_silicon" {
		platformIcon = IconApple + " Apple Silicon"
	} else {
		platformIcon = IconChip + " Intel (T2)"
	}
	sb.WriteString(BoldText("Platform: "))
	sb.WriteString(Info(platformIcon))
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
		PadRight(IconShield+" TPM/SE Present", 28),
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

	// Hardware Key Support
	sb.WriteString(TableRowColored(
		PadRight(IconKey+" Hardware Key Support", 28),
		PadRight(BoolToStatusColored(result.HardwareKeySupport), 22),
	))
	sb.WriteString("\n")

	sb.WriteString(TableBottom(28, 22))
	sb.WriteString("\n\n")

	// Capabilities section
	sb.WriteString(BoldText("Capabilities:"))
	sb.WriteString("\n")
	sb.WriteString(Muted(strings.Repeat("─", 35)))
	sb.WriteString("\n")
	for _, cap := range result.Capabilities {
		sb.WriteString(fmt.Sprintf("  %s %s\n", Success(IconCheck), cap))
	}
	sb.WriteString("\n")

	return sb.String()
}

// FormatTPM formats TPM status in the specified format
func FormatTPM(result *TPMResult, format string) string {
	return FormatOutput(result, func() string {
		return FormatTPMTable(result)
	}, format)
}

// IsTPMSupported returns true on macOS (Secure Enclave serves as TPM)
func IsTPMSupported() bool {
	return true
}
