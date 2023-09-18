# Loggy

![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/ksdfg/loggy)
[![Test](https://github.com/ksdfg/loggy/actions/workflows/test.yml/badge.svg)](https://github.com/ksdfg/loggy/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/ksdfg/loggy.svg)](https://pkg.go.dev/github.com/ksdfg/loggy)

Loggy provides various handlers for logging messages in different formats and to destinations at the same time. This
allows users to run separate configurations for logging to console, files, sentry, newrelic etc. with just a single
logger object.

## Installation

```shell
go get github.com/ksdfg/loggy
```

## Usage

The `NewCombinedHandler` function creates a logging handler that combines multiple loggers into a single handler. The
package also has other util functions like the `NewConsoleLogHandler` function which return handlers you can use for
logging to various sources.

Here is an example of how you can use it:

```go
package main

import (
	"log/slog"
	"os"

	"github.com/ksdfg/loggy"
)

func main() {
	// Open logfile
	logFile, err := os.Create("app.log")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// Initialize your handlers
	consoleLogHandler := loggy.NewConsoleLogHandler()
	fileLogHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true})

	// Create a combined handler
	combinedLogHandler := loggy.NewCombinedHandler(consoleLogHandler, fileLogHandler)

	// Use handler to create a logger and set it as default
	logger := slog.New(combinedLogHandler)
	slog.SetDefault(logger)

	// Log data
	slog.Debug("this is a debug log")
	slog.Info("this is an info log")
	slog.Warn("this is a warning log")
	slog.Error("this is an error log")
}
```

For more details, you can check the [godoc for this package](https://pkg.go.dev/github.com/ksdfg/loggy).
