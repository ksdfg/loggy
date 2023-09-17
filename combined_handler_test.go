package loggy_test

import (
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ksdfg/loggy"
)

// setupLoggerWithCombinedHandler initializes a logger with a combined handler that outputs logs to the provided
// outputStream with four different loggers
//   - A text logger set to debug level
//   - A JSON logger set to info level
//   - A text logger set to warn level with sources added
//   - A JSON logger set to error level with sources added
//
// It removes the 'time' attribute from the logs to ensure consistent output regardless of when the test is run.
func setupLoggerWithCombinedHandler(outputStream io.Writer) {
	// Replace time from the logs so that the output is exactly the same no matter when the test is run
	removeTimeAttr := func(group []string, attr slog.Attr) slog.Attr {
		if len(group) == 0 && attr.Key == "time" {
			return slog.Attr{}
		}
		return attr
	}

	// Create handlers for different log levels
	debugTextLogHandler := slog.NewTextHandler(
		outputStream, &slog.HandlerOptions{Level: slog.LevelDebug, ReplaceAttr: removeTimeAttr},
	)
	infoJSONLogHandler := slog.NewJSONHandler(
		outputStream, &slog.HandlerOptions{Level: slog.LevelInfo, ReplaceAttr: removeTimeAttr},
	)
	warnTextLogHandler := slog.NewTextHandler(
		outputStream, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: removeTimeAttr},
	)
	errorJSONLogHandler := slog.NewJSONHandler(
		outputStream, &slog.HandlerOptions{Level: slog.LevelError, ReplaceAttr: removeTimeAttr},
	)

	// Add attributes to each handler so that we know which logs came from which handler
	debugTextLogHandlerWithAttr := debugTextLogHandler.WithAttrs(
		[]slog.Attr{slog.String("logger", "debugTextLogHandler")},
	)
	infoJSONLogHandlerWithAttr := infoJSONLogHandler.WithAttrs(
		[]slog.Attr{slog.String("logger", "infoJSONLogHandler")},
	)
	warnTextLogHandlerWithAttr := warnTextLogHandler.WithAttrs(
		[]slog.Attr{slog.String("logger", "warnTextLogHandler")},
	)
	errorJSONLogHandlerWithAttr := errorJSONLogHandler.WithAttrs(
		[]slog.Attr{slog.String("logger", "errorJSONLogHandler")},
	)

	// Create a combined handler with the above handlers
	combinedHandler := loggy.NewCombinedHandler(
		debugTextLogHandlerWithAttr,
		infoJSONLogHandlerWithAttr,
		warnTextLogHandlerWithAttr,
		errorJSONLogHandlerWithAttr,
	)

	// Create a new logger with the combined handler, and set it as the default
	logger := slog.New(combinedHandler)
	slog.SetDefault(logger)
}

// TestNewCombinedHandler_Debug tests the functionality of the NewCombinedHandler function with a debug log.
// This should output a single text log from the logger set using setupLoggerWithCombinedHandler.
func TestNewCombinedHandler_Debug(t *testing.T) {
	// Create a strings.Builder to capture the output
	var outputStream strings.Builder

	// Set up the logger with the combined handler and redirect the output to outputStream
	setupLoggerWithCombinedHandler(&outputStream)

	// Log a debug message
	slog.Debug("this is a debug log")

	// Check the output
	expectedOutput := "level=DEBUG msg=\"this is a debug log\" logger=debugTextLogHandler\n"
	assert.Equal(t, expectedOutput, outputStream.String())
}

// TestNewCombinedHandler_Info tests the functionality of the NewCombinedHandler function with an info log.
// This should output a text log and a JSON log from the logger set using setupLoggerWithCombinedHandler.
func TestNewCombinedHandler_Info(t *testing.T) {
	// Create a strings.Builder to capture the output
	var outputStream strings.Builder

	// Set up the logger with the combined handler and redirect the output to outputStream
	setupLoggerWithCombinedHandler(&outputStream)

	// Log an info message
	slog.Info("this is an info log")

	// Check the output
	expectedOutput := strings.Join(
		[]string{
			"level=INFO msg=\"this is an info log\" logger=debugTextLogHandler",
			"{\"level\":\"INFO\",\"msg\":\"this is an info log\",\"logger\":\"infoJSONLogHandler\"}\n",
		},
		"\n",
	)
	assert.Equal(t, expectedOutput, outputStream.String())
}

// TestNewCombinedHandler_Warn tests the functionality of the NewCombinedHandler function with a warn log.
// This should output two text logs and a JSON log from the logger set using setupLoggerWithCombinedHandler.
func TestNewCombinedHandler_Warn(t *testing.T) {
	// Create a strings.Builder to capture the output
	var outputStream strings.Builder

	// Set up the logger with the combined handler and redirect the output to outputStream
	setupLoggerWithCombinedHandler(&outputStream)

	// Log an info message
	slog.Warn("this is a warning log")

	// Check the output
	expectedOutput := strings.Join(
		[]string{
			"level=WARN msg=\"this is a warning log\" logger=debugTextLogHandler",
			"{\"level\":\"WARN\",\"msg\":\"this is a warning log\",\"logger\":\"infoJSONLogHandler\"}",
			"level=WARN msg=\"this is a warning log\" logger=warnTextLogHandler\n",
		},
		"\n",
	)
	assert.Equal(t, expectedOutput, outputStream.String())
}

// TestNewCombinedHandler_Error tests the functionality of the NewCombinedHandler function with an error log.
// This should output two text logs and two JSON logs from the logger set using setupLoggerWithCombinedHandler.
func TestNewCombinedHandler_Error(t *testing.T) {
	// Create a strings.Builder to capture the output
	var outputStream strings.Builder

	// Set up the logger with the combined handler and redirect the output to outputStream
	setupLoggerWithCombinedHandler(&outputStream)

	// Log an info message
	slog.Error("this is an error log")

	// Check the output
	expectedOutput := strings.Join(
		[]string{
			"level=ERROR msg=\"this is an error log\" logger=debugTextLogHandler",
			"{\"level\":\"ERROR\",\"msg\":\"this is an error log\",\"logger\":\"infoJSONLogHandler\"}",
			"level=ERROR msg=\"this is an error log\" logger=warnTextLogHandler",
			"{\"level\":\"ERROR\",\"msg\":\"this is an error log\",\"logger\":\"errorJSONLogHandler\"}\n",
		},
		"\n",
	)
	assert.Equal(t, expectedOutput, outputStream.String())
}
