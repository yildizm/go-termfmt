package termfmt

import (
	"os"
	"strings"
)

const (
	// ConfidenceBarLength is the standard length for confidence bars
	ConfidenceBarLength = 10
)

// Color codes for terminal output
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"
	Italic = "\033[3m"

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
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

// getEmojiMap returns emoji mappings with fallbacks
func getEmojiMap() map[string][2]string {
	return map[string][2]string{
		// [emoji, fallback]
		"error":            {"❌", "[ERR]"},
		"warning":          {"⚠️", "[WRN]"},
		"info":             {"ℹ️", "[INF]"},
		"success":          {"✅", "[OK]"},
		"insight":          {"💡", "[INS]"},
		"pattern":          {"🔍", "[PAT]"},
		"statistics":       {"📊", "[STATS]"},
		"recommendations":  {"📋", "[REC]"},
		"rocket":           {"🚀", "[LOG]"},
		"error_pattern":    {"🔴", "[ERR]"},
		"anomaly_pattern":  {"🟡", "[ANO]"},
		"perf_pattern":     {"⚠️", "[PRF]"},
		"security_pattern": {"🔒", "[SEC]"},
		"help":             {"❓", "[?]"},
		"target":           {"🎯", "[>]"},
		"brain":            {"🧠", "[AI]"},
		"tag":              {"🏷️", "[TAG]"},
		"scale":            {"⚖️", "[BAL]"},
		"door":             {"🚪", "[EXIT]"},
		"number":           {"🔢", "[#]"},
		"patterns":         {"🔍", "[PAT]"},
		"insights":         {"💡", "[INS]"},
		"clock":            {"⏰", "[TIME]"},
		"chart":            {"📊", "[CHART]"},
		"list":             {"📝", "[LIST]"},
	}
}

// ColorProfile represents a color scheme
type ColorProfile struct {
	Error   string
	Warning string
	Info    string
	Success string
	Accent  string
	Muted   string
}

// DefaultColorProfile returns the default color scheme
func DefaultColorProfile() *ColorProfile {
	return &ColorProfile{
		Error:   Red,
		Warning: Yellow,
		Info:    Blue,
		Success: Green,
		Accent:  Cyan,
		Muted:   BrightBlack,
	}
}

// HighContrastColorProfile returns a high contrast color scheme
func HighContrastColorProfile() *ColorProfile {
	return &ColorProfile{
		Error:   BrightRed,
		Warning: BrightYellow,
		Info:    BrightBlue,
		Success: BrightGreen,
		Accent:  BrightCyan,
		Muted:   White,
	}
}

// Colorize applies color to text if color is enabled
func Colorize(text, color string, opts *TerminalOptions) string {
	if !opts.Color || !supportsColor() {
		return text
	}

	return color + text + Reset
}

// ColorizeWithProfile applies color using a color profile
func ColorizeWithProfile(
	text, colorType string,
	profile *ColorProfile,
	opts *TerminalOptions,
) string {
	if !opts.Color || !supportsColor() {
		return text
	}

	var color string

	switch colorType {
	case "error":
		color = profile.Error
	case "warning":
		color = profile.Warning
	case "info":
		color = profile.Info
	case "success":
		color = profile.Success
	case "accent":
		color = profile.Accent
	case "muted":
		color = profile.Muted
	default:
		return text
	}

	return color + text + Reset
}

// GetEmoji returns emoji or fallback based on options
func GetEmoji(key string, opts *TerminalOptions) string {
	if mapping, exists := getEmojiMap()[key]; exists {
		// Respect user preference first - if they explicitly disable emoji, use fallback
		if !opts.Emoji {
			return mapping[1] // fallback
		}

		// If user wants emoji but terminal doesn't support it, graceful fallback
		if !supportsEmoji() {
			return mapping[1] // fallback
		}

		return mapping[0] // emoji
	}

	return "[?]" // unknown key
}

// GetSymbol is an alias for GetEmoji for compatibility
func GetSymbol(key string, opts *TerminalOptions) string {
	return GetEmoji(key, opts)
}

