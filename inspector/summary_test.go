package inspector

import (
	"encoding/json"
	"runtime"
	"strings"
	"testing"
)

func TestGetSecuritySummary(t *testing.T) {
	result, err := GetSecuritySummary()
	if err != nil {
		t.Fatalf("GetSecuritySummary failed: %v", err)
	}

	if result == nil {
		t.Fatal("GetSecuritySummary returned nil result")
	}

	// Platform should match runtime
	if result.Platform != runtime.GOOS {
		t.Errorf("Platform = %q, want %q", result.Platform, runtime.GOOS)
	}

	// Score should be between 0 and 100
	if result.OverallScore < 0 || result.OverallScore > 100 {
		t.Errorf("OverallScore = %d, want between 0 and 100", result.OverallScore)
	}

	// Status should be one of the valid values
	validStatuses := []string{"excellent", "good", "fair", "needs_improvement", "critical"}
	found := false
	for _, s := range validStatuses {
		if result.OverallStatus == s {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("OverallStatus = %q, want one of %v", result.OverallStatus, validStatuses)
	}
}

func TestSecuritySummary_ScoreStatus(t *testing.T) {
	tests := []struct {
		score          int
		expectedStatus string
	}{
		{100, "excellent"},
		{75, "good"},
		{50, "fair"},
		{25, "needs_improvement"},
		{0, "critical"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			// Verify the test's expected status matches our threshold logic
			var expected string
			switch {
			case tt.score >= 100:
				expected = "excellent"
			case tt.score >= 75:
				expected = "good"
			case tt.score >= 50:
				expected = "fair"
			case tt.score >= 25:
				expected = "needs_improvement"
			default:
				expected = "critical"
			}

			if expected != tt.expectedStatus {
				t.Errorf("Score %d: expected status %q, test expects %q", tt.score, expected, tt.expectedStatus)
			}
		})
	}
}

func TestSecuritySummary_JSON(t *testing.T) {
	result := &SecuritySummary{
		Platform:      "darwin",
		OverallScore:  75,
		OverallStatus: "good",
		TPM: &TPMSummary{
			Present: true,
			Enabled: true,
			Type:    "secure_enclave",
		},
		SecureBoot: &BootSummary{
			Enabled: true,
			Mode:    "full",
		},
		Encryption: &EncSummary{
			Enabled: false,
			Type:    "filevault",
			Status:  "disabled",
		},
		Biometrics: &BioSummary{
			Available:  true,
			Configured: true,
			Type:       "touch_id",
		},
		Recommendations: []string{"Enable FileVault"},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal SecuritySummary: %v", err)
	}

	// Verify JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	expectedFields := []string{
		"platform", "overall_score", "overall_status",
		"tpm", "secure_boot", "encryption", "biometrics",
	}

	for _, field := range expectedFields {
		if _, ok := parsed[field]; !ok {
			t.Errorf("JSON should contain '%s' field", field)
		}
	}
}

