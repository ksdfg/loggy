package loggy

import (
	"context"
	"log/slog"
)

// CombinedHandler is a handler that delegates to multiple other handlers.
type CombinedHandler struct {
	handlers []slog.Handler
}

// Enabled reports whether the CombinedHandler handles records at the given level.
//
// The CombinedHandler ignores records whose level is lower.
// This method is called early, before any arguments are processed,
// to save effort if the log event should be discarded.
//
// If called from a Logger method, the first argument is the context
// passed to that method, or context.Background() if nil was passed
// or the method does not take a context.
// The context is passed so Enabled can use its values
// to make a decision.
//
// Parameters:
//   - ctx: The context passed to the method, or context.Background() if nil.
//   - level: The log level to check.
//
// Returns:
//   - enabled: A boolean indicating whether the CombinedHandler is enabled for the given level.
func (h CombinedHandler) Enabled(ctx context.Context, level slog.Level) (enabled bool) {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}

	return false
}

// Handle handles the Record.
// It will only be called when Enabled returns true.
// The Context argument is as for Enabled.
// It is present solely to provide Handlers access to the context's values.
// Canceling the context should not affect record processing.
// (Among other things, log messages may be necessary to debug a
// cancellation-related problem.)
func (h CombinedHandler) Handle(ctx context.Context, record slog.Record) error {
	// Iterate over each handler
	for _, handler := range h.handlers {
		// Check if the handler is enabled for the given context and record level
		if !handler.Enabled(ctx, record.Level) {
			continue
		}

		// Call the handler's Handle function
		err := handler.Handle(ctx, record)
		if err != nil {
			return err
		}
	}

	return nil
}

// WithAttrs returns a new CombinedHandler whose child handlers' attributes consist of
// both the child handlers' attributes and the arguments.
// The CombinedHandler owns the slice: it may retain, modify or discard it.
func (h CombinedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// If no attributes are passed, return the receiver
	if len(attrs) == 0 {
		return h
	}

	// Create a new CombinedHandler
	newHandler := CombinedHandler{}

	// Iterate over each handler in the receiver's handlers slice
	for _, handler := range h.handlers {
		// Call the WithAttrs method on each handler with the given attributes
		// and append the result to the new CombinedHandler's handlers slice
		newHandler.handlers = append(newHandler.handlers, handler.WithAttrs(attrs))
	}

	// Return the new CombinedHandler
	return newHandler
}

// WithGroup returns a new CombinedHandler with the given group appended to
// the child handlers' existing groups.
// The keys of all subsequent attributes, whether added by With or in a
// Record, should be qualified by the sequence of group names.
//
// If the name is empty, WithGroup returns the receiver.
func (h CombinedHandler) WithGroup(name string) slog.Handler {
	// If the name is empty, return the receiver
	if name == "" {
		return h
	}

	// Create a new CombinedHandler
	newHandler := CombinedHandler{}

	// Iterate over all the child handlers
	for _, handler := range h.handlers {
		// Append a new handler with the given group name to the newHandler handlers slice
		newHandler.handlers = append(newHandler.handlers, handler.WithGroup(name))
	}

	// Return the new CombinedHandler
	return newHandler
}

// NewCombinedHandler will return a single CombinedHandler that writes logs to multiple streams via all the handlers
// passed in the arguments.
func NewCombinedHandler(handlers ...slog.Handler) slog.Handler {
	return CombinedHandler{handlers: handlers}
}
