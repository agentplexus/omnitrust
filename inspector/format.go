package inspector

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

// OutputFormat constants
const (
	FormatJSON  = "json"
	FormatTable = "table"
)

// ANSI color codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	// Foreground colors
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// Bright foreground colors
	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	// Background colors
	BgBlue  = "\033[44m"
	BgCyan  = "\033[46m"
	BgWhite = "\033[47m"
)

// UTF-8 icons
const (
	IconCPU         = "🖥️ "
	IconMemory      = "💾"
	IconProcess     = "⚙️ "
	IconCheck       = "✓"
	IconCross       = "✗"
	IconCircle      = "●"
	IconDiamond     = "◆"
	IconArrow       = "→"
	IconBar         = "█"
	IconBarLight    = "░"
	IconBarMed      = "▒"
	IconCore        = "◉"
	IconPID         = "⬡"
	IconStatus      = "◈"
	IconWarning     = "⚠️ "
	IconInfo        = "ℹ️ "
	IconLock        = "🔒"
	IconUnlock      = "🔓"
	IconKey         = "🔑"
	IconShield      = "🛡️ "
	IconFingerprint = "👆"
	IconFace        = "👤"
	IconApple       = "🍎"
	IconChip        = "🔲"
)

// Colorize wraps text with a color and reset
func Colorize(color, text string) string {
	return color + text + Reset
}

// Bold makes text bold
func BoldText(text string) string {
	return Bold + text + Reset
}

// Dim makes text dimmed
func DimText(text string) string {
	return Dim + text + Reset
}

// Header formats text as a header (bold cyan)
func Header(text string) string {
	return Bold + Cyan + text + Reset
}

// Success formats text as success (green)
func Success(text string) string {
	return Green + text + Reset
}

// Warning formats text as warning (yellow)
func Warning(text string) string {
	return Yellow + text + Reset
}

// Danger formats text as danger (red)
func Danger(text string) string {
	return Red + text + Reset
}

// Info formats text as info (blue)
func Info(text string) string {
	return Blue + text + Reset
}

// Muted formats text as muted (gray)
func Muted(text string) string {
	return BrightBlack + text + Reset
}

// FormatBytes converts bytes to human-readable format
func FormatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

// TableRow creates a formatted table row with box-drawing characters
func TableRow(cols ...string) string {
	return "│ " + strings.Join(cols, " │ ") + " │"
}

// TableRowColored creates a colored table row
func TableRowColored(cols ...string) string {
	return Muted("│") + " " + strings.Join(cols, " "+Muted("│")+" ") + " " + Muted("│")
}

// TableSeparator creates a separator line for tables
func TableSeparator(widths ...int) string {
	parts := make([]string, len(widths))
	for i, w := range widths {
		parts[i] = strings.Repeat("─", w)
	}
	return Muted("├─" + strings.Join(parts, "─┼─") + "─┤")
}

// TableTop creates a top border for tables
func TableTop(widths ...int) string {
	parts := make([]string, len(widths))
	for i, w := range widths {
		parts[i] = strings.Repeat("─", w)
	}
	return Muted("┌─" + strings.Join(parts, "─┬─") + "─┐")
}

// TableBottom creates a bottom border for tables
func TableBottom(widths ...int) string {
	parts := make([]string, len(widths))
	for i, w := range widths {
		parts[i] = strings.Repeat("─", w)
	}
	return Muted("└─" + strings.Join(parts, "─┴─") + "─┘")
}

// PadRight pads a string to the right to reach the specified width
func PadRight(s string, width int) string {
	visLen := VisibleLen(s)
	if visLen >= width {
		return s
	}
	return s + strings.Repeat(" ", width-visLen)
}

// PadLeft pads a string to the left to reach the specified width
func PadLeft(s string, width int) string {
	visLen := VisibleLen(s)
	if visLen >= width {
		return s
	}
	return strings.Repeat(" ", width-visLen) + s
}

// StripANSI removes ANSI escape codes from a string
func StripANSI(s string) string {
	var result strings.Builder
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}

// VisibleLen calculates the visible display width of a string
// (excluding ANSI codes and accounting for wide characters like emojis)
func VisibleLen(s string) int {
	return runewidth.StringWidth(StripANSI(s))
}

// ProgressBar creates a colored progress bar
func ProgressBar(percent float64, width int) string {
	filled := int(percent / 100 * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	var color string
	switch {
	case percent >= 90:
		color = Red
	case percent >= 70:
		color = Yellow
	default:
		color = Green
	}

	bar := color + strings.Repeat(IconBar, filled) + Reset
	bar += Muted(strings.Repeat(IconBarLight, width-filled))
	return bar
}

// BoolToStatusColored returns a colored status string
func BoolToStatusColored(b bool) string {
	if b {
		return Success(IconCheck + " Yes")
	}
	return Danger(IconCross + " No")
}

// BoolToCheckbox returns a checkbox icon
func BoolToCheckbox(b bool) string {
	if b {
		return Success("☑")
	}
	return Muted("☐")
}

// FormatOutput returns the result in the requested format (json or table)
func FormatOutput(data any, tableFunc func() string, format string) string {
	if strings.ToLower(format) == FormatTable {
		return tableFunc()
	}
	resultJSON, _ := json.MarshalIndent(data, "", "  ")
	return string(resultJSON)
}

// UsageColor returns the appropriate color based on usage percentage
func UsageColor(percent float64) string {
	switch {
	case percent >= 90:
		return Red
	case percent >= 70:
		return Yellow
	default:
		return Green
	}
}
