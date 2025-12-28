package inspector

import (
	"context"
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v4/mem"
)

// MemoryResult contains memory usage information
type MemoryResult struct {
	TotalBytes     uint64  `json:"total_bytes"`
	UsedBytes      uint64  `json:"used_bytes"`
	FreeBytes      uint64  `json:"free_bytes"`
	AvailableBytes uint64  `json:"available_bytes"`
	UsedPercent    float64 `json:"used_percent"`
	TotalHuman     string  `json:"total_human"`
	UsedHuman      string  `json:"used_human"`
	AvailableHuman string  `json:"available_human"`
}

// GetMemory returns current memory usage
func GetMemory(ctx context.Context) (*MemoryResult, error) {
	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory stats: %w", err)
	}

	return &MemoryResult{
		TotalBytes:     vmStat.Total,
		UsedBytes:      vmStat.Used,
		FreeBytes:      vmStat.Free,
		AvailableBytes: vmStat.Available,
		UsedPercent:    vmStat.UsedPercent,
		TotalHuman:     FormatBytes(vmStat.Total),
		UsedHuman:      FormatBytes(vmStat.Used),
		AvailableHuman: FormatBytes(vmStat.Available),
	}, nil
}

// FormatMemoryTable formats memory usage as a colored table
func FormatMemoryTable(result *MemoryResult) string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(Header(IconMemory + " Memory Usage"))
	sb.WriteString("\n")
	sb.WriteString(Muted(strings.Repeat("─", 50)))
	sb.WriteString("\n\n")

	// Usage summary with progress bar
	sb.WriteString(BoldText("Usage: "))
	usageColor := UsageColor(result.UsedPercent)
	sb.WriteString(Colorize(usageColor+Bold, fmt.Sprintf("%.1f%%", result.UsedPercent)))
	sb.WriteString(Muted(" of "))
	sb.WriteString(Info(result.TotalHuman))
	sb.WriteString("\n")
	sb.WriteString(ProgressBar(result.UsedPercent, 40))
	sb.WriteString("\n\n")

	// Memory details table
	sb.WriteString(TableTop(12, 14, 20))
	sb.WriteString("\n")
	sb.WriteString(TableRowColored(
		Header(PadRight("Metric", 12)),
		Header(PadLeft("Size", 14)),
		Header(PadLeft("Bytes", 20)),
	))
	sb.WriteString("\n")
	sb.WriteString(TableSeparator(12, 14, 20))
	sb.WriteString("\n")

	// Total
	sb.WriteString(TableRowColored(
		Info(PadRight(IconDiamond+" Total", 12)),
		PadLeft(result.TotalHuman, 14),
		Muted(PadLeft(fmt.Sprintf("%d", result.TotalBytes), 20)),
	))
	sb.WriteString("\n")

	// Used
	usedColor := UsageColor(result.UsedPercent)
	sb.WriteString(TableRowColored(
		Colorize(usedColor, PadRight(IconCircle+" Used", 12)),
		Colorize(usedColor, PadLeft(result.UsedHuman, 14)),
		Muted(PadLeft(fmt.Sprintf("%d", result.UsedBytes), 20)),
	))
	sb.WriteString("\n")

	// Free
	sb.WriteString(TableRowColored(
		Success(PadRight(IconCircle+" Free", 12)),
		Success(PadLeft(FormatBytes(result.FreeBytes), 14)),
		Muted(PadLeft(fmt.Sprintf("%d", result.FreeBytes), 20)),
	))
	sb.WriteString("\n")

	// Available
	sb.WriteString(TableRowColored(
		Success(PadRight(IconCircle+" Available", 12)),
		Success(PadLeft(result.AvailableHuman, 14)),
		Muted(PadLeft(fmt.Sprintf("%d", result.AvailableBytes), 20)),
	))
	sb.WriteString("\n")

	sb.WriteString(TableBottom(12, 14, 20))
	sb.WriteString("\n")
	return sb.String()
}

// FormatMemory formats memory usage in the specified format
func FormatMemory(result *MemoryResult, format string) string {
	return FormatOutput(result, func() string {
		return FormatMemoryTable(result)
	}, format)
}
