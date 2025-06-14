package termfmt

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	// MaxFieldLength is the maximum length for field values before truncation
	MaxFieldLength = 50
)

// terminalFormatter formats data for terminal display
type terminalFormatter struct {
	options      *TerminalOptions
	colorProfile *ColorProfile
}

// NewTerminal creates a new terminal formatter with optional color support
func NewTerminal(color bool) Formatter {
	opts := DefaultOptions()
	opts.Color = color

	return &terminalFormatter{
		options:      opts,
		colorProfile: DefaultColorProfile(),
	}
}

// NewTerminalWithOptions creates a new terminal formatter with custom options
func NewTerminalWithOptions(opts *TerminalOptions) Formatter {
	if opts == nil {
		opts = DefaultOptions()
	}

	return &terminalFormatter{
		options:      opts,
		colorProfile: DefaultColorProfile(),
	}
}

// SetColorProfile sets the color profile for the formatter
func (f *terminalFormatter) SetColorProfile(profile *ColorProfile) {
	f.colorProfile = profile
}

// Format formats the given data for terminal display
func (f *terminalFormatter) Format(data interface{}) ([]byte, error) {
	if data == nil {
		return []byte(""), nil
	}

	var output strings.Builder

	// Handle different data types
	switch v := data.(type) {
	case map[string]interface{}:
		f.formatMap(&output, v, "")
	case []interface{}:
		f.formatSlice(&output, v)
	case string:
		output.WriteString(v)
	default:
		// Use reflection for struct types
		if err := f.formatStruct(&output, data); err != nil {
			// Fallback to JSON representation
			jsonData, jsonErr := json.MarshalIndent(data, "", "  ")
			if jsonErr != nil {
				return nil, fmt.Errorf("failed to format data: %w", jsonErr)
			}

			output.Write(jsonData)
		}
	}

	return []byte(output.String()), nil
}

// formatStruct formats a struct using reflection
func (f *terminalFormatter) formatStruct(output *strings.Builder, data interface{}) error {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			output.WriteString("nil")
			return nil
		}

		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %s", v.Kind())
	}

	// Write struct name as header
	structName := t.Name()
	if structName == "" {
		structName = "Data"
	}

	header := Header(structName, f.options)
	output.WriteString(header + "\n")
	output.WriteString(strings.Repeat("─", len(structName)) + "\n\n")

	// Format fields as tree view
	items := f.structToTreeItems(v, t)
	treeOutput := TreeViewWithOptions(items, f.options)
	output.WriteString(treeOutput)

	return nil
}

// structToTreeItems converts struct fields to tree items
func (f *terminalFormatter) structToTreeItems(v reflect.Value, t reflect.Type) []TreeItem {
	items := make([]TreeItem, 0, v.NumField())

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}

		label := fieldType.Name
		value := f.formatFieldValue(field)

		item := TreeItem{
			Label: label,
			Value: value,
			Last:  i == v.NumField()-1,
		}

		// Handle nested structs
		if field.Kind() == reflect.Struct ||
			(field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct) {
			if field.Kind() == reflect.Ptr && !field.IsNil() {
				item.Children = f.structToTreeItems(field.Elem(), field.Elem().Type())
			} else if field.Kind() == reflect.Struct {
				item.Children = f.structToTreeItems(field, field.Type())
			}
		}

		items = append(items, item)
	}

	return items
}

// formatFieldValue formats a field value for display
func (f *terminalFormatter) formatFieldValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return f.formatStringValue(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', 2, 64)
	case reflect.Bool:
		return f.formatBoolValue(v)
	case reflect.Slice, reflect.Array:
		return fmt.Sprintf("[%d items]", v.Len())
	case reflect.Map:
		return fmt.Sprintf("{%d keys}", v.Len())
	case reflect.Ptr, reflect.Interface:
		return f.formatPointerValue(v)
	case reflect.Struct:
		return fmt.Sprintf("{%s}", v.Type().Name())
	case reflect.Invalid, reflect.Uintptr, reflect.Complex64, reflect.Complex128,
		reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return f.formatSpecialValue(v)
	default:
		return fmt.Sprintf("<%s>", v.Kind())
	}
}

// formatStringValue formats string values with truncation
func (f *terminalFormatter) formatStringValue(v reflect.Value) string {
	str := v.String()
	if len(str) > MaxFieldLength {
		str = str[:MaxFieldLength-3] + "..."
	}

	return fmt.Sprintf("%q", str)
}

