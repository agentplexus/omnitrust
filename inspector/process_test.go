package inspector

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestListProcesses(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := ListProcesses(ctx, 0)
	if err != nil {
		t.Fatalf("ListProcesses failed: %v", err)
	}

	if result == nil {
		t.Fatal("ListProcesses returned nil result")
	}

	// Should have at least one process (this test itself)
	if result.Total == 0 {
		t.Error("Total should be > 0")
	}

	if len(result.Processes) == 0 {
		t.Error("Should have at least one process")
	}

	// Verify process structure
	for i, proc := range result.Processes {
		// PID 0 is valid on Windows (System Idle Process)
		if proc.PID < 0 {
			t.Errorf("Process[%d].PID = %d, want >= 0", i, proc.PID)
		}
		// Name can be empty for some system processes, so we don't check it
		// CPU and memory percentages should be non-negative
		if proc.CPUPercent < 0 {
			t.Errorf("Process[%d].CPUPercent = %.2f, want >= 0", i, proc.CPUPercent)
		}
		if proc.MemoryPercent < 0 {
			t.Errorf("Process[%d].MemoryPercent = %.2f, want >= 0", i, proc.MemoryPercent)
		}
	}
}

func TestListProcesses_WithLimit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	limit := 5
	result, err := ListProcesses(ctx, limit)
	if err != nil {
		t.Fatalf("ListProcesses with limit failed: %v", err)
	}

	if result == nil {
		t.Fatal("ListProcesses returned nil result")
	}

	// Should respect the limit
	if len(result.Processes) > limit {
		t.Errorf("len(Processes) = %d, want <= %d", len(result.Processes), limit)
	}

	// Total should still reflect all processes
	if result.Total < len(result.Processes) {
		t.Errorf("Total (%d) should be >= len(Processes) (%d)", result.Total, len(result.Processes))
	}
}

func TestListProcesses_Sorted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := ListProcesses(ctx, 10)
	if err != nil {
		t.Fatalf("ListProcesses failed: %v", err)
	}

	// Verify sorted by CPU descending
	for i := 1; i < len(result.Processes); i++ {
		if result.Processes[i].CPUPercent > result.Processes[i-1].CPUPercent {
			t.Errorf("Processes not sorted by CPU: [%d]=%.2f > [%d]=%.2f",
				i, result.Processes[i].CPUPercent,
				i-1, result.Processes[i-1].CPUPercent)
		}
	}
}

func TestProcessListResult_JSON(t *testing.T) {
	result := &ProcessListResult{
		Total: 100,
		Processes: []ProcessInfo{
			{PID: 1, Name: "init", CPUPercent: 0.1, MemoryPercent: 0.5, Status: "S"},
			{PID: 2, Name: "kernel", CPUPercent: 0.0, MemoryPercent: 0.0, Status: "S"},
		},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal ProcessListResult: %v", err)
	}

	// Verify JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if _, ok := parsed["total"]; !ok {
		t.Error("JSON should contain 'total' field")
	}
	if _, ok := parsed["processes"]; !ok {
		t.Error("JSON should contain 'processes' field")
	}
}

