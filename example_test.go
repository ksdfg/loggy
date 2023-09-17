package loggy_test

import (
	"log/slog"
	"os"

	"github.com/ksdfg/loggy"
)

func ExampleNewCombinedHandler() {
	// Create handlers for different log levels
	jsonLogHandler := slog.NewJSONHandler(os.Stderr, nil)
	debugLogHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	errorLogHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError, AddSource: true})

	// Add attributes to each handler so that we know which logs came from which handler
	jsonLogHandlerWithAttr := jsonLogHandler.WithAttrs([]slog.Attr{slog.String("logger", "jsonLogHandler")})
	debugLogHandlerWithAttr := debugLogHandler.WithAttrs([]slog.Attr{slog.String("logger", "debugLogHandler")})
	errorLogHandlerWithAttr := errorLogHandler.WithAttrs([]slog.Attr{slog.String("logger", "errorLogHandler")})

	// Create a combined handler with the above handlers
	combinedHandler := loggy.NewCombinedHandler(
		jsonLogHandlerWithAttr, debugLogHandlerWithAttr, errorLogHandlerWithAttr,
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
	// time=2023-09-17T17:44:01.603+05:30 level=DEBUG msg="this is a debug log" logger=debugLogHandler
	// {"time":"2023-09-17T17:44:01.603497782+05:30","level":"INFO","msg":"this is an info log","logger":"jsonLogHandler"}
	// time=2023-09-17T17:44:01.603+05:30 level=INFO msg="this is an info log" logger=debugLogHandler
	// {"time":"2023-09-17T17:44:01.603500753+05:30","level":"WARN","msg":"this is a warning log","logger":"jsonLogHandler"}
	// time=2023-09-17T17:44:01.603+05:30 level=WARN msg="this is a warning log" logger=debugLogHandler
	// {"time":"2023-09-17T17:44:01.603502747+05:30","level":"ERROR","msg":"this is an error log","logger":"jsonLogHandler"}
	// time=2023-09-17T17:44:01.603+05:30 level=ERROR msg="this is an error log" logger=debugLogHandler
	// time=2023-09-17T17:44:01.603+05:30 level=ERROR source=/home/ksdfg/code/loggy/example_test.go:37 msg="this is an error log" logger=errorLogHandler
}

func ExampleNewStderrLogHandler() {
	// Create a new slog.Handler that writes JSON text logs with source info to stderr
	opts := loggy.StderrLogWriterOpts{
		JSON: true, HandlerOptions: slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug},
	}
	handler := loggy.NewStderrLogHandler(opts)

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
