# go-termfmt

[![Go Reference](https://pkg.go.dev/badge/github.com/yildizm/go-termfmt.svg)](https://pkg.go.dev/github.com/yildizm/go-termfmt)
[![CI](https://github.com/yildizm/go-termfmt/workflows/CI/badge.svg)](https://github.com/yildizm/go-termfmt/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/yildizm/go-termfmt)](https://goreportcard.com/report/github.com/yildizm/go-termfmt)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/yildizm/go-termfmt)](https://github.com/yildizm/go-termfmt)
[![Release](https://img.shields.io/github/v/release/yildizm/go-termfmt)](https://github.com/yildizm/go-termfmt/releases)
[![Issues](https://img.shields.io/github/issues/yildizm/go-termfmt)](https://github.com/yildizm/go-termfmt/issues)
[![Stars](https://img.shields.io/github/stars/yildizm/go-termfmt)](https://github.com/yildizm/go-termfmt/stargazers)
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/yildizm/go-termfmt/graphs/commit-activity)

A lightweight Go library for beautiful terminal formatting with zero dependencies.

## Features

- 📦 **Terminal Components**: Boxes, tables, bar charts, progress bars, tree views
- 🎨 **Color Support**: Automatic color detection with fallbacks
- 🔤 **Emoji Support**: Unicode emojis with text fallbacks
- 🏗️ **Struct Formatting**: Automatic formatting of Go structs using reflection
- ⚙️ **Configurable**: Flexible options for colors, emojis, and styling
- 🚀 **Zero Dependencies**: Pure Go implementation

## Installation

```bash
go get github.com/yildizm/go-termfmt
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/yildizm/go-termfmt"
)

func main() {
    // Create a simple box
    box := termfmt.Box("Hello", "Welcome to go-termfmt!")
    fmt.Println(box)
    
    // Create a table
    headers := []string{"Name", "Status", "Value"}
    rows := [][]string{
        {"Server", "Running", "100%"},
        {"Database", "Connected", "99%"},
    }
    table := termfmt.Table(headers, rows)
    fmt.Println(table)
    
    // Format any Go struct
    type Config struct {
        Host string
        Port int
        SSL  bool
    }
    
    config := Config{Host: "localhost", Port: 8080, SSL: true}
    formatter := termfmt.NewTerminal(true) // with colors
    output, _ := formatter.Format(config)
    fmt.Print(string(output))
}
```

## Components

### Boxes

Create bordered boxes with optional titles:

```go
// Simple box
simple := termfmt.Box("", "Simple content")

// Titled box
titled := termfmt.Box("Configuration", "Host: localhost\nPort: 8080")
```

### Tables

Create formatted tables from headers and rows:

```go
headers := []string{"Service", "Status", "CPU"}
rows := [][]string{
    {"api", "Running", "45%"},
    {"db", "Connected", "23%"},
}
table := termfmt.Table(headers, rows)
```

### Bar Charts

Create horizontal bar charts from data:

```go
data := map[string]int{
    "Errors":   45,
    "Warnings": 123,
    "Info":     456,
}
chart := termfmt.BarChart(data, 50) // width of 50 chars
```

### Tree Views

Create hierarchical tree displays:

```go
items := []termfmt.TreeItem{
    {Label: "Root", Children: []termfmt.TreeItem{
        {Label: "Child 1", Value: "value1"},
        {Label: "Child 2", Value: "value2"},
    }},
}
tree := termfmt.TreeView(items)
```

### Progress Bars

Create progress indicators:

```go
progress := termfmt.ProgressBar(75, 100, 60) // 75/100, width 60
```

## Styling

### Colors

The library automatically detects terminal color support:

```go
opts := termfmt.DefaultOptions()
opts.Color = true

// Style text
header := termfmt.Header("Important", opts)
success := termfmt.Success("Done!", opts)
warning := termfmt.Warning("Careful", opts)
error := termfmt.Error("Failed", opts)
```

### Emojis

Unicode emojis with automatic fallbacks:

```go
opts := termfmt.DefaultOptions()
opts.Emoji = true

// Get emojis or fallbacks
emoji := termfmt.GetEmoji("success", opts) // ✅ or [OK]
```

## Configuration

Customize formatting behavior with options:

```go
opts := &termfmt.TerminalOptions{
    Color:     true,  // Enable colors
    Emoji:     true,  // Enable emojis
    Width:     80,    // Terminal width
    Compact:   false, // Use compact formatting
    ShowIcons: true,  // Show icons/symbols
}

formatter := termfmt.NewTerminalWithOptions(opts)
```

## Struct Formatting

The library can automatically format any Go struct:

```go
type User struct {
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    Age      int       `json:"age"`
    Active   bool      `json:"active"`
    Created  time.Time `json:"created"`
}

user := User{
    Name:    "John Doe",
    Email:   "john@example.com",
    Age:     30,
    Active:  true,
    Created: time.Now(),
}

formatter := termfmt.NewTerminal(true)
output, _ := formatter.Format(user)
fmt.Print(string(output))
```

Output:
```
User
────

├─ Name: "John Doe"
├─ Email: "john@example.com"
├─ Age: 30
├─ Active: true
└─ Created: 2023-12-07 10:30:45
```

## Examples

See the [examples](examples/) directory for comprehensive usage examples:

```bash
cd examples/demo
go run main.go
```

## Environment Detection

The library automatically detects terminal capabilities:

- **Color Support**: Checks `TERM`, `NO_COLOR`, `FORCE_COLOR` environment variables
- **Emoji Support**: Conservative detection based on terminal type
- **Fallbacks**: Always provides text-based alternatives

## API Reference

### Core Types

```go
type Formatter interface {
    Format(data interface{}) ([]byte, error)
}

type TerminalOptions struct {
    Color      bool
    Emoji      bool  
    Width      int
    Compact    bool
    ShowIcons  bool
}
```

### Component Functions

```go
// Boxes
func Box(title, content string) string
func BoxWithOptions(title, content string, opts *TerminalOptions) string

// Tables  
func Table(headers []string, rows [][]string) string
func TableWithOptions(headers []string, rows [][]string, opts *TerminalOptions) string

// Charts
func BarChart(data map[string]int, width int) string
func BarChartWithOptions(data map[string]int, width int, opts *TerminalOptions) string

// Trees
func TreeView(items []TreeItem) string
func TreeViewWithOptions(items []TreeItem, opts *TerminalOptions) string

// Progress
func ProgressBar(current, total int, width int) string
func ProgressBarWithOptions(current, total int, width int, opts *TerminalOptions) string
```

### Styling Functions

```go
// Text styling
func Header(title string, opts *TerminalOptions) string
func Success(text string, opts *TerminalOptions) string
func Warning(text string, opts *TerminalOptions) string
func Error(text string, opts *TerminalOptions) string
func Info(text string, opts *TerminalOptions) string
func Muted(text string, opts *TerminalOptions) string

// Emojis and symbols
func GetEmoji(key string, opts *TerminalOptions) string
func GetSymbol(key string, opts *TerminalOptions) string

// Confidence bars
func CreateConfidenceBar(confidence float64, opts *TerminalOptions) string
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Used By

go-termfmt powers the beautiful terminal output in:

- **[LogSum](https://github.com/yildizm/LogSum)** - High-performance log analysis tool with beautiful terminal UI

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.