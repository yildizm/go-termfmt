package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yildizm/go-termfmt"
)

// SampleData represents some sample data to format
type SampleData struct {
	Name        string            `json:"name"`
	Value       int               `json:"value"`
	Timestamp   time.Time         `json:"timestamp"`
	Active      bool              `json:"active"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
	Performance *PerformanceData  `json:"performance,omitempty"`
}

// PerformanceData represents nested performance metrics
type PerformanceData struct {
	CPU    float64 `json:"cpu"`
	Memory int64   `json:"memory"`
	Disk   float64 `json:"disk"`
}

func main() {
	fmt.Println("🎨 go-termfmt Library Demo")
	fmt.Println("==========================\n")

	// Test different color and emoji options
	demoBasicComponents()
	demoAdvancedComponents()
	demoStructFormatting()
	demoDataFormatting()
	demoColorAndStyling()
}

func demoBasicComponents() {
	fmt.Println("📦 Basic Components")
	fmt.Println("-------------------\n")

	// Basic box
	opts := termfmt.DefaultOptions()
	boxContent := termfmt.Box("Simple Box", "This is a simple box with some content inside.")
	fmt.Println(boxContent)
	fmt.Println()

	// Titled box
	titledBox := termfmt.BoxWithOptions("Configuration",
		"Server: localhost:8080\nDatabase: postgresql://...\nCacheSize: 1024MB", opts)
	fmt.Println(titledBox)
	fmt.Println()

	// Progress bar
	progressBar := termfmt.ProgressBar(75, 100, 60)
	fmt.Println("Progress:", progressBar)
	fmt.Println()
}

func demoAdvancedComponents() {
	fmt.Println("📊 Advanced Components")
	fmt.Println("----------------------\n")

	// Table
	headers := []string{"Service", "Status", "CPU", "Memory"}
	rows := [][]string{
		{"api-server", "Running", "45%", "512MB"},
		{"database", "Running", "23%", "2.1GB"},
		{"cache", "Stopped", "0%", "0MB"},
		{"worker", "Running", "67%", "256MB"},
	}

	table := termfmt.Table(headers, rows)
	fmt.Println("Service Status Table:")
	fmt.Println(table)
	fmt.Println()

	// Bar chart
	chartData := map[string]int{
		"Errors":   45,
		"Warnings": 123,
		"Info":     456,
		"Debug":    234,
	}

	chart := termfmt.BarChart(chartData, 50)
	fmt.Println("Log Level Distribution:")
	fmt.Println(chart)
	fmt.Println()

	// Tree view
	treeItems := []termfmt.TreeItem{
		{Label: "Application", Value: "", Children: []termfmt.TreeItem{
			{Label: "API Server", Value: "running"},
			{Label: "Database", Value: "connected"},
			{Label: "Cache", Value: "redis://localhost:6379"},
		}},
		{Label: "Metrics", Value: "", Children: []termfmt.TreeItem{
			{Label: "Requests/sec", Value: "1,234"},
			{Label: "Error rate", Value: "0.05%"},
			{Label: "Avg response", Value: "45ms"},
		}},
	}

	tree := termfmt.TreeView(treeItems)
	fmt.Println("System Overview:")
	fmt.Println(tree)
	fmt.Println()
}

func demoStructFormatting() {
	fmt.Println("🏗️  Struct Formatting")
	fmt.Println("---------------------\n")

	// Sample data
	sampleData := &SampleData{
		Name:      "web-server-01",
		Value:     42,
		Timestamp: time.Now(),
		Active:    true,
		Tags:      []string{"production", "web", "critical"},
		Metadata: map[string]string{
			"region":      "us-west-2",
			"environment": "production",
			"version":     "1.2.3",
		},
		Performance: &PerformanceData{
			CPU:    78.5,
			Memory: 2048576,
			Disk:   45.2,
		},
	}

	// Format using the terminal formatter
	formatter := termfmt.NewTerminal(true)
	output, err := formatter.Format(sampleData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))
	fmt.Println()
}

func demoDataFormatting() {
	fmt.Println("📋 Data Formatting")
	fmt.Println("------------------\n")

	// Summary box
	summaryData := map[string]interface{}{
		"Total Entries": 1543,
		"Errors":        12,
		"Warnings":      45,
		"Success Rate":  "98.2%",
		"Uptime":        "99.9%",
	}

	summary := termfmt.Summary("System Summary", summaryData, termfmt.DefaultOptions())
	fmt.Println(summary)
	fmt.Println()

	// Confidence bar
	opts := termfmt.DefaultOptions()
	fmt.Println("Confidence Levels:")
	for i, label := range []string{"Low", "Medium", "High", "Very High"} {
		confidence := float64(i+1) / 4.0
		bar := termfmt.CreateConfidenceBar(confidence, opts)
		fmt.Printf("%-10s %s %.0f%%\n", label+":", bar, confidence*100)
	}
	fmt.Println()
}

func demoColorAndStyling() {
	fmt.Println("🎨 Color and Styling")
	fmt.Println("-------------------\n")

	opts := termfmt.DefaultOptions()

	// Test different text styles
	fmt.Println("Text Styles:")
	fmt.Println("  Header:  " + termfmt.Header("Important Header", opts))
	fmt.Println("  Success: " + termfmt.Success("Operation completed", opts))
	fmt.Println("  Warning: " + termfmt.Warning("Deprecated feature", opts))
	fmt.Println("  Error:   " + termfmt.Error("Something went wrong", opts))
	fmt.Println("  Info:    " + termfmt.Info("Informational message", opts))
	fmt.Println("  Muted:   " + termfmt.Muted("Less important info", opts))
	fmt.Println()

	// Test emojis vs fallbacks
	fmt.Println("Emoji Support:")
	emojiOpts := termfmt.DefaultOptions()
	noEmojiOpts := &termfmt.TerminalOptions{Color: true, Emoji: false}

	symbols := []string{"error", "warning", "success", "info", "statistics", "rocket"}
	for _, symbol := range symbols {
		emoji := termfmt.GetEmoji(symbol, emojiOpts)
		fallback := termfmt.GetEmoji(symbol, noEmojiOpts)
		fmt.Printf("  %-12s: %s (with emoji) | %s (fallback)\n", symbol, emoji, fallback)
	}
	fmt.Println()

	// Test color vs no-color
	fmt.Println("Color Support:")
	colorOpts := &termfmt.TerminalOptions{Color: true, Emoji: true}
	noColorOpts := &termfmt.TerminalOptions{Color: false, Emoji: true}

	testText := "Colored Text"
	fmt.Println("  With color:    " + termfmt.Error(testText, colorOpts))
	fmt.Println("  Without color: " + termfmt.Error(testText, noColorOpts))
	fmt.Println()
}
