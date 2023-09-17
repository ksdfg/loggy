package loggy

import (
	"log/slog"
	"os"

	"github.com/fatih/color"
)

// StderrLogWriter represents a logger that writes log messages to stderr.
//
// It implements the `slog.Writer` interface, allowing it to be used as a logger handler.
type StderrLogWriter struct {
	channel chan []byte
}

// Write writes the log message to the standard error output.
//
// It colorizes the log message according to the log levels: red for error,
// yellow for warning, and blue for info.
//
// Parameters:
// - p: The byte slice containing the log message.
//
// Returns:
// - n: The number of bytes written.
// - err: An error, if any occurred.
func (w StderrLogWriter) Write(p []byte) (n int, err error) {
	log := string(p)

	// Colorize according to log levels
	switch {
	case checkLevel(log, slog.LevelError):
		return color.New(color.FgRed).Fprint(os.Stderr, log)
	case checkLevel(log, slog.LevelWarn):
		return color.New(color.FgYellow).Fprint(os.Stderr, log)
	case checkLevel(log, slog.LevelInfo):
		return color.New(color.FgBlue).Fprint(os.Stderr, log)
	default:
		return os.Stderr.Write([]byte(log))
	}
}

// StderrLogWriterOpts represents the options for configuring the behavior of the `StderrLogWriter`.
type StderrLogWriterOpts struct {
	// JSON specifies whether to use JSON format for log messages.
	JSON bool

	// HandlerOptions contains additional options for the logger handler.
	HandlerOptions slog.HandlerOptions
}

// NewStderrLogHandler initializes a new stderr log writer based on the given options.
// It returns a slog.Handler that uses the log writer to write log messages to stderr.
func NewStderrLogHandler(options ...StderrLogWriterOpts) slog.Handler {
	// If options are provided, assign the first option to opts
	var opts StderrLogWriterOpts
	if len(options) > 0 {
		opts = options[0]
	}

	writer := StderrLogWriter{}

	// If the JSON option is enabled, create a new JSON handler using the writer and opts.HandlerOptions
	if opts.JSON {
		return slog.NewJSONHandler(writer, &opts.HandlerOptions)
	}

	// Otherwise, create a new text handler using the writer and opts.HandlerOptions
	return slog.NewTextHandler(writer, &opts.HandlerOptions)
}
