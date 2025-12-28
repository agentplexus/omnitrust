package inspector

import (
	"strings"
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    uint64
		expected string
	}{
		{"zero bytes", 0, "0 bytes"},
		{"bytes", 500, "500 bytes"},
		{"one KB", 1024, "1.00 KB"},
		{"KB", 2048, "2.00 KB"},
		{"one MB", 1024 * 1024, "1.00 MB"},
		{"MB", 5 * 1024 * 1024, "5.00 MB"},
		{"one GB", 1024 * 1024 * 1024, "1.00 GB"},
		{"GB", 8 * 1024 * 1024 * 1024, "8.00 GB"},
		{"large GB", 16 * 1024 * 1024 * 1024, "16.00 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("FormatBytes(%d) = %q, want %q", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no ANSI", "hello world", "hello world"},
		{"single color", "\033[31mred\033[0m", "red"},
		{"bold", "\033[1mbold\033[0m", "bold"},
		{"multiple codes", "\033[1m\033[32mgreen bold\033[0m", "green bold"},
		{"mixed content", "before \033[31mred\033[0m after", "before red after"},
		{"empty string", "", ""},
		{"only reset", "\033[0m", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripANSI(tt.input)
			if result != tt.expected {
				t.Errorf("StripANSI(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestVisibleLen(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"empty", "", 0},
		{"ascii", "hello", 5},
		{"with ANSI", "\033[31mhello\033[0m", 5},
		{"unicode", "世界", 4}, // 2 wide chars
		{"emoji", "👍", 2},    // emoji is 2 wide
		{"mixed", "hi 👋", 5}, // 2 + 1 + 2
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VisibleLen(tt.input)
			if result != tt.expected {
				t.Errorf("VisibleLen(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		expected string
	}{
		{"no padding needed", "hello", 5, "hello"},
		{"pad short string", "hi", 5, "hi   "},
		{"already longer", "hello world", 5, "hello world"},
		{"empty string", "", 3, "   "},
		{"with ANSI", "\033[31mhi\033[0m", 5, "\033[31mhi\033[0m   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PadRight(tt.input, tt.width)
			if result != tt.expected {
				t.Errorf("PadRight(%q, %d) = %q, want %q", tt.input, tt.width, result, tt.expected)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		expected string
	}{
		{"no padding needed", "hello", 5, "hello"},
		{"pad short string", "hi", 5, "   hi"},
		{"already longer", "hello world", 5, "hello world"},
		{"empty string", "", 3, "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PadLeft(tt.input, tt.width)
			if result != tt.expected {
				t.Errorf("PadLeft(%q, %d) = %q, want %q", tt.input, tt.width, result, tt.expected)
			}
		})
	}
}

func TestColorize(t *testing.T) {
	result := Colorize(Red, "error")
	if !strings.HasPrefix(result, Red) {
		t.Errorf("Colorize should start with color code")
	}
	if !strings.HasSuffix(result, Reset) {
		t.Errorf("Colorize should end with reset code")
	}
	if !strings.Contains(result, "error") {
		t.Errorf("Colorize should contain the text")
	}
}

func TestFormattingFunctions(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string) string
		text string
	}{
		{"BoldText", BoldText, "bold"},
		{"DimText", DimText, "dim"},
		{"Header", Header, "header"},
		{"Success", Success, "success"},
		{"Warning", Warning, "warning"},
		{"Danger", Danger, "danger"},
		{"Info", Info, "info"},
		{"Muted", Muted, "muted"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.text)
			// Should contain the original text
			stripped := StripANSI(result)
			if stripped != tt.text {
				t.Errorf("%s(%q) stripped = %q, want %q", tt.name, tt.text, stripped, tt.text)
			}
			// Should have ANSI codes
			if result == tt.text {
				t.Errorf("%s(%q) should contain ANSI codes", tt.name, tt.text)
			}
		})
	}
}

func TestProgressBar(t *testing.T) {
	tests := []struct {
		name        string
		percent     float64
		width       int
		expectGreen bool
		expectRed   bool
	}{
		{"low usage", 30.0, 20, true, false},
		{"medium usage", 75.0, 20, false, false}, // yellow
		{"high usage", 95.0, 20, false, true},
		{"zero", 0.0, 10, true, false},
		{"full", 100.0, 10, false, true},
		{"over 100", 150.0, 10, false, true},
		{"negative", -10.0, 10, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProgressBar(tt.percent, tt.width)
			if tt.expectGreen && !strings.Contains(result, Green) {
				t.Errorf("ProgressBar(%.1f) should contain green", tt.percent)
			}
			if tt.expectRed && !strings.Contains(result, Red) {
				t.Errorf("ProgressBar(%.1f) should contain red", tt.percent)
			}
		})
	}
}

