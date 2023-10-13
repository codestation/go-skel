package logger

import (
	"context"
	"log/slog"
)

var _ slog.Handler = Handler{}

type Handler struct {
	handler slog.Handler
	fields  []any
}

// NewContextHandler creates a new Handler with a list of fields to include in each log record from the context.
func NewContextHandler(handler slog.Handler, fields ...any) slog.Handler {
	return Handler{
		handler: handler,
		fields:  fields,
	}
}

func (h Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h Handler) Handle(ctx context.Context, record slog.Record) error {
	for _, field := range h.fields {
		if v := ctx.Value(field); v != nil {
			if key, ok := field.(interface{ String() string }); ok {
				record.AddAttrs(slog.Any(key.String(), v))
			}
		}
	}

	return h.handler.Handle(ctx, record)
}

func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return Handler{handler: h.handler.WithAttrs(attrs), fields: h.fields}
}

func (h Handler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}
