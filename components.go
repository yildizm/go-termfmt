package termfmt

import (
	"fmt"
	"strings"
)

const (
	// UI spacing constants
	borderPadding     = 2   // Padding around box borders
	labelSpacing      = 10  // Space reserved for labels in bar charts
	progressSpacing   = 20  // Space reserved for progress bar metadata
	tableRowPadding   = 2   // Extra padding for table rows
	percentMultiplier = 100 // Multiplier for percentage calculations
)

// Box creates a bordered box around content with an optional title
func Box(title, content string) string {
	return BoxWithOptions(title, content, DefaultOptions())
}

// BoxWithOptions creates a bordered box with custom options
func BoxWithOptions(title, content string, opts *TerminalOptions) string {
	if title == "" {
		return simpleBox(content)
	}

	return titledBox(title, content)
}

// simpleBox creates a simple bordered box
func simpleBox(content string) string {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return ""
	}

	// Find the maximum line length
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	var b strings.Builder

	// Top border
	b.WriteString("┌" + strings.Repeat("─", maxLen+borderPadding) + "┐\n")

	// Content lines
	for _, line := range lines {
		padding := maxLen - len(line)
		b.WriteString("│ " + line + strings.Repeat(" ", padding) + " │\n")
	}

	// Bottom border
	b.WriteString("└" + strings.Repeat("─", maxLen+borderPadding) + "┘")

	return b.String()
}

// titledBox creates a box with a title
func titledBox(title, content string) string {
	titleLen := len(title)
	lines := strings.Split(content, "\n")

	// Find the maximum line length
	maxLen := titleLen + borderPadding // Title + padding
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	var b strings.Builder

	// Top border with title
	b.WriteString("╔" + strings.Repeat("═", maxLen+borderPadding) + "╗\n")
	titlePadding := maxLen - titleLen
	b.WriteString("║ " + title + strings.Repeat(" ", titlePadding) + " ║\n")
	b.WriteString("╠" + strings.Repeat("═", maxLen+borderPadding) + "╣\n")

	// Content lines
	for _, line := range lines {
		padding := maxLen - len(line)
		b.WriteString("║ " + line + strings.Repeat(" ", padding) + " ║\n")
	}

	// Bottom border
	b.WriteString("╚" + strings.Repeat("═", maxLen+borderPadding) + "╝")

	return b.String()
}

// Table creates a formatted table from headers and rows
func Table(headers []string, rows [][]string) string {
	return TableWithOptions(headers, rows, DefaultOptions())
}

// TableWithOptions creates a formatted table with custom options
func TableWithOptions(headers []string, rows [][]string, opts *TerminalOptions) string {
	if len(headers) == 0 {
		return ""
	}

	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	var b strings.Builder

	// Header row
	b.WriteString("│")

	for i, header := range headers {
		padding := colWidths[i] - len(header)
		b.WriteString(" " + header + strings.Repeat(" ", padding) + " │")
	}

	b.WriteString("\n")

	// Separator
	b.WriteString("├")

	for i, width := range colWidths {
		b.WriteString(strings.Repeat("─", width+tableRowPadding))

		if i < len(colWidths)-1 {
			b.WriteString("┼")
		}
	}

	b.WriteString("┤\n")

	// Data rows
	for _, row := range rows {
		b.WriteString("│")

		for i, cell := range row {
			if i < len(colWidths) {
				padding := colWidths[i] - len(cell)
				b.WriteString(" " + cell + strings.Repeat(" ", padding) + " │")
			}
		}

		b.WriteString("\n")
	}

	return strings.TrimRight(b.String(), "\n")
}

// BarChart creates a horizontal bar chart from data
func BarChart(data map[string]int, width int) string {
	return BarChartWithOptions(data, width, DefaultOptions())
}

// BarChartWithOptions creates a horizontal bar chart with custom options
func BarChartWithOptions(data map[string]int, width int, opts *TerminalOptions) string {
	if len(data) == 0 {
		return ""
	}

	// Find max value for scaling
	maxValue := 0
	maxLabelLen := 0

	for label, value := range data {
		if value > maxValue {
			maxValue = value
		}

		if len(label) > maxLabelLen {
			maxLabelLen = len(label)
		}
	}

	if maxValue == 0 {
		return ""
	}

	var b strings.Builder

	barWidth := width - maxLabelLen - labelSpacing // Leave space for label and value

	for label, value := range data {
		// Label (right-padded)
		labelPadding := maxLabelLen - len(label)
		b.WriteString(label + strings.Repeat(" ", labelPadding))

		// Bar
		barLength := int(float64(value) / float64(maxValue) * float64(barWidth))

		b.WriteString(" │")

		if opts.Emoji {
			b.WriteString(strings.Repeat("█", barLength))
			b.WriteString(strings.Repeat("░", barWidth-barLength))
		} else {
			b.WriteString(strings.Repeat("#", barLength))
			b.WriteString(strings.Repeat("-", barWidth-barLength))
		}

		// Value
		b.WriteString(fmt.Sprintf("│ %d\n", value))
	}

	return strings.TrimRight(b.String(), "\n")
}

// TreeView creates a tree-style view with prefix indicators
func TreeView(items []TreeItem) string {
	return TreeViewWithOptions(items, DefaultOptions())
}

// TreeItem represents an item in a tree view
type TreeItem struct {
	Label    string
	Value    string
	Children []TreeItem
	Last     bool // Whether this is the last item in its group
}

// TreeViewWithOptions creates a tree-style view with custom options
func TreeViewWithOptions(items []TreeItem, opts *TerminalOptions) string {
	var b strings.Builder

	renderTreeItems(&b, items, "", opts)

	return strings.TrimRight(b.String(), "\n")
}

// renderTreeItems recursively renders tree items
func renderTreeItems(b *strings.Builder, items []TreeItem, prefix string, _ *TerminalOptions) {
	for i, item := range items {
		isLast := i == len(items)-1

		// Choose prefix based on position
		var itemPrefix, childPrefix string
		if isLast {
			itemPrefix = "└─ "
			childPrefix = "   "
		} else {
			itemPrefix = "├─ "
			childPrefix = "│  "
		}

		// Write the item
		b.WriteString(prefix + itemPrefix + item.Label)

		if item.Value != "" {
			b.WriteString(": " + item.Value)
		}

		b.WriteString("\n")

		// Render children
		if len(item.Children) > 0 {
			renderTreeItems(b, item.Children, prefix+childPrefix, nil)
		}
	}
}

// ProgressBar creates a progress bar
func ProgressBar(current, total, width int) string {
	return ProgressBarWithOptions(current, total, width, DefaultOptions())
}

// ProgressBarWithOptions creates a progress bar with custom options
func ProgressBarWithOptions(current, total, width int, opts *TerminalOptions) string {
	if total <= 0 {
		return ""
	}

	percentage := float64(current) / float64(total)
	if percentage > 1.0 {
		percentage = 1.0
	}

	barWidth := width - progressSpacing // Leave space for percentage and brackets
	filledWidth := int(percentage * float64(barWidth))

	var b strings.Builder

	b.WriteString("[")

	if opts.Emoji {
		b.WriteString(strings.Repeat("█", filledWidth))
		b.WriteString(strings.Repeat("░", barWidth-filledWidth))
	} else {
		b.WriteString(strings.Repeat("#", filledWidth))
		b.WriteString(strings.Repeat("-", barWidth-filledWidth))
	}

	b.WriteString(fmt.Sprintf("] %.1f%% (%d/%d)", percentage*percentMultiplier, current, total))

	return b.String()
}
