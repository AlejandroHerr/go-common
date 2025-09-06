package logging

import (
	"context"
	"log/slog"
)

var _ slog.Handler = (*DiscardHandler)(nil)

// NewTestLogger returns a no-op logger that discards all output.
// Perfect for unit tests where you don't want to see log output.
func NewTestLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

// DiscardHandler implements slog.Handler but discards all log records.
type DiscardHandler struct{}

// NewDiscardHandler creates a new handler that discards all logs.
func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// Enabled always returns false to minimize overhead.
func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

// Handle discards the record and does nothing.
func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// WithAttrs returns a new handler with the same behavior.
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns a new handler with the same behavior.
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}
