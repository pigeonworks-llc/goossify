// Package output provides consistent output formatting for CLI commands.
package output

import (
	"encoding/json"
	"io"
	"os"
)

// Format represents the output format type.
type Format string

// FormatHuman and FormatJSON define supported output formats.
const (
	FormatHuman Format = "human"
	FormatJSON  Format = "json"
)

// Formatter handles output formatting for CLI commands.
type Formatter struct {
	format Format
	writer io.Writer
}

// New creates a new Formatter with the specified format.
func New(format string) *Formatter {
	f := Format(format)
	if f != FormatHuman && f != FormatJSON {
		f = FormatHuman
	}
	return &Formatter{format: f, writer: os.Stdout}
}

// SetWriter sets the output writer.
func (f *Formatter) SetWriter(w io.Writer) {
	f.writer = w
}

// IsJSON returns true if the output format is JSON.
func (f *Formatter) IsJSON() bool {
	return f.format == FormatJSON
}

// Output writes data in the configured format.
// For JSON format, it marshals the data.
// For human format, it calls the provided humanFn.
func (f *Formatter) Output(data interface{}, humanFn func()) error {
	switch f.format {
	case FormatJSON:
		return f.outputJSON(data)
	default:
		humanFn()
		return nil
	}
}

// JSON outputs data as formatted JSON.
func (f *Formatter) JSON(data interface{}) error {
	return f.outputJSON(data)
}

func (f *Formatter) outputJSON(data interface{}) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
