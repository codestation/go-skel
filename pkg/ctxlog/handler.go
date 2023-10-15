// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package ctxlog

import (
	"context"
	"log/slog"
)

var _ slog.Handler = Handler{}

type Handler struct {
	handler slog.Handler
}

func NewHandler(handler slog.Handler) slog.Handler {
	return Handler{handler: handler}
}

func (h Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h Handler) Handle(ctx context.Context, record slog.Record) error {
	record.AddAttrs(GetContextAttrs(ctx)...)
	return h.handler.Handle(ctx, record)
}

func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return Handler{h.handler.WithAttrs(attrs)}
}

func (h Handler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}
