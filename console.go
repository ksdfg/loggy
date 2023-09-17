package loggy

import (
	"io"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

// ConsoleLogWriter represents a logger that writes log messages to stderr.
//
// It implements the `slog.Writer` interface, allowing it to be used as a logger handler.
type ConsoleLogWriter struct {
	outputStream io.Writer
}

// Write writes the log message to the standard error output.
//
// If noColour is false, it colorizes the log message according to the log levels: red for error,
// yellow for warning, and blue for info.
//
// Parameters:
// - p: The byte slice containing the log message.
//
// Returns:
// - n: The number of bytes written.
// - err: An error, if any occurred.
func (w ConsoleLogWriter) Write(p []byte) (n int, err error) {
	log := string(p)

	// Colourise according to log levels
	switch {
	case checkLevel(log, slog.LevelError):
		return color.New(color.FgRed).Fprint(w.outputStream, log)
	case checkLevel(log, slog.LevelWarn):
		return color.New(color.FgYellow).Fprint(w.outputStream, log)
	case checkLevel(log, slog.LevelInfo):
		return color.New(color.FgBlue).Fprint(w.outputStream, log)
	default:
		return w.outputStream.Write([]byte(log))
	}
}

// ConsoleLogWriterOpts represents the options for configuring the behavior of the `ConsoleLogWriter`.
type ConsoleLogWriterOpts struct {
	// JSON specifies whether to use JSON format for log messages.
	JSON bool

	// LogToStdout specifies whether to write log messages to stdout. By default,
	// the ConsoleLogWriter will write logs to stderr.
	LogToStdout bool

	// HandlerOptions contains additional options for the logger handler.
	HandlerOptions slog.HandlerOptions
}

// NewConsoleLogHandler initializes a new stderr log writer based on the given options.
// It returns a slog.Handler that uses the log writer to write log messages to stderr.
func NewConsoleLogHandler(options ...ConsoleLogWriterOpts) slog.Handler {
	// If options are provided, assign the first option to opts
	var opts ConsoleLogWriterOpts
	if len(options) > 0 {
		opts = options[0]
	}

	// Select the stream to output to
	outputStream := os.Stderr
	if opts.LogToStdout {
		outputStream = os.Stdout
	}

	// Create a new ConsoleLogWriter with all the required params
	writer := ConsoleLogWriter{outputStream: outputStream}

	// If the JSON option is enabled, create a new JSON handler using the writer and opts.HandlerOptions
	if opts.JSON {
		return slog.NewJSONHandler(writer, &opts.HandlerOptions)
	}

	// Otherwise, create a new text handler using the writer and opts.HandlerOptions
	return slog.NewTextHandler(writer, &opts.HandlerOptions)
}
