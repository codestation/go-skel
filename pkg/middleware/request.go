// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package middleware

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.megpoid.dev/go-skel/pkg/ctxlog"
)

func SlogRequestID() echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		RequestIDHandler: func(ctx echo.Context, id string) {
			// add request id to context
			requestCtx := ctxlog.AddContextAttrs(ctx.Request().Context(), slog.String("request_id", id))
			ctx.SetRequest(ctx.Request().WithContext(requestCtx))
		},
	})
}
