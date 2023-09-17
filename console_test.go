package loggy_test

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"

	"github.com/ksdfg/loggy"
)

// initializeLogger initializes the logger with the given a handler created by NewConsoleLogHandler,
// using options passed in the parameters.
//
// It replaces the "time" attribute from the logs to ensure consistent output regardless of when the test is run.
// The function accepts a single parameter opts of type loggy.ConsoleLogWriterOpts.
// There is no return value.
func initializeLogger(opts loggy.ConsoleLogWriterOpts) {
	// Replace the "time" attribute from the logs so that the output is consistent
	// regardless of when the test is run
	opts.HandlerOptions.ReplaceAttr = func(group []string, attr slog.Attr) slog.Attr {
		if len(group) == 0 && attr.Key == "time" {
			return slog.Attr{}
		}
		return attr
	}

	// Create a new console log handler with the attribute replacement function
	handler := loggy.NewConsoleLogHandler(opts)

	// Create a new logger with the console log handler
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// captureConsoleOutput captures the output from stdout or stderr and returns it as a string.
// If captureStdout is true, it captures stdout. Otherwise, it captures stderr.
// It also calls the logTestCommand function to log the test command.
// Returns the captured output as a string and any error encountered.
func captureConsoleOutput(t *testing.T, captureStdout bool, logTestCommand func()) (string, error) {
	// Create a pipe to capture the output
	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
		return "", err
	}

	// Redirect stdout or stderr to the pipe
	if captureStdout {
		backup := os.Stdout
		os.Stdout = w
		defer func() {
			os.Stdout = backup
		}()
	} else {
		backup := os.Stderr
		os.Stderr = w
		defer func() {
			os.Stderr = backup
		}()
	}

	// Call the logTestCommand function to log the test command
	logTestCommand()

	// Close the write end of the pipe
	err = w.Close()
	if err != nil {
		t.Error(err)
		return "", err
	}

	// Read the captured output from the pipe into a buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Error(err)
		return "", err
	}

	// Return the captured output as a string
	return buf.String(), nil
}

