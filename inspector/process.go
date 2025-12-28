package inspector

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

// ProcessInfo contains information about a single process
type ProcessInfo struct {
	PID           int32   `json:"pid"`
	Name          string  `json:"name"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float32 `json:"memory_percent"`
	Status        string  `json:"status"`
}

// ProcessListResult contains the process list result
type ProcessListResult struct {
	Processes []ProcessInfo `json:"processes"`
	Total     int           `json:"total"`
}

// ListProcesses returns a list of running processes
func ListProcesses(ctx context.Context, limit int) (*ProcessListResult, error) {
	procs, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list processes: %w", err)
	}

	var procInfos []ProcessInfo
	for _, p := range procs {
		name, _ := p.NameWithContext(ctx)
		cpuPercent, _ := p.CPUPercentWithContext(ctx)
		memPercent, _ := p.MemoryPercentWithContext(ctx)
		status, _ := p.StatusWithContext(ctx)

		statusStr := "unknown"
		if len(status) > 0 {
			statusStr = status[0]
		}

		procInfos = append(procInfos, ProcessInfo{
			PID:           p.Pid,
			Name:          name,
			CPUPercent:    cpuPercent,
			MemoryPercent: memPercent,
			Status:        statusStr,
		})
	}

	// Sort by CPU usage descending
	sort.Slice(procInfos, func(i, j int) bool {
		return procInfos[i].CPUPercent > procInfos[j].CPUPercent
	})

	total := len(procInfos)
	if limit > 0 && limit < len(procInfos) {
		procInfos = procInfos[:limit]
	}

	return &ProcessListResult{
		Processes: procInfos,
		Total:     total,
	}, nil
}

// formatStatus returns a colored status string
func formatStatus(status string) string {
	switch status {
	case "R", "running":
		return Success(IconCircle + " Run")
	case "S", "sleep":
		return Info(IconCircle + " Sleep")
	case "I", "idle":
		return Muted(IconCircle + " Idle")
	case "Z", "zombie":
		return Danger(IconCircle + " Zombie")
	case "T", "stop":
		return Warning(IconCircle + " Stop")
	default:
		return Muted(IconCircle + " " + status)
	}
}

// FormatProcessListTable formats process list as a colored table
func FormatProcessListTable(result *ProcessListResult) string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(Header(fmt.Sprintf("%s Processes (Total: %d)", IconProcess, result.Total)))
	sb.WriteString("\n")
	sb.WriteString(Muted(strings.Repeat("─", 70)))
	sb.WriteString("\n\n")

	// Process table
	sb.WriteString(TableTop(8, 28, 9, 9, 10))
	sb.WriteString("\n")
	sb.WriteString(TableRowColored(
		Header(PadRight("PID", 8)),
		Header(PadRight("Name", 28)),
		Header(PadLeft("CPU %", 9)),
		Header(PadLeft("Mem %", 9)),
		Header(PadRight("Status", 10)),
	))
	sb.WriteString("\n")
	sb.WriteString(TableSeparator(8, 28, 9, 9, 10))
	sb.WriteString("\n")

	for _, proc := range result.Processes {
		// Color CPU based on usage
		var cpuStr string
		switch {
		case proc.CPUPercent >= 50:
			cpuStr = Danger(fmt.Sprintf("%9.1f", proc.CPUPercent))
		case proc.CPUPercent >= 25:
			cpuStr = Warning(fmt.Sprintf("%9.1f", proc.CPUPercent))
		default:
			cpuStr = fmt.Sprintf("%9.1f", proc.CPUPercent)
		}

		// Color memory based on usage
		var memStr string
		switch {
		case proc.MemoryPercent >= 10:
			memStr = Danger(fmt.Sprintf("%9.1f", proc.MemoryPercent))
		case proc.MemoryPercent >= 5:
			memStr = Warning(fmt.Sprintf("%9.1f", proc.MemoryPercent))
		default:
			memStr = fmt.Sprintf("%9.1f", proc.MemoryPercent)
		}

		// Truncate name if too long
		name := proc.Name
		if len(name) > 28 {
			name = name[:25] + "..."
		}

		sb.WriteString(TableRowColored(
			Info(PadRight(fmt.Sprintf("%d", proc.PID), 8)),
			PadRight(name, 28),
			cpuStr,
			memStr,
			PadRight(formatStatus(proc.Status), 10),
		))
		sb.WriteString("\n")
	}

	sb.WriteString(TableBottom(8, 28, 9, 9, 10))
	sb.WriteString("\n")
	return sb.String()
}

// FormatProcessList formats process list in the specified format
func FormatProcessList(result *ProcessListResult, format string) string {
	return FormatOutput(result, func() string {
		return FormatProcessListTable(result)
	}, format)
}
