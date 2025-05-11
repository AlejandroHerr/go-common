package logging

import (
	"context"
	"fmt"
	"log/slog"
)

type ContextHandler struct {
	handler slog.Handler
	keys    map[any]string // Maps context keys to attribute names
}

func NewContextHandler(handler slog.Handler, keys map[any]string) *ContextHandler {
	return &ContextHandler{
		handler: handler,
		keys:    keys,
	}
}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	for ctxKey, attrName := range h.keys {
		if value := ctx.Value(ctxKey); value != nil {
			switch v := value.(type) {
			case string:
				if v != "" {
					r.AddAttrs(slog.String(attrName, v))
				}
			case int:
				r.AddAttrs(slog.Int(attrName, v))
			case int64:
				r.AddAttrs(slog.Int64(attrName, v))
			case float64:
				r.AddAttrs(slog.Float64(attrName, v))
			case bool:
				r.AddAttrs(slog.Bool(attrName, v))
			default:
				// For other types, convert to string
				r.AddAttrs(slog.Any(attrName, v))
			}
		}
	}

	err := h.handler.Handle(ctx, r)
	if err != nil {
		return fmt.Errorf("contextHandler Handle: %w", err)
	}

	return nil
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{
		handler: h.handler.WithAttrs(attrs),
		keys:    h.keys,
	}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{
		handler: h.handler.WithGroup(name),
		keys:    h.keys,
	}
}
