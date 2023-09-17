package loggy_test

import (
	"log/slog"
	"os"

	"github.com/ksdfg/loggy"
)

func ExampleNewCombinedHandler() {
	// Create handlers for different log levels
	debugTextLogHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	infoJSONLogHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	warnTextLogHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn, AddSource: true})
	errorJSONLogHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError, AddSource: true})

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

	// Log a debug message
	slog.Debug("this is a debug log")
	// Log an info message
	slog.Info("this is an info log")
	// Log a warning message
	slog.Warn("this is a warning log")
	// Log an error message
	slog.Error("this is an error log")
	// Output:
	// time=2023-09-17T20:01:50.364+05:30 level=DEBUG msg="this is a debug log" logger=debugTextLogHandler
	// time=2023-09-17T20:01:50.364+05:30 level=INFO msg="this is an info log" logger=debugTextLogHandler
	// {"time":"2023-09-17T20:01:50.364653675+05:30","level":"INFO","msg":"this is an info log","logger":"infoJSONLogHandler"}
	// time=2023-09-17T20:01:50.364+05:30 level=WARN msg="this is a warning log" logger=debugTextLogHandler
	// {"time":"2023-09-17T20:01:50.364658189+05:30","level":"WARN","msg":"this is a warning log","logger":"infoJSONLogHandler"}
	// time=2023-09-17T20:01:50.364+05:30 level=WARN source=/home/ksdfg/code/loggy/combined_handler_test.go:56 msg="this is a warning log" logger=warnTextLogHandler
	// time=2023-09-17T20:01:50.364+05:30 level=ERROR msg="this is an error log" logger=debugTextLogHandler
	// {"time":"2023-09-17T20:01:50.36466392+05:30","level":"ERROR","msg":"this is an error log","logger":"infoJSONLogHandler"}
	// time=2023-09-17T20:01:50.364+05:30 level=ERROR source=/home/ksdfg/code/loggy/combined_handler_test.go:58 msg="this is an error log" logger=warnTextLogHandler
	// {"time":"2023-09-17T20:01:50.36466392+05:30","level":"ERROR","source":{"function":"github.com/ksdfg/loggy_test.TestNewCombinedHandler","file":"/home/ksdfg/code/loggy/combined_handler_test.go","line":58},"msg":"this is an error log","logger":"errorJSONLogHandler"}
}

func ExampleNewConsoleLogHandler() {
	// Create a new slog.Handler that writes JSON text logs with source info to stderr
	opts := loggy.ConsoleLogWriterOpts{
		JSON:           true,
		LogToStdout:    true,
		HandlerOptions: slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug},
	}
	handler := loggy.NewConsoleLogHandler(opts)

	// Create a new logger with the above handler, and set it as the default
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Log a debug message
	slog.Debug("this is a debug log")
	// Log an info message
	slog.Info("this is an info log")
	// Log a warning message
	slog.Warn("this is a warning log")
	// Log an error message
	slog.Error("this is an error log")
	// Output:
	// {"time":"2023-09-17T17:44:01.603539649+05:30","level":"DEBUG","source":{"function":"github.com/ksdfg/loggy_test.ExampleNewStderrLogHandler","file":"/home/ksdfg/code/loggy/example_test.go","line":63},"msg":"this is a debug log"}
	// {"time":"2023-09-17T17:44:01.60354488+05:30","level":"INFO","source":{"function":"github.com/ksdfg/loggy_test.ExampleNewStderrLogHandler","file":"/home/ksdfg/code/loggy/example_test.go","line":65},"msg":"this is an info log"}
	// {"time":"2023-09-17T17:44:01.603548573+05:30","level":"WARN","source":{"function":"github.com/ksdfg/loggy_test.ExampleNewStderrLogHandler","file":"/home/ksdfg/code/loggy/example_test.go","line":67},"msg":"this is a warning log"}
	// {"time":"2023-09-17T17:44:01.603551182+05:30","level":"ERROR","source":{"function":"github.com/ksdfg/loggy_test.ExampleNewStderrLogHandler","file":"/home/ksdfg/code/loggy/example_test.go","line":69},"msg":"this is an error log"}
}