func TestProcessInfo_JSON(t *testing.T) {
	proc := ProcessInfo{
		PID:           12345,
		Name:          "test_process",
		CPUPercent:    25.5,
		MemoryPercent: 3.2,
		Status:        "R",
	}

	data, err := json.Marshal(proc)
	if err != nil {
		t.Fatalf("Failed to marshal ProcessInfo: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	expectedFields := []string{"pid", "name", "cpu_percent", "memory_percent", "status"}
	for _, field := range expectedFields {
		if _, ok := parsed[field]; !ok {
			t.Errorf("JSON should contain '%s' field", field)
		}
	}
}

func TestFormatProcessListTable(t *testing.T) {
	result := &ProcessListResult{
		Total: 100,
		Processes: []ProcessInfo{
			{PID: 1234, Name: "high_cpu_process", CPUPercent: 75.5, MemoryPercent: 2.0, Status: "R"},
			{PID: 5678, Name: "normal_process", CPUPercent: 5.0, MemoryPercent: 1.0, Status: "S"},
			{PID: 9012, Name: "idle_process", CPUPercent: 0.0, MemoryPercent: 0.5, Status: "I"},
		},
	}

	output := FormatProcessListTable(result)

	// Should contain header
	if !strings.Contains(output, "Processes") {
		t.Error("Output should contain 'Processes' header")
	}

	// Should contain total count
	if !strings.Contains(output, "100") {
		t.Error("Output should contain total count")
	}

	// Should contain column headers
	headers := []string{"PID", "Name", "CPU", "Mem", "Status"}
	for _, header := range headers {
		if !strings.Contains(output, header) {
			t.Errorf("Output should contain '%s' header", header)
		}
	}

	// Should contain process names
	if !strings.Contains(output, "high_cpu_process") {
		t.Error("Output should contain process name")
	}

	// Should have table characters
	if !strings.Contains(output, "┌") || !strings.Contains(output, "└") {
		t.Error("Output should contain table border characters")
	}
}

func TestFormatProcessList(t *testing.T) {
	result := &ProcessListResult{
		Total: 50,
		Processes: []ProcessInfo{
			{PID: 100, Name: "test", CPUPercent: 10.0, MemoryPercent: 5.0, Status: "R"},
		},
	}

	// Test JSON format
	jsonOutput := FormatProcessList(result, "json")
	if !strings.Contains(jsonOutput, "processes") {
		t.Error("JSON format should contain 'processes'")
	}
	if !strings.Contains(jsonOutput, "total") {
		t.Error("JSON format should contain 'total'")
	}

	// Test table format
	tableOutput := FormatProcessList(result, "table")
	if !strings.Contains(tableOutput, "Processes") {
		t.Error("Table format should contain 'Processes'")
	}
}

func TestFormatStatus(t *testing.T) {
	tests := []struct {
		status   string
		contains string
	}{
		{"R", "Run"},
		{"running", "Run"},
		{"S", "Sleep"},
		{"sleep", "Sleep"},
		{"I", "Idle"},
		{"idle", "Idle"},
		{"Z", "Zombie"},
		{"zombie", "Zombie"},
		{"T", "Stop"},
		{"stop", "Stop"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			result := formatStatus(tt.status)
			stripped := StripANSI(result)
			if !strings.Contains(stripped, tt.contains) {
				t.Errorf("formatStatus(%q) = %q, should contain %q", tt.status, stripped, tt.contains)
			}
		})
	}
}

func TestFormatProcessListTable_LongName(t *testing.T) {
	result := &ProcessListResult{
		Total: 1,
		Processes: []ProcessInfo{
			{
				PID:           1,
				Name:          "this_is_a_very_long_process_name_that_should_be_truncated",
				CPUPercent:    1.0,
				MemoryPercent: 1.0,
				Status:        "R",
			},
		},
	}

	output := FormatProcessListTable(result)

	// Should truncate long names with "..."
	if !strings.Contains(output, "...") {
		t.Error("Long process names should be truncated with ...")
	}
}

func TestFormatProcessListTable_Empty(t *testing.T) {
	result := &ProcessListResult{
		Total:     0,
		Processes: []ProcessInfo{},
	}

	// Should not panic with empty list
	output := FormatProcessListTable(result)
	if output == "" {
		t.Error("Output should not be empty even with no processes")
	}
}

func TestFormatProcessListTable_HighUsage(t *testing.T) {
	result := &ProcessListResult{
		Total: 2,
		Processes: []ProcessInfo{
			{PID: 1, Name: "cpu_hog", CPUPercent: 95.0, MemoryPercent: 2.0, Status: "R"},
			{PID: 2, Name: "mem_hog", CPUPercent: 5.0, MemoryPercent: 15.0, Status: "R"},
		},
	}

	output := FormatProcessListTable(result)

	// High CPU/memory should use warning/danger colors
	if !strings.Contains(output, Red) && !strings.Contains(output, Yellow) {
		t.Error("High usage should use warning/danger colors")
	}
}