// formatBoolValue formats boolean values with styling
func (f *terminalFormatter) formatBoolValue(v reflect.Value) string {
	if v.Bool() {
		return Success("true", f.options)
	}

	return Muted("false", f.options)
}

// formatPointerValue formats pointer and interface values
func (f *terminalFormatter) formatPointerValue(v reflect.Value) string {
	if v.IsNil() {
		return Muted("nil", f.options)
	}

	return f.formatFieldValue(v.Elem())
}

// formatSpecialValue formats special types (complex, chan, func, etc.)
func (f *terminalFormatter) formatSpecialValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return Muted("invalid", f.options)
	case reflect.Uintptr:
		return fmt.Sprintf("0x%x", v.Uint())
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%v", v.Complex())
	case reflect.Chan:
		return fmt.Sprintf("chan<%s>", v.Type().Elem())
	case reflect.Func:
		return fmt.Sprintf("func<%s>", v.Type())
	case reflect.UnsafePointer:
		return fmt.Sprintf("unsafe.Pointer(0x%x)", v.Pointer())
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Array, reflect.Interface,
		reflect.Map, reflect.Ptr, reflect.Slice, reflect.String, reflect.Struct:
		// These should be handled by the main switch, fallback to generic handling
		return fmt.Sprintf("<%s>", v.Kind())
	default:
		return fmt.Sprintf("<%s>", v.Kind())
	}
}

// formatMap formats a map for display
func (f *terminalFormatter) formatMap(
	output *strings.Builder,
	data map[string]interface{},
	prefix string,
) {
	for key, value := range data {
		output.WriteString(prefix)

		keyStr := ColorizeWithProfile(key, "accent", f.colorProfile, f.options)
		output.WriteString(keyStr + ": ")

		switch v := value.(type) {
		case map[string]interface{}:
			output.WriteString("\n")
			f.formatMap(output, v, prefix+"  ")
		case []interface{}:
			fmt.Fprintf(output, "[%d items]\n", len(v))
		case string:
			if len(v) > MaxFieldLength {
				v = v[:MaxFieldLength-3] + "..."
			}

			output.WriteString(fmt.Sprintf("%q", v) + "\n")
		case nil:
			output.WriteString(Muted("nil", f.options) + "\n")
		default:
			fmt.Fprintf(output, "%v\n", v)
		}
	}
}

// formatSlice formats a slice for display
func (f *terminalFormatter) formatSlice(output *strings.Builder, data []interface{}) {
	header := Header("Array Data", f.options)
	output.WriteString(header + "\n")
	output.WriteString("──────────\n\n")

	for i, item := range data {
		fmt.Fprintf(output, "[%d] ", i)

		switch v := item.(type) {
		case map[string]interface{}:
			output.WriteString("{\n")
			f.formatMap(output, v, "  ")
			output.WriteString("}\n")
		case string:
			output.WriteString(fmt.Sprintf("%q", v) + "\n")
		case nil:
			output.WriteString(Muted("nil", f.options) + "\n")
		default:
			fmt.Fprintf(output, "%v\n", v)
		}
	}
}

// FormatAsTable formats data as a table (utility function)
func FormatAsTable(headers []string, rows [][]string, opts *TerminalOptions) string {
	if opts == nil {
		opts = DefaultOptions()
	}

	return TableWithOptions(headers, rows, opts)
}

// FormatAsBox formats content in a box (utility function)
func FormatAsBox(title, content string, opts *TerminalOptions) string {
	if opts == nil {
		opts = DefaultOptions()
	}

	return BoxWithOptions(title, content, opts)
}

// FormatAsBarChart formats data as a bar chart (utility function)
func FormatAsBarChart(data map[string]int, width int, opts *TerminalOptions) string {
	if opts == nil {
		opts = DefaultOptions()
	}

	return BarChartWithOptions(data, width, opts)
}

// Summary creates a formatted summary box
func Summary(title string, items map[string]interface{}, opts *TerminalOptions) string {
	if opts == nil {
		opts = DefaultOptions()
	}

	var content strings.Builder

	maxKeyLen := 0

	// Find max key length for alignment
	for key := range items {
		if len(key) > maxKeyLen {
			maxKeyLen = len(key)
		}
	}

	for key, value := range items {
		padding := maxKeyLen - len(key)
		keyStr := ColorizeWithProfile(key, "info", DefaultColorProfile(), opts)
		valueStr := fmt.Sprintf("%v", value)

		content.WriteString(keyStr + strings.Repeat(" ", padding) + ": " + valueStr + "\n")
	}

	return BoxWithOptions(title, strings.TrimRight(content.String(), "\n"), opts)
}
