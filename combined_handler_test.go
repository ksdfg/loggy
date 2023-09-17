package loggy_test

import (
	"context"
	"io"
	"log/slog"
	"os"
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

// TestNewCombinedHandler_WithAttrs tests the WithAttr method of the CombinedHandler returned by NewCombinedHandler
// with a string attribute added to the logger that is using the generated CombinedHandler.
func TestNewCombinedHandler_WithAttrs(t *testing.T) {
	// Create a strings.Builder to capture the output
	var outputStream strings.Builder

	// Set up the logger with the combined handler and redirect the output to outputStream
	setupLoggerWithCombinedHandler(&outputStream)

	// Update default logger with attributes
	logger := slog.Default().With(slog.String("test_key", "test_value"))
	slog.SetDefault(logger)

	// Log an error message
	slog.Error("this is an error log")

	// Check the output
	expectedOutput := strings.Join(
		[]string{
			"level=ERROR msg=\"this is an error log\" logger=debugTextLogHandler test_key=test_value",
			"{\"level\":\"ERROR\",\"msg\":\"this is an error log\",\"logger\":\"infoJSONLogHandler\",\"test_key\":\"test_value\"}",
			"level=ERROR msg=\"this is an error log\" logger=warnTextLogHandler test_key=test_value",
			"{\"level\":\"ERROR\",\"msg\":\"this is an error log\",\"logger\":\"errorJSONLogHandler\",\"test_key\":\"test_value\"}\n",
		},
		"\n",
	)
	assert.Equal(t, expectedOutput, outputStream.String())
}

// TestNewCombinedHandler_WithAttrs_EmptyAttrs tests the WithAttr method of the CombinedHandler returned by
// NewCombinedHandler with an empty attribute list added to the logger that is using the generated CombinedHandler.
func TestNewCombinedHandler_WithAttrs_EmptyAttrs(t *testing.T) {
	// Create a strings.Builder to capture the output
	var outputStream strings.Builder

	// Set up the logger with the combined handler and redirect the output to outputStream
	setupLoggerWithCombinedHandler(&outputStream)

	// Update default logger with no attributes
	logger := slog.Default().With()
	slog.SetDefault(logger)

	// Log an error message
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

// TestNewCombinedHandler_WithGroup tests the WithGroup method of the CombinedHandler returned by NewCombinedHandler
// with a test group added to the logger that is using the generated CombinedHandler.
func TestNewCombinedHandler_WithGroup(t *testing.T) {
	// Create a strings.Builder to capture the output
	var outputStream strings.Builder

	// Set up the logger with the combined handler and redirect the output to outputStream
	setupLoggerWithCombinedHandler(&outputStream)

	// Update default logger with new group
	logger := slog.Default().WithGroup("test")
	slog.SetDefault(logger)

	// Log an error message
	slog.Error("this is an error log", slog.String("test_key", "test_value"))

	// Check the output
	expectedOutput := strings.Join(
		[]string{
			"level=ERROR msg=\"this is an error log\" logger=debugTextLogHandler test.test_key=test_value",
			"{\"level\":\"ERROR\",\"msg\":\"this is an error log\",\"logger\":\"infoJSONLogHandler\",\"test\":{\"test_key\":\"test_value\"}}",
			"level=ERROR msg=\"this is an error log\" logger=warnTextLogHandler test.test_key=test_value",
			"{\"level\":\"ERROR\",\"msg\":\"this is an error log\",\"logger\":\"errorJSONLogHandler\",\"test\":{\"test_key\":\"test_value\"}}\n",
		},
		"\n",
	)
	assert.Equal(t, expectedOutput, outputStream.String())
}

// TestNewCombinedHandler_WithGroup_EmptyName tests the WithGroup method of the CombinedHandler returned by NewCombinedHandler
// with an empty group added to the logger that is using the generated CombinedHandler.
func TestNewCombinedHandler_WithGroup_EmptyName(t *testing.T) {
	// Create a strings.Builder to capture the output
	var outputStream strings.Builder

	// Set up the logger with the combined handler and redirect the output to outputStream
	setupLoggerWithCombinedHandler(&outputStream)

	// Update default logger with new group
	logger := slog.Default().WithGroup("")
	slog.SetDefault(logger)

	// Log an error message
	slog.Error("this is an error log", slog.String("test_key", "test_value"))

	// Check the output
	expectedOutput := strings.Join(
		[]string{
			"level=ERROR msg=\"this is an error log\" logger=debugTextLogHandler test_key=test_value",
			"{\"level\":\"ERROR\",\"msg\":\"this is an error log\",\"logger\":\"infoJSONLogHandler\",\"test_key\":\"test_value\"}",
			"level=ERROR msg=\"this is an error log\" logger=warnTextLogHandler test_key=test_value",
			"{\"level\":\"ERROR\",\"msg\":\"this is an error log\",\"logger\":\"errorJSONLogHandler\",\"test_key\":\"test_value\"}\n",
		},
		"\n",
	)
	assert.Equal(t, expectedOutput, outputStream.String())
}

// TestNewCombinedHandler_HandleError tests the CombinedHandler returned by NewCombinedHandler with a closed
// io.Writer passed to one of the handlers passed to NewCombinedHandler.
func TestNewCombinedHandler_HandleError(t *testing.T) {
	// Create a strings.Builder to capture the output
	writer, _, err := os.Pipe()
	if err != nil {
		t.Error(err)
		return
	}

	// Set up the logger with the combined handler and redirect the output to writer generated above
	handler := loggy.NewCombinedHandler(slog.NewTextHandler(writer, nil))
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Close the writer to force an error
	err = writer.Close()
	if err != nil {
		t.Error(err)
		return
	}

	// Log an error message
	slog.Error("this is an error log", slog.String("test_key", "test_value"))
}

// TestNewCombinedHandler_Enabled_False tests the Enabled method of the CombinedHandler returned by NewCombinedHandler
// which has a WARN and ERROR level TextHandlers and checks if INFO level logs are disabled as expected.
func TestNewCombinedHandler_Enabled_False(t *testing.T) {
	// Create new combined handler with loggers for Warn and Error levels
	handler := loggy.NewCombinedHandler(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}),
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}),
	)

	// Set up the logger with the combined handler
	logger := slog.New(handler)

	// Check if Info level is disabled
	assert.Equal(t, false, logger.Enabled(context.Background(), slog.LevelInfo))
}

// TestNewCombinedHandler_Enabled_False tests the Enabled method of the CombinedHandler returned by NewCombinedHandler
// which has an INFO, WARN and ERROR level TextHandlers and checks if INFO level logs are disabled as expected.
func TestNewCombinedHandler_Enabled_True(t *testing.T) {
	// Create new combined handler with loggers for Info, Warn and Error levels
	handler := loggy.NewCombinedHandler(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}),
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}),
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}),
	)

	// Set up the logger with the combined handler
	logger := slog.New(handler)

	// Check if Info level is enabled
	assert.Equal(t, true, logger.Enabled(context.Background(), slog.LevelInfo))
}
