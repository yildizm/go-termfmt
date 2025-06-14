package termfmt

// Formatter defines the interface for terminal output formatting
type Formatter interface {
	Format(data interface{}) ([]byte, error)
}

// TerminalOptions configures terminal formatting behavior
type TerminalOptions struct {
	Color     bool // Enable colored output
	Emoji     bool // Enable emoji output
	Width     int  // Terminal width for formatting
	Compact   bool // Use compact formatting
	ShowIcons bool // Show icons/symbols
}

const (
	// DefaultTerminalWidth is the standard terminal width
	DefaultTerminalWidth = 80
)

// DefaultOptions returns sensible default options
func DefaultOptions() *TerminalOptions {
	return &TerminalOptions{
		Color:     true,
		Emoji:     true,
		Width:     DefaultTerminalWidth,
		Compact:   false,
		ShowIcons: true,
	}
}
