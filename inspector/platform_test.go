package inspector

import (
	"runtime"
	"testing"
)

func TestIsTPMSupported(t *testing.T) {
	// IsTPMSupported should return a boolean without panicking
	supported := IsTPMSupported()
	t.Logf("IsTPMSupported() = %v (platform: %s)", supported, runtime.GOOS)

	// On known platforms, it should return true
	switch runtime.GOOS {
	case "darwin", "windows", "linux":
		if !supported {
			t.Logf("Warning: TPM support expected on %s but returned false", runtime.GOOS)
		}
	}
}

func TestIsSecureBootSupported(t *testing.T) {
	supported := IsSecureBootSupported()
	t.Logf("IsSecureBootSupported() = %v (platform: %s)", supported, runtime.GOOS)

	switch runtime.GOOS {
	case "darwin", "windows", "linux":
		if !supported {
			t.Logf("Warning: Secure Boot support expected on %s but returned false", runtime.GOOS)
		}
	}
}

func TestIsEncryptionSupported(t *testing.T) {
	supported := IsEncryptionSupported()
	t.Logf("IsEncryptionSupported() = %v (platform: %s)", supported, runtime.GOOS)

	switch runtime.GOOS {
	case "darwin", "windows", "linux":
		if !supported {
			t.Logf("Warning: Encryption support expected on %s but returned false", runtime.GOOS)
		}
	}
}

func TestIsBiometricsSupported(t *testing.T) {
	supported := IsBiometricsSupported()
	t.Logf("IsBiometricsSupported() = %v (platform: %s)", supported, runtime.GOOS)

	switch runtime.GOOS {
	case "darwin", "windows", "linux":
		if !supported {
			t.Logf("Warning: Biometrics support expected on %s but returned false", runtime.GOOS)
		}
	}
}

func TestGetTPMStatus_WhenSupported(t *testing.T) {
	if !IsTPMSupported() {
		t.Skip("TPM not supported on this platform")
	}

	result, err := GetTPMStatus()
	if err != nil {
		t.Fatalf("GetTPMStatus failed: %v", err)
	}

	if result == nil {
		t.Fatal("GetTPMStatus returned nil result")
	}

	// Verify structure
	t.Logf("TPM Present: %v", result.Present)
	t.Logf("TPM Enabled: %v", result.Enabled)
	t.Logf("TPM Type: %s", result.Type)
	t.Logf("TPM Platform: %s", result.Platform)
}

func TestGetSecureBootStatus_WhenSupported(t *testing.T) {
	if !IsSecureBootSupported() {
		t.Skip("Secure Boot not supported on this platform")
	}

	result, err := GetSecureBootStatus()
	if err != nil {
		t.Fatalf("GetSecureBootStatus failed: %v", err)
	}

	if result == nil {
		t.Fatal("GetSecureBootStatus returned nil result")
	}

	t.Logf("Secure Boot Enabled: %v", result.Enabled)
	t.Logf("Secure Boot Mode: %s", result.Mode)
}

func TestGetEncryptionStatus_WhenSupported(t *testing.T) {
	if !IsEncryptionSupported() {
		t.Skip("Encryption not supported on this platform")
	}

	result, err := GetEncryptionStatus()
	if err != nil {
		t.Fatalf("GetEncryptionStatus failed: %v", err)
	}

	if result == nil {
		t.Fatal("GetEncryptionStatus returned nil result")
	}

	t.Logf("Encryption Enabled: %v", result.Enabled)
	t.Logf("Encryption Type: %s", result.Type)
	t.Logf("Encryption Status: %s", result.Status)
}

func TestGetBiometricCapabilities_WhenSupported(t *testing.T) {
	if !IsBiometricsSupported() {
		t.Skip("Biometrics not supported on this platform")
	}

	result, err := GetBiometricCapabilities()
	if err != nil {
		t.Fatalf("GetBiometricCapabilities failed: %v", err)
	}

	if result == nil {
		t.Fatal("GetBiometricCapabilities returned nil result")
	}

	t.Logf("Touch ID Available: %v", result.TouchIDAvailable)
	t.Logf("Touch ID Enrolled: %v", result.TouchIDEnrolled)
	t.Logf("Face ID Available: %v", result.FaceIDAvailable)
	t.Logf("Face ID Enrolled: %v", result.FaceIDEnrolled)
	t.Logf("Biometry Type: %s", result.BiometryType)
}

func TestFormatTPM_WhenSupported(t *testing.T) {
	if !IsTPMSupported() {
		t.Skip("TPM not supported on this platform")
	}

	result, err := GetTPMStatus()
	if err != nil {
		t.Fatalf("GetTPMStatus failed: %v", err)
	}

	// Test JSON format
	jsonOutput := FormatTPM(result, "json")
	if jsonOutput == "" {
		t.Error("JSON output should not be empty")
	}

	// Test table format
	tableOutput := FormatTPM(result, "table")
	if tableOutput == "" {
		t.Error("Table output should not be empty")
	}
}

func TestFormatSecureBoot_WhenSupported(t *testing.T) {
	if !IsSecureBootSupported() {
		t.Skip("Secure Boot not supported on this platform")
	}

	result, err := GetSecureBootStatus()
	if err != nil {
		t.Fatalf("GetSecureBootStatus failed: %v", err)
	}

	// Test JSON format
	jsonOutput := FormatSecureBoot(result, "json")
	if jsonOutput == "" {
		t.Error("JSON output should not be empty")
	}

	// Test table format
	tableOutput := FormatSecureBoot(result, "table")
	if tableOutput == "" {
		t.Error("Table output should not be empty")
	}
}

func TestFormatEncryption_WhenSupported(t *testing.T) {
	if !IsEncryptionSupported() {
		t.Skip("Encryption not supported on this platform")
	}

	result, err := GetEncryptionStatus()
	if err != nil {
		t.Fatalf("GetEncryptionStatus failed: %v", err)
	}

	// Test JSON format
	jsonOutput := FormatEncryption(result, "json")
	if jsonOutput == "" {
		t.Error("JSON output should not be empty")
	}

	// Test table format
	tableOutput := FormatEncryption(result, "table")
	if tableOutput == "" {
		t.Error("Table output should not be empty")
	}
}

func TestFormatBiometricCapabilities_WhenSupported(t *testing.T) {
	if !IsBiometricsSupported() {
		t.Skip("Biometrics not supported on this platform")
	}

	result, err := GetBiometricCapabilities()
	if err != nil {
		t.Fatalf("GetBiometricCapabilities failed: %v", err)
	}

	// Test JSON format
	jsonOutput := FormatBiometricCapabilities(result, "json")
	if jsonOutput == "" {
		t.Error("JSON output should not be empty")
	}

	// Test table format
	tableOutput := FormatBiometricCapabilities(result, "table")
	if tableOutput == "" {
		t.Error("Table output should not be empty")
	}
}

// Benchmark tests

func BenchmarkGetTPMStatus(b *testing.B) {
	if !IsTPMSupported() {
		b.Skip("TPM not supported on this platform")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetTPMStatus()
	}
}

func BenchmarkGetSecureBootStatus(b *testing.B) {
	if !IsSecureBootSupported() {
		b.Skip("Secure Boot not supported on this platform")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetSecureBootStatus()
	}
}

func BenchmarkGetEncryptionStatus(b *testing.B) {
	if !IsEncryptionSupported() {
		b.Skip("Encryption not supported on this platform")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetEncryptionStatus()
	}
}

func BenchmarkGetSecuritySummary(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetSecuritySummary()
	}
}