func TestUsageColor(t *testing.T) {
	tests := []struct {
		percent  float64
		expected string
	}{
		{0, Green},
		{50, Green},
		{69, Green},
		{70, Yellow},
		{85, Yellow},
		{89, Yellow},
		{90, Red},
		{95, Red},
		{100, Red},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := UsageColor(tt.percent)
			if result != tt.expected {
				t.Errorf("UsageColor(%.0f) = %q, want %q", tt.percent, result, tt.expected)
			}
		})
	}
}

func TestBoolToStatusColored(t *testing.T) {
	trueResult := BoolToStatusColored(true)
	if !strings.Contains(trueResult, IconCheck) {
		t.Error("BoolToStatusColored(true) should contain check icon")
	}
	if !strings.Contains(trueResult, "Yes") {
		t.Error("BoolToStatusColored(true) should contain 'Yes'")
	}

	falseResult := BoolToStatusColored(false)
	if !strings.Contains(falseResult, IconCross) {
		t.Error("BoolToStatusColored(false) should contain cross icon")
	}
	if !strings.Contains(falseResult, "No") {
		t.Error("BoolToStatusColored(false) should contain 'No'")
	}
}

func TestBoolToCheckbox(t *testing.T) {
	trueResult := BoolToCheckbox(true)
	if !strings.Contains(trueResult, "☑") {
		t.Error("BoolToCheckbox(true) should contain checked box")
	}

	falseResult := BoolToCheckbox(false)
	if !strings.Contains(falseResult, "☐") {
		t.Error("BoolToCheckbox(false) should contain unchecked box")
	}
}

func TestTableFunctions(t *testing.T) {
	widths := []int{10, 15, 20}

	top := TableTop(widths...)
	if !strings.Contains(top, "┌") {
		t.Error("TableTop should contain top-left corner")
	}
	if !strings.Contains(top, "┬") {
		t.Error("TableTop should contain top separator")
	}
	if !strings.Contains(top, "┐") {
		t.Error("TableTop should contain top-right corner")
	}

	sep := TableSeparator(widths...)
	if !strings.Contains(sep, "├") {
		t.Error("TableSeparator should contain left junction")
	}
	if !strings.Contains(sep, "┼") {
		t.Error("TableSeparator should contain cross junction")
	}
	if !strings.Contains(sep, "┤") {
		t.Error("TableSeparator should contain right junction")
	}

	bottom := TableBottom(widths...)
	if !strings.Contains(bottom, "└") {
		t.Error("TableBottom should contain bottom-left corner")
	}
	if !strings.Contains(bottom, "┴") {
		t.Error("TableBottom should contain bottom separator")
	}
	if !strings.Contains(bottom, "┘") {
		t.Error("TableBottom should contain bottom-right corner")
	}
}

func TestTableRow(t *testing.T) {
	result := TableRow("col1", "col2", "col3")
	if !strings.HasPrefix(result, "│") {
		t.Error("TableRow should start with vertical bar")
	}
	if !strings.HasSuffix(result, "│") {
		t.Error("TableRow should end with vertical bar")
	}
	if !strings.Contains(result, "col1") {
		t.Error("TableRow should contain column content")
	}
}

func TestFormatOutput(t *testing.T) {
	data := map[string]string{"key": "value"}
	tableFunc := func() string { return "table output" }

	// Test JSON format (default)
	jsonResult := FormatOutput(data, tableFunc, "json")
	if !strings.Contains(jsonResult, "key") || !strings.Contains(jsonResult, "value") {
		t.Error("FormatOutput with json format should return JSON")
	}

	// Test table format
	tableResult := FormatOutput(data, tableFunc, "table")
	if tableResult != "table output" {
		t.Errorf("FormatOutput with table format = %q, want %q", tableResult, "table output")
	}

	// Test case insensitivity
	tableResult2 := FormatOutput(data, tableFunc, "TABLE")
	if tableResult2 != "table output" {
		t.Error("FormatOutput should be case insensitive for format")
	}
}

func TestConstants(t *testing.T) {
	// Verify format constants
	if FormatJSON != "json" {
		t.Errorf("FormatJSON = %q, want %q", FormatJSON, "json")
	}
	if FormatTable != "table" {
		t.Errorf("FormatTable = %q, want %q", FormatTable, "table")
	}

	// Verify ANSI codes are non-empty
	codes := []struct {
		name string
		code string
	}{
		{"Reset", Reset},
		{"Bold", Bold},
		{"Red", Red},
		{"Green", Green},
		{"Yellow", Yellow},
		{"Blue", Blue},
		{"Cyan", Cyan},
	}

	for _, c := range codes {
		if c.code == "" {
			t.Errorf("%s should not be empty", c.name)
		}
		if !strings.HasPrefix(c.code, "\033[") {
			t.Errorf("%s should start with escape sequence", c.name)
		}
	}

	// Verify icons are non-empty
	icons := []struct {
		name string
		icon string
	}{
		{"IconCheck", IconCheck},
		{"IconCross", IconCross},
		{"IconBar", IconBar},
		{"IconBarLight", IconBarLight},
	}

	for _, i := range icons {
		if i.icon == "" {
			t.Errorf("%s should not be empty", i.name)
		}
	}
}