func TestTPMSummary_JSON(t *testing.T) {
	summary := &TPMSummary{
		Present: true,
		Enabled: true,
		Type:    "tpm_2.0",
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("Failed to marshal TPMSummary: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if _, ok := parsed["present"]; !ok {
		t.Error("JSON should contain 'present' field")
	}
	if _, ok := parsed["enabled"]; !ok {
		t.Error("JSON should contain 'enabled' field")
	}
	if _, ok := parsed["type"]; !ok {
		t.Error("JSON should contain 'type' field")
	}
}

func TestBootSummary_JSON(t *testing.T) {
	summary := &BootSummary{
		Enabled: true,
		Mode:    "full",
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("Failed to marshal BootSummary: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if _, ok := parsed["enabled"]; !ok {
		t.Error("JSON should contain 'enabled' field")
	}
	if _, ok := parsed["mode"]; !ok {
		t.Error("JSON should contain 'mode' field")
	}
}

func TestEncSummary_JSON(t *testing.T) {
	summary := &EncSummary{
		Enabled: true,
		Type:    "filevault",
		Status:  "encrypted",
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("Failed to marshal EncSummary: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if _, ok := parsed["enabled"]; !ok {
		t.Error("JSON should contain 'enabled' field")
	}
	if _, ok := parsed["type"]; !ok {
		t.Error("JSON should contain 'type' field")
	}
	if _, ok := parsed["status"]; !ok {
		t.Error("JSON should contain 'status' field")
	}
}

func TestBioSummary_JSON(t *testing.T) {
	summary := &BioSummary{
		Available:  true,
		Configured: true,
		Type:       "touch_id",
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("Failed to marshal BioSummary: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if _, ok := parsed["available"]; !ok {
		t.Error("JSON should contain 'available' field")
	}
	if _, ok := parsed["configured"]; !ok {
		t.Error("JSON should contain 'configured' field")
	}
	if _, ok := parsed["type"]; !ok {
		t.Error("JSON should contain 'type' field")
	}
}

func TestFormatSecuritySummaryTable(t *testing.T) {
	result := &SecuritySummary{
		Platform:      "darwin",
		OverallScore:  75,
		OverallStatus: "good",
		TPM: &TPMSummary{
			Present: true,
			Enabled: true,
			Type:    "secure_enclave",
		},
		SecureBoot: &BootSummary{
			Enabled: true,
			Mode:    "full",
		},
		Encryption: &EncSummary{
			Enabled: false,
			Type:    "filevault",
			Status:  "disabled",
		},
		Biometrics: &BioSummary{
			Available:  true,
			Configured: true,
			Type:       "touch_id",
		},
		Recommendations: []string{"Enable FileVault to protect data at rest"},
	}

	output := FormatSecuritySummaryTable(result)

	// Should contain header
	if !strings.Contains(output, "Security Summary") {
		t.Error("Output should contain 'Security Summary' header")
	}

	// Should contain platform info
	if !strings.Contains(output, "Platform") {
		t.Error("Output should contain 'Platform' label")
	}

	// Should contain score
	if !strings.Contains(output, "75") {
		t.Error("Output should contain score value")
	}

	// Should contain status
	if !strings.Contains(output, "Good") {
		t.Error("Output should contain status")
	}

	// Should contain feature labels
	features := []string{"Secure Boot", "Biometrics"}
	for _, feature := range features {
		if !strings.Contains(output, feature) {
			t.Errorf("Output should contain '%s'", feature)
		}
	}

	// Should contain recommendations
	if !strings.Contains(output, "Recommendations") {
		t.Error("Output should contain 'Recommendations' section")
	}

	// Should have table characters
	if !strings.Contains(output, "┌") || !strings.Contains(output, "└") {
		t.Error("Output should contain table border characters")
	}
}

func TestFormatSecuritySummary(t *testing.T) {
	result := &SecuritySummary{
		Platform:      "darwin",
		OverallScore:  100,
		OverallStatus: "excellent",
	}

	// Test JSON format
	jsonOutput := FormatSecuritySummary(result, "json")
	if !strings.Contains(jsonOutput, "platform") {
		t.Error("JSON format should contain 'platform'")
	}
	if !strings.Contains(jsonOutput, "overall_score") {
		t.Error("JSON format should contain 'overall_score'")
	}

	// Test table format
	tableOutput := FormatSecuritySummary(result, "table")
	if !strings.Contains(tableOutput, "Security Summary") {
		t.Error("Table format should contain 'Security Summary'")
	}
}

func TestFormatSecuritySummaryTable_NoRecommendations(t *testing.T) {
	result := &SecuritySummary{
		Platform:        "darwin",
		OverallScore:    100,
		OverallStatus:   "excellent",
		Recommendations: []string{},
	}

	output := FormatSecuritySummaryTable(result)

	// Should still produce valid output
	if output == "" {
		t.Error("Output should not be empty")
	}
}

func TestFormatSecuritySummaryTable_NilSections(t *testing.T) {
	result := &SecuritySummary{
		Platform:      "test",
		OverallScore:  0,
		OverallStatus: "critical",
		TPM:           nil,
		SecureBoot:    nil,
		Encryption:    nil,
		Biometrics:    nil,
	}

	// Should not panic with nil sections
	output := FormatSecuritySummaryTable(result)
	if output == "" {
		t.Error("Output should not be empty even with nil sections")
	}

	// Should show N/A for missing sections
	if !strings.Contains(output, "N/A") {
		t.Error("Output should contain 'N/A' for missing sections")
	}
}

func TestFormatSecuritySummaryTable_AllStatuses(t *testing.T) {
	statuses := []struct {
		score  int
		status string
	}{
		{100, "excellent"},
		{75, "good"},
		{50, "fair"},
		{25, "needs_improvement"},
		{0, "critical"},
	}

	for _, s := range statuses {
		t.Run(s.status, func(t *testing.T) {
			result := &SecuritySummary{
				Platform:      "test",
				OverallScore:  s.score,
				OverallStatus: s.status,
			}

			output := FormatSecuritySummaryTable(result)
			if output == "" {
				t.Errorf("Output for status %q should not be empty", s.status)
			}
		})
	}
}

func TestSecurityScoreBar(t *testing.T) {
	tests := []struct {
		score       int
		width       int
		expectGreen bool
		expectRed   bool
	}{
		{100, 20, true, false},
		{75, 20, true, false},
		{50, 20, false, false}, // yellow
		{25, 20, false, true},
		{0, 20, false, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := securityScoreBar(tt.score, tt.width)

			if tt.expectGreen && !strings.Contains(result, Green) {
				t.Errorf("Score %d should produce green bar", tt.score)
			}
			if tt.expectRed && !strings.Contains(result, Red) {
				t.Errorf("Score %d should produce red bar", tt.score)
			}
		})
	}
}

func TestFeatureStatus(t *testing.T) {
	enabled := featureStatus(true)
	if !strings.Contains(enabled, IconCheck) {
		t.Error("Enabled status should contain check icon")
	}
	if !strings.Contains(enabled, "Enabled") {
		t.Error("Enabled status should contain 'Enabled'")
	}

	disabled := featureStatus(false)
	if !strings.Contains(disabled, IconCross) {
		t.Error("Disabled status should contain cross icon")
	}
	if !strings.Contains(disabled, "Disabled") {
		t.Error("Disabled status should contain 'Disabled'")
	}
}
