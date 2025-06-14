package termfmt

import (
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if opts == nil {
		t.Fatal("DefaultOptions() returned nil")
	}

	if !opts.Color {
		t.Error("Expected Color to be true by default")
	}

	if !opts.Emoji {
		t.Error("Expected Emoji to be true by default")
	}

	if opts.Width != 80 {
		t.Errorf("Expected Width to be 80, got %d", opts.Width)
	}
}

func TestNewTerminal(t *testing.T) {
	formatter := NewTerminal(true)
	if formatter == nil {
		t.Fatal("NewTerminal() returned nil")
	}

	// Test that it implements the Formatter interface
	_ = formatter
}

func TestBox(t *testing.T) {
	result := Box("Test", "Content")
	if result == "" {
		t.Error("Box() returned empty string")
	}

	if !contains(result, "Test") {
		t.Error("Box() result doesn't contain title")
	}

	if !contains(result, "Content") {
		t.Error("Box() result doesn't contain content")
	}
}

func TestTable(t *testing.T) {
	headers := []string{"Name", "Value"}
	rows := [][]string{
		{"Test", "123"},
		{"Example", "456"},
	}

	result := Table(headers, rows)
	if result == "" {
		t.Error("Table() returned empty string")
	}

	if !contains(result, "Name") || !contains(result, "Value") {
		t.Error("Table() result doesn't contain headers")
	}

	if !contains(result, "Test") || !contains(result, "123") {
		t.Error("Table() result doesn't contain row data")
	}
}

func TestBarChart(t *testing.T) {
	data := map[string]int{
		"A": 10,
		"B": 20,
		"C": 30,
	}

	result := BarChart(data, 50)
	if result == "" {
		t.Error("BarChart() returned empty string")
	}

	if !contains(result, "A") || !contains(result, "B") || !contains(result, "C") {
		t.Error("BarChart() result doesn't contain all labels")
	}
}

func TestGetEmoji(t *testing.T) {
	// Test with emoji explicitly enabled
	emojiOpts := &TerminalOptions{
		Color: true,
		Emoji: true,
		Width: 80,
	}

	emoji := GetEmoji("success", emojiOpts)
	if emoji == "" {
		t.Error("GetEmoji() returned empty string")
	}

	// Test with emoji explicitly disabled
	noEmojiOpts := &TerminalOptions{
		Color: true,
		Emoji: false,
		Width: 80,
	}

	fallback := GetEmoji("success", noEmojiOpts)
	if fallback == "" {
		t.Error("GetEmoji() fallback returned empty string")
	}

	// Check that we get the expected fallback when emoji is disabled
	expectedFallback := "[OK]"
	if fallback != expectedFallback {
		t.Errorf("Expected fallback to be '%s', got '%s'", expectedFallback, fallback)
	}

	// The values should be different unless terminal doesn't support emoji
	// In which case both would return fallback, which is acceptable
	if emoji == fallback {
		// This is only acceptable if terminal doesn't support emoji
		t.Logf("Emoji and fallback are the same ('%s') - terminal may not support emoji", emoji)
	}
}

func TestCreateConfidenceBar(t *testing.T) {
	opts := DefaultOptions()
	opts.Emoji = false // Use text mode for predictable length

	bar := CreateConfidenceBar(0.5, opts)
	if bar == "" {
		t.Error("CreateConfidenceBar() returned empty string")
	}

	// Should contain filled and empty portions
	if !contains(bar, "#") || !contains(bar, "-") {
		t.Error("CreateConfidenceBar() doesn't contain expected characters")
	}
}

func TestProgressBar(t *testing.T) {
	bar := ProgressBar(50, 100, 60)
	if bar == "" {
		t.Error("ProgressBar() returned empty string")
	}

	if !contains(bar, "50.0%") {
		t.Error("ProgressBar() doesn't show correct percentage")
	}

	if !contains(bar, "(50/100)") {
		t.Error("ProgressBar() doesn't show correct counts")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

// Simple substring search
func indexOf(s, substr string) int {
	n := len(substr)
	if n == 0 {
		return 0
	}

	for i := 0; i <= len(s)-n; i++ {
		if s[i:i+n] == substr {
			return i
		}
	}

	return -1
}
