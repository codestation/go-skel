// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package ctxlog

import (
	"context"
	"log/slog"
)

type ctxlogKey struct{}

func GetContextAttrs(ctx context.Context) []slog.Attr {
	currentAttrs, _ := ctx.Value(ctxlogKey{}).([]slog.Attr)
	return currentAttrs
}

func AddContextAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	currentAttrs := GetContextAttrs(ctx)
	currentAttrs = append(currentAttrs, attrs...)
	return context.WithValue(ctx, ctxlogKey{}, currentAttrs)
}