// TestNewConsoleLogHandler_Text_Stdout_Debug tests the NewConsoleLogHandler function with a text debug log to stdout.
func TestNewConsoleLogHandler_Text_Stdout_Debug(t *testing.T) {
	// Capture the console output and any error that occurs
	output, err := captureConsoleOutput(
		t,
		true,
		func() {
			// Set up options for the console log writer
			opts := loggy.ConsoleLogWriterOpts{
				LogToStdout: true,
				HandlerOptions: slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			}
			// Initialize the logger with the options
			initializeLogger(opts)

			// Log a debug message
			slog.Debug("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Check that the output matches the expected value
	assert.Equal(t, "level=DEBUG msg=\"this is a test log\"\n", output)
}

// TestNewConsoleLogHandler_Text_Stdout_Info tests the NewConsoleLogHandler function with a text info log to stdout.
func TestNewConsoleLogHandler_Text_Stdout_Info(t *testing.T) {
	// Capture the console output for testing
	output, err := captureConsoleOutput(
		t,
		true,
		func() {
			// Set up the console log writer options
			opts := loggy.ConsoleLogWriterOpts{LogToStdout: true}
			initializeLogger(opts)

			// Log an info message
			slog.Info("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Define the expected output
	expectedOutput := color.New(color.FgBlue).Sprint("level=INFO msg=\"this is a test log\"\n")

	// Assert that the output matches the expected output
	assert.Equal(t, expectedOutput, output)
}

// TestNewConsoleLogHandler_Text_Stdout_Warn tests the NewConsoleLogHandler function with a text warning log to stdout.
func TestNewConsoleLogHandler_Text_Stdout_Warn(t *testing.T) {
	// Capture the console output and any error
	output, err := captureConsoleOutput(
		t,
		true,
		func() {
			// Set up the options for console log writer
			opts := loggy.ConsoleLogWriterOpts{LogToStdout: true}
			initializeLogger(opts)

			// Log a warning message
			slog.Warn("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Define the expected console output for a warning log
	expectedOutput := color.New(color.FgYellow).Sprint("level=WARN msg=\"this is a test log\"\n")

	// Check if the actual output matches the expected output
	assert.Equal(t, expectedOutput, output)
}

// TestNewConsoleLogHandler_Text_Stdout_Error tests the NewConsoleLogHandler function with a text error log to stdout.
func TestNewConsoleLogHandler_Text_Stdout_Error(t *testing.T) {
	// Capture the console output and any error that occurs
	output, err := captureConsoleOutput(
		t, true, func() {
			// Initialize the logger with console log writer options
			opts := loggy.ConsoleLogWriterOpts{LogToStdout: true}
			initializeLogger(opts)

			// Log an error message
			slog.Error("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Define the expected output with red color formatting
	expectedOutput := color.New(color.FgRed).Sprint("level=ERROR msg=\"this is a test log\"\n")

	// Assert that the output matches the expected output
	assert.Equal(t, expectedOutput, output)
}

// TestNewConsoleLogHandler_JSON_Stdout_Debug tests the NewConsoleLogHandler function with a JSON debug log to stdout.
func TestNewConsoleLogHandler_JSON_Stdout_Debug(t *testing.T) {
	// Capture the console output and any potential errors
	output, err := captureConsoleOutput(
		t, true, func() {
			// Initialize the logger with console log writer options
			opts := loggy.ConsoleLogWriterOpts{
				LogToStdout: true,
				JSON:        true,
				HandlerOptions: slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			}
			initializeLogger(opts)

			// Log a debug message
			slog.Debug("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Assert that the output matches the expected JSON log message
	assert.Equal(t, "{\"level\":\"DEBUG\",\"msg\":\"this is a test log\"}\n", output)
}

// TestNewConsoleLogHandler_JSON_Stdout_Info tests the NewConsoleLogHandler function with a JSON info log to stdout.
func TestNewConsoleLogHandler_JSON_Stdout_Info(t *testing.T) {
	// Capture the console output and any error
	output, err := captureConsoleOutput(
		t, true, func() {
			// Initialize the logger with console log writer options
			opts := loggy.ConsoleLogWriterOpts{LogToStdout: true, JSON: true}
			initializeLogger(opts)

			// Log an info message
			slog.Info("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Define the expected output
	expectedOutput := color.New(color.FgBlue).Sprint("{\"level\":\"INFO\",\"msg\":\"this is a test log\"}\n")

	// Assert that the output matches the expected output
	assert.Equal(t, expectedOutput, output)
}

// TestNewConsoleLogHandler_JSON_Stdout_Warn tests the NewConsoleLogHandler function with a JSON warning log to stdout.
func TestNewConsoleLogHandler_JSON_Stdout_Warn(t *testing.T) {
	// Capture the console output and any errors
	output, err := captureConsoleOutput(
		t, true, func() {
			// Initialize the logger with console log writer options
			opts := loggy.ConsoleLogWriterOpts{LogToStdout: true, JSON: true}
			initializeLogger(opts)

			// Log a warning message
			slog.Warn("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Define the expected output
	expectedOutput := color.New(color.FgYellow).Sprint("{\"level\":\"WARN\",\"msg\":\"this is a test log\"}\n")

	// Assert that the output matches the expected output
	assert.Equal(t, expectedOutput, output)
}

// TestNewConsoleLogHandler_JSON_Stdout tests the NewConsoleLogHandler function with a JSON error log to stdout.
func TestNewConsoleLogHandler_JSON_Stdout_Error(t *testing.T) {
	// Capture the console output and any error that occurs
	output, err := captureConsoleOutput(
		t, true, func() {
			// Initialize the logger with console log writer options
			opts := loggy.ConsoleLogWriterOpts{LogToStdout: true, JSON: true}
			initializeLogger(opts)

			// Log an error message
			slog.Error("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Define the expected output with red color
	expectedOutput := color.New(color.FgRed).Sprint("{\"level\":\"ERROR\",\"msg\":\"this is a test log\"}\n")
	// Assert that the actual output matches the expected output
	assert.Equal(t, expectedOutput, output)
}

// TestNewConsoleLogHandler_Text_Stderr_Debug tests the NewConsoleLogHandler function with a text debug log to stderr.
func TestNewConsoleLogHandler_Text_Stderr_Debug(t *testing.T) {
	// Capture the console output
	output, err := captureConsoleOutput(
		t,
		false,
		func() {
			// Initialize the logger with console log writer options
			opts := loggy.ConsoleLogWriterOpts{HandlerOptions: slog.HandlerOptions{Level: slog.LevelDebug}}
			initializeLogger(opts)

			// Log a debug message
			slog.Debug("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Assert that the output matches the expected log message
	assert.Equal(t, "level=DEBUG msg=\"this is a test log\"\n", output)
}

// TestNewConsoleLogHandler_Text_Stderr_Info tests the NewConsoleLogHandler function with a text info log to stderr.
func TestNewConsoleLogHandler_Text_Stderr_Info(t *testing.T) {
	// Capture the output of the function
	output, err := captureConsoleOutput(
		t, false, func() {
			// Initialize the logger with console log writer options
			opts := loggy.ConsoleLogWriterOpts{}
			initializeLogger(opts)

			// Log an info message
			slog.Info("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Set the expected output of the function
	expectedOutput := color.New(color.FgBlue).Sprint("level=INFO msg=\"this is a test log\"\n")

	// Assert that the output matches the expected output
	assert.Equal(t, expectedOutput, output)
}

// TestNewConsoleLogHandler_Text_Stderr tests the NewConsoleLogHandler function with a text warning log to stderr.
func TestNewConsoleLogHandler_Text_Stderr_Warn(t *testing.T) {
	// Capture the console output and any error
	output, err := captureConsoleOutput(
		t, false, func() {
			// Initialize the logger with console log writer options
			opts := loggy.ConsoleLogWriterOpts{}
			initializeLogger(opts)

			// Log a warning message
			slog.Warn("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Define the expected output
	expectedOutput := color.New(color.FgYellow).Sprint("level=WARN msg=\"this is a test log\"\n")

	// Assert that the output matches the expected output
	assert.Equal(t, expectedOutput, output)
}

// TestNewConsoleLogHandler_Text_Stderr_Error tests the NewConsoleLogHandler function with a text error log to stderr.
func TestNewConsoleLogHandler_Text_Stderr_Error(t *testing.T) {
	// Capture the console output and any error
	output, err := captureConsoleOutput(
		t, false, func() {
			// Initialize the logger with console log writer options
			opts := loggy.ConsoleLogWriterOpts{}
			initializeLogger(opts)

			// Log an error message
			slog.Error("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Define the expected console output
	expectedOutput := color.New(color.FgRed).Sprint("level=ERROR msg=\"this is a test log\"\n")

	// Assert that the actual output matches the expected output
	assert.Equal(t, expectedOutput, output)
}

// TestNewConsoleLogHandler_JSON_Stderr tests the NewConsoleLogHandler function with a JSON debug log to stderr.
func TestNewConsoleLogHandler_JSON_Stderr_Debug(t *testing.T) {
	// Capture the console output and any error
	output, err := captureConsoleOutput(
		t, false, func() {
			// Initialize the logger with console log writer options
			opts := loggy.ConsoleLogWriterOpts{
				JSON:           true,
				HandlerOptions: slog.HandlerOptions{Level: slog.LevelDebug},
			}
			initializeLogger(opts)

			// Log a debug message
			slog.Debug("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Assert that the captured output matches the expected value
	assert.Equal(t, "{\"level\":\"DEBUG\",\"msg\":\"this is a test log\"}\n", output)
}

// TestNewConsoleLogHandler_JSON_Stderr_Info tests the NewConsoleLogHandler function with a JSON info log to stdout.
func TestNewConsoleLogHandler_JSON_Stderr_Info(t *testing.T) {
	output, err := captureConsoleOutput(
		t, false, func() {
			// Initialize the logger with console log writer options
			opts := loggy.ConsoleLogWriterOpts{JSON: true}
			initializeLogger(opts)

			// Log an info message
			slog.Info("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Define the expected console output
	expectedOutput := color.New(color.FgBlue).Sprint("{\"level\":\"INFO\",\"msg\":\"this is a test log\"}\n")

	// Assert that the output matches the expected format
	assert.Equal(t, expectedOutput, output)
}

// TestNewConsoleLogHandler_JSON_Stderr_Warn tests the NewConsoleLogHandler function with a JSON warning log to stdout.
func TestNewConsoleLogHandler_JSON_Stderr_Warn(t *testing.T) {
	// Capture the console output and any error
	output, err := captureConsoleOutput(
		t, false, func() {
			// Initialize the logger with console log writer options
			opts := loggy.ConsoleLogWriterOpts{JSON: true}
			initializeLogger(opts)

			// Log a warning message
			slog.Warn("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Define the expected output in JSON format with a yellow color
	expectedOutput := color.New(color.FgYellow).Sprint("{\"level\":\"WARN\",\"msg\":\"this is a test log\"}\n")

	// Assert that the actual output matches the expected output
	assert.Equal(t, expectedOutput, output)
}

// TestNewConsoleLogHandler_JSON_Stderr_Error tests the NewConsoleLogHandler function with a JSON error log to stdout.
func TestNewConsoleLogHandler_JSON_Stderr_Error(t *testing.T) {
	// Capture the console output and any error
	output, err := captureConsoleOutput(
		t,
		false,
		func() {
			// Initialize the logger with console log writer options
			opts := loggy.ConsoleLogWriterOpts{JSON: true}
			initializeLogger(opts)

			// Log an error message
			slog.Error("this is a test log")
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	// Define the expected output
	expectedOutput := color.New(color.FgRed).Sprint("{\"level\":\"ERROR\",\"msg\":\"this is a test log\"}\n")

	// Assert that the output matches the expected output
	assert.Equal(t, expectedOutput, output)
}
