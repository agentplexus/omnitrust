package inspector

import (
	"context"
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v4/cpu"
)

// CPUUsageResult contains CPU usage information
type CPUUsageResult struct {
	UsagePercent float64   `json:"usage_percent"`
	PerCore      []float64 `json:"per_core"`
}

// GetCPUUsage returns current CPU usage
func GetCPUUsage(ctx context.Context) (*CPUUsageResult, error) {
	overall, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get overall CPU usage: %w", err)
	}

	perCore, err := cpu.PercentWithContext(ctx, 0, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get per-core CPU usage: %w", err)
	}

	var overallUsage float64
	if len(overall) > 0 {
		overallUsage = overall[0]
	}

	return &CPUUsageResult{
		UsagePercent: overallUsage,
		PerCore:      perCore,
	}, nil
}

// FormatCPUUsageTable formats CPU usage as a colored table
func FormatCPUUsageTable(result *CPUUsageResult) string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(Header(IconCPU + " CPU Usage"))
	sb.WriteString("\n")
	sb.WriteString(Muted(strings.Repeat("─", 40)))
	sb.WriteString("\n\n")

	// Overall usage with progress bar
	sb.WriteString(BoldText("Overall: "))
	usageColor := UsageColor(result.UsagePercent)
	sb.WriteString(Colorize(usageColor+Bold, fmt.Sprintf("%.1f%%", result.UsagePercent)))
	sb.WriteString("\n")
	sb.WriteString(ProgressBar(result.UsagePercent, 30))
	sb.WriteString("\n\n")

	// Per-core table
	sb.WriteString(BoldText("Per-Core Usage:"))
	sb.WriteString("\n")
	sb.WriteString(TableTop(6, 10, 20))
	sb.WriteString("\n")
	sb.WriteString(TableRowColored(
		Header(PadRight("Core", 6)),
		Header(PadLeft("Usage", 10)),
		Header(PadRight("", 20)),
	))
	sb.WriteString("\n")
	sb.WriteString(TableSeparator(6, 10, 20))
	sb.WriteString("\n")

	for i, usage := range result.PerCore {
		var usageStr string
		switch {
		case usage >= 90:
			usageStr = Danger(fmt.Sprintf("%6.1f%%", usage))
		case usage >= 70:
			usageStr = Warning(fmt.Sprintf("%6.1f%%", usage))
		default:
			usageStr = Success(fmt.Sprintf("%6.1f%%", usage))
		}
		sb.WriteString(TableRowColored(
			Info(PadRight(fmt.Sprintf("%s %d", IconCore, i), 6)),
			PadLeft(usageStr, 10),
			ProgressBar(usage, 20),
		))
		sb.WriteString("\n")
	}
	sb.WriteString(TableBottom(6, 10, 20))
	sb.WriteString("\n")
	return sb.String()
}

// FormatCPUUsage formats CPU usage in the specified format
func FormatCPUUsage(result *CPUUsageResult, format string) string {
	return FormatOutput(result, func() string {
		return FormatCPUUsageTable(result)
	}, format)
}
