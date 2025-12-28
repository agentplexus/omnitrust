package inspector

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestGetMemory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := GetMemory(ctx)
	if err != nil {
		t.Fatalf("GetMemory failed: %v", err)
	}

	if result == nil {
		t.Fatal("GetMemory returned nil result")
	}

	// Total should be > 0
	if result.TotalBytes == 0 {
		t.Error("TotalBytes should be > 0")
	}

	// Used should be <= Total
	if result.UsedBytes > result.TotalBytes {
		t.Errorf("UsedBytes (%d) should be <= TotalBytes (%d)", result.UsedBytes, result.TotalBytes)
	}

	// Available should be <= Total
	if result.AvailableBytes > result.TotalBytes {
		t.Errorf("AvailableBytes (%d) should be <= TotalBytes (%d)", result.AvailableBytes, result.TotalBytes)
	}

	// UsedPercent should be between 0 and 100
	if result.UsedPercent < 0 || result.UsedPercent > 100 {
		t.Errorf("UsedPercent = %.2f, want between 0 and 100", result.UsedPercent)
	}

	// Human-readable strings should not be empty
	if result.TotalHuman == "" {
		t.Error("TotalHuman should not be empty")
	}
	if result.UsedHuman == "" {
		t.Error("UsedHuman should not be empty")
	}
	if result.AvailableHuman == "" {
		t.Error("AvailableHuman should not be empty")
	}
}

func TestMemoryResult_JSON(t *testing.T) {
	result := &MemoryResult{
		TotalBytes:     16 * 1024 * 1024 * 1024,
		UsedBytes:      8 * 1024 * 1024 * 1024,
		FreeBytes:      2 * 1024 * 1024 * 1024,
		AvailableBytes: 6 * 1024 * 1024 * 1024,
		UsedPercent:    50.0,
		TotalHuman:     "16.00 GB",
		UsedHuman:      "8.00 GB",
		AvailableHuman: "6.00 GB",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal MemoryResult: %v", err)
	}

	// Verify JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	expectedFields := []string{
		"total_bytes", "used_bytes", "free_bytes", "available_bytes",
		"used_percent", "total_human", "used_human", "available_human",
	}

	for _, field := range expectedFields {
		if _, ok := parsed[field]; !ok {
			t.Errorf("JSON should contain '%s' field", field)
		}
	}
}

func TestFormatMemoryTable(t *testing.T) {
	result := &MemoryResult{
		TotalBytes:     16 * 1024 * 1024 * 1024,
		UsedBytes:      8 * 1024 * 1024 * 1024,
		FreeBytes:      2 * 1024 * 1024 * 1024,
		AvailableBytes: 6 * 1024 * 1024 * 1024,
		UsedPercent:    50.0,
		TotalHuman:     "16.00 GB",
		UsedHuman:      "8.00 GB",
		AvailableHuman: "6.00 GB",
	}

	output := FormatMemoryTable(result)

	// Should contain header
	if !strings.Contains(output, "Memory Usage") {
		t.Error("Output should contain 'Memory Usage' header")
	}

	// Should contain usage summary
	if !strings.Contains(output, "Usage") {
		t.Error("Output should contain 'Usage' label")
	}

	// Should contain memory metrics
	metrics := []string{"Total", "Used", "Free", "Available"}
	for _, metric := range metrics {
		if !strings.Contains(output, metric) {
			t.Errorf("Output should contain '%s' metric", metric)
		}
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

func TestFormatMemory(t *testing.T) {
	result := &MemoryResult{
		TotalBytes:     16 * 1024 * 1024 * 1024,
		UsedBytes:      8 * 1024 * 1024 * 1024,
		FreeBytes:      2 * 1024 * 1024 * 1024,
		AvailableBytes: 6 * 1024 * 1024 * 1024,
		UsedPercent:    50.0,
		TotalHuman:     "16.00 GB",
		UsedHuman:      "8.00 GB",
		AvailableHuman: "6.00 GB",
	}

	// Test JSON format
	jsonOutput := FormatMemory(result, "json")
	if !strings.Contains(jsonOutput, "total_bytes") {
		t.Error("JSON format should contain 'total_bytes'")
	}

	// Test table format
	tableOutput := FormatMemory(result, "table")
	if !strings.Contains(tableOutput, "Memory Usage") {
		t.Error("Table format should contain 'Memory Usage'")
	}
}

func TestFormatMemoryTable_HighUsage(t *testing.T) {
	// Test with high memory usage
	result := &MemoryResult{
		TotalBytes:     16 * 1024 * 1024 * 1024,
		UsedBytes:      15 * 1024 * 1024 * 1024,
		FreeBytes:      512 * 1024 * 1024,
		AvailableBytes: 1 * 1024 * 1024 * 1024,
		UsedPercent:    93.75,
		TotalHuman:     "16.00 GB",
		UsedHuman:      "15.00 GB",
		AvailableHuman: "1.00 GB",
	}

	output := FormatMemoryTable(result)

	// Should still produce valid output
	if output == "" {
		t.Error("Output should not be empty")
	}

	// High usage should trigger red color (contained in ANSI codes)
	if !strings.Contains(output, Red) {
		t.Error("High usage should use red color")
	}
}

func TestFormatMemoryTable_LowUsage(t *testing.T) {
	// Test with low memory usage
	result := &MemoryResult{
		TotalBytes:     16 * 1024 * 1024 * 1024,
		UsedBytes:      4 * 1024 * 1024 * 1024,
		FreeBytes:      8 * 1024 * 1024 * 1024,
		AvailableBytes: 12 * 1024 * 1024 * 1024,
		UsedPercent:    25.0,
		TotalHuman:     "16.00 GB",
		UsedHuman:      "4.00 GB",
		AvailableHuman: "12.00 GB",
	}

	output := FormatMemoryTable(result)

	// Should still produce valid output
	if output == "" {
		t.Error("Output should not be empty")
	}

	// Low usage should trigger green color
	if !strings.Contains(output, Green) {
		t.Error("Low usage should use green color")
	}
}
