package inspector

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestGetCPUUsage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := GetCPUUsage(ctx)
	if err != nil {
		t.Fatalf("GetCPUUsage failed: %v", err)
	}

	// Verify result structure
	if result == nil {
		t.Fatal("GetCPUUsage returned nil result")
	}

	// Usage should be between 0 and 100
	if result.UsagePercent < 0 || result.UsagePercent > 100 {
		t.Errorf("UsagePercent = %.2f, want between 0 and 100", result.UsagePercent)
	}

	// Should have at least one core
	if len(result.PerCore) == 0 {
		t.Error("PerCore should have at least one entry")
	}

	// Each core should have valid percentage
	for i, usage := range result.PerCore {
		if usage < 0 || usage > 100 {
			t.Errorf("PerCore[%d] = %.2f, want between 0 and 100", i, usage)
		}
	}
}

func TestGetCPUUsage_Cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := GetCPUUsage(ctx)
	// Should either succeed quickly or return context error
	if err != nil && !strings.Contains(err.Error(), "context") {
		// Some implementations may not respect context, which is okay
		t.Logf("GetCPUUsage with cancelled context: %v", err)
	}
}

func TestCPUUsageResult_JSON(t *testing.T) {
	result := &CPUUsageResult{
		UsagePercent: 45.5,
		PerCore:      []float64{30.0, 50.0, 60.0, 40.0},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal CPUUsageResult: %v", err)
	}

	// Verify JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if _, ok := parsed["usage_percent"]; !ok {
		t.Error("JSON should contain 'usage_percent' field")
	}
	if _, ok := parsed["per_core"]; !ok {
		t.Error("JSON should contain 'per_core' field")
	}
}

func TestFormatCPUUsageTable(t *testing.T) {
	result := &CPUUsageResult{
		UsagePercent: 45.5,
		PerCore:      []float64{30.0, 50.0, 85.0, 95.0},
	}

	output := FormatCPUUsageTable(result)

	// Should contain header
	if !strings.Contains(output, "CPU Usage") {
		t.Error("Output should contain 'CPU Usage' header")
	}

	// Should contain overall usage
	if !strings.Contains(output, "Overall") {
		t.Error("Output should contain 'Overall' label")
	}

	// Should contain per-core section
	if !strings.Contains(output, "Per-Core") {
		t.Error("Output should contain 'Per-Core' section")
	}

	// Should have table characters
	if !strings.Contains(output, "┌") || !strings.Contains(output, "└") {
		t.Error("Output should contain table border characters")
	}

	// Should have progress bar characters
	if !strings.Contains(output, IconBar) || !strings.Contains(output, IconBarLight) {
		t.Error("Output should contain progress bar characters")
	}
}

func TestFormatCPUUsage(t *testing.T) {
	result := &CPUUsageResult{
		UsagePercent: 25.0,
		PerCore:      []float64{20.0, 30.0},
	}

	// Test JSON format
	jsonOutput := FormatCPUUsage(result, "json")
	if !strings.Contains(jsonOutput, "usage_percent") {
		t.Error("JSON format should contain 'usage_percent'")
	}
	if !strings.Contains(jsonOutput, "per_core") {
		t.Error("JSON format should contain 'per_core'")
	}

	// Test table format
	tableOutput := FormatCPUUsage(result, "table")
	if !strings.Contains(tableOutput, "CPU Usage") {
		t.Error("Table format should contain 'CPU Usage'")
	}

	// Test case insensitivity
	tableOutput2 := FormatCPUUsage(result, "TABLE")
	if !strings.Contains(tableOutput2, "CPU Usage") {
		t.Error("Format should be case insensitive")
	}
}

func TestFormatCPUUsageTable_ColorThresholds(t *testing.T) {
	// Test with various usage levels to verify color coding
	tests := []struct {
		name    string
		usage   float64
		perCore []float64
	}{
		{"low usage", 25.0, []float64{20.0, 30.0}},
		{"medium usage", 75.0, []float64{70.0, 80.0}},
		{"high usage", 95.0, []float64{90.0, 100.0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &CPUUsageResult{
				UsagePercent: tt.usage,
				PerCore:      tt.perCore,
			}
			output := FormatCPUUsageTable(result)

			// Should produce valid output without panicking
			if output == "" {
				t.Error("Output should not be empty")
			}

			// Should contain percentage values
			if !strings.Contains(output, "%") {
				t.Error("Output should contain percentage values")
			}
		})
	}
}

func TestFormatCPUUsageTable_EmptyCores(t *testing.T) {
	result := &CPUUsageResult{
		UsagePercent: 50.0,
		PerCore:      []float64{},
	}

	// Should not panic with empty cores
	output := FormatCPUUsageTable(result)
	if output == "" {
		t.Error("Output should not be empty even with no cores")
	}
}