// CreateConfidenceBar creates ASCII confidence bar
func CreateConfidenceBar(confidence float64, opts *TerminalOptions) string {
	barLength := int(confidence * ConfidenceBarLength)

	var filled, empty rune
	if opts.Emoji && supportsEmoji() {
		filled = '█'
		empty = '░'
	} else {
		filled = '#'
		empty = '-'
	}

	var b strings.Builder
	if !opts.Emoji {
		b.WriteString("[")
	}

	for i := range ConfidenceBarLength {
		if i < barLength {
			b.WriteRune(filled)
		} else {
			b.WriteRune(empty)
		}
	}

	if !opts.Emoji {
		b.WriteString("]")
	}

	return b.String()
}

// supportsColor checks if the terminal supports color output
func supportsColor() bool {
	// Check common environment variables
	term := os.Getenv("TERM")
	if term == "" {
		return false
	}

	// Check for known color-supporting terminals
	colorTerms := []string{
		"xterm", "xterm-256color", "xterm-color",
		"screen", "screen-256color",
		"tmux", "tmux-256color",
		"rxvt", "rxvt-unicode",
		"linux", "ansi",
	}

	for _, colorTerm := range colorTerms {
		if strings.Contains(term, colorTerm) {
			return true
		}
	}

	// Check for NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check for FORCE_COLOR environment variable
	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}

	return false
}

// supportsEmoji checks if the terminal likely supports emoji
func supportsEmoji() bool {
	// Similar to color support, but more conservative
	term := os.Getenv("TERM")
	if strings.Contains(term, "256color") {
		return true
	}

	// Check for modern terminals
	if os.Getenv("TERM_PROGRAM") != "" {
		return true
	}

	return false
}

// Stylize applies multiple styles to text
func Stylize(text string, styles []string, opts *TerminalOptions) string {
	if !opts.Color || !supportsColor() {
		return text
	}

	var codes []string

	for _, style := range styles {
		switch style {
		case "bold":
			codes = append(codes, Bold)
		case "dim":
			codes = append(codes, Dim)
		case "italic":
			codes = append(codes, Italic)
		case "red":
			codes = append(codes, Red)
		case "green":
			codes = append(codes, Green)
		case "yellow":
			codes = append(codes, Yellow)
		case "blue":
			codes = append(codes, Blue)
		case "magenta":
			codes = append(codes, Magenta)
		case "cyan":
			codes = append(codes, Cyan)
		case "white":
			codes = append(codes, White)
		}
	}

	if len(codes) == 0 {
		return text
	}

	return strings.Join(codes, "") + text + Reset
}

// Header creates a styled header with optional color
func Header(title string, opts *TerminalOptions) string {
	if opts.Color && supportsColor() {
		return Stylize(title, []string{"bold", "cyan"}, opts)
	}

	return title
}

// Subtitle creates a styled subtitle
func Subtitle(title string, opts *TerminalOptions) string {
	if opts.Color && supportsColor() {
		return Stylize(title, []string{"bold"}, opts)
	}

	return title
}

// Muted creates muted/dimmed text
func Muted(text string, opts *TerminalOptions) string {
	if opts.Color && supportsColor() {
		return Stylize(text, []string{"dim"}, opts)
	}

	return text
}

// Success creates success-styled text
func Success(text string, opts *TerminalOptions) string {
	if opts.Color && supportsColor() {
		return Stylize(text, []string{"green"}, opts)
	}

	return text
}

// Warning creates warning-styled text
func Warning(text string, opts *TerminalOptions) string {
	if opts.Color && supportsColor() {
		return Stylize(text, []string{"yellow"}, opts)
	}

	return text
}

// Error creates error-styled text
func Error(text string, opts *TerminalOptions) string {
	if opts.Color && supportsColor() {
		return Stylize(text, []string{"red"}, opts)
	}

	return text
}

// Info creates info-styled text
func Info(text string, opts *TerminalOptions) string {
	if opts.Color && supportsColor() {
		return Stylize(text, []string{"blue"}, opts)
	}

	return text
}
