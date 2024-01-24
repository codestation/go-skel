// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package middleware

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapimw "github.com/oapi-codegen/echo-middleware"
)

type ValidatorOption func(*Validator)

func WithSkipperFunc(skipFn middleware.Skipper) ValidatorOption {
	return func(o *Validator) {
		o.skipper = skipFn
	}
}

func MultiError() ValidatorOption {
	return func(o *Validator) {
		o.multiError = true
	}
}

func ErrorHandler(fn oapimw.ErrorHandler) ValidatorOption {
	return func(o *Validator) {
		o.errorHandler = fn
	}
}

func MultiErrorHandler(fn oapimw.MultiErrorHandler) ValidatorOption {
	return func(o *Validator) {
		o.multiErrorHandler = fn
	}
}

func JWTAuth(signingKey any) ValidatorOption {
	return func(o *Validator) {
		o.jwt = echojwt.JWT(signingKey)
	}
}

func JWTAuthWithConfig(config echojwt.Config) ValidatorOption {
	return func(o *Validator) {
		o.jwt = echojwt.WithConfig(config)
	}
}

func BasicAuth(fn middleware.BasicAuthValidator) ValidatorOption {
	return func(o *Validator) {
		o.basic = middleware.BasicAuth(fn)
	}
}

func BasicAuthWithConfig(config middleware.BasicAuthConfig) ValidatorOption {
	return func(o *Validator) {
		o.basic = middleware.BasicAuthWithConfig(config)
	}
}

func KeyAuth(fn middleware.KeyAuthValidator) ValidatorOption {
	return func(o *Validator) {
		o.apiKey = middleware.KeyAuth(fn)
	}
}

func KeyAuthWithConfig(config middleware.KeyAuthConfig) ValidatorOption {
	return func(o *Validator) {
		o.apiKey = middleware.KeyAuthWithConfig(config)
	}
}

func OpenIDConnect(auth *Auth) ValidatorOption {
	return func(o *Validator) {
		o.openIdConnect = NewOIDC(auth)
	}
}

func OapiValidator(spec *openapi3.T, opts ...ValidatorOption) echo.MiddlewareFunc {
	options := Validator{}

	for _, opt := range opts {
		opt(&options)
	}

	oapiOptions := &oapimw.Options{
		ErrorHandler:          options.errorHandler,
		MultiErrorHandler:     options.multiErrorHandler,
		SilenceServersWarning: true,
		Skipper:               options.skipper,
		Options: openapi3filter.Options{
			MultiError:         options.multiError,
			AuthenticationFunc: options.AuthenticatorFunc(),
		},
	}

	return oapimw.OapiRequestValidatorWithOptions(spec, oapiOptions)
}

type Validator struct {
	jwt               echo.MiddlewareFunc
	basic             echo.MiddlewareFunc
	apiKey            echo.MiddlewareFunc
	oauth2            echo.MiddlewareFunc
	openIdConnect     echo.MiddlewareFunc
	skipper           middleware.Skipper
	multiError        bool
	errorHandler      oapimw.ErrorHandler
	multiErrorHandler oapimw.MultiErrorHandler
}

func (v *Validator) next(c echo.Context) error {
	return nil
}

func (v *Validator) AuthenticatorFunc() openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		switch input.SecurityScheme.Type {
		case "http":
			switch input.SecurityScheme.Scheme {
			case "basic":
				if v.basic != nil {
					echoCtx := oapimw.GetEchoContext(ctx)
					return v.basic(v.next)(echoCtx)
				}
			case "bearer":
				switch input.SecurityScheme.BearerFormat {
				case "JWT", "jwt":
					if v.jwt != nil {
						echoCtx := oapimw.GetEchoContext(ctx)
						return v.jwt(v.next)(echoCtx)
					}
				default:
					if v.apiKey != nil {
						echoCtx := oapimw.GetEchoContext(ctx)
						return v.apiKey(v.next)(echoCtx)
					}
				}

			default:
				return fmt.Errorf("unknown http security scheme: %s", input.SecurityScheme.Scheme)
			}
		case "apiKey":
			if v.apiKey != nil {
				echoCtx := oapimw.GetEchoContext(ctx)
				return v.apiKey(v.next)(echoCtx)
			}
		case "oauth2":
			if v.oauth2 != nil {
				echoCtx := oapimw.GetEchoContext(ctx)
				return v.oauth2(v.next)(echoCtx)
			}
		case "openIdConnect":
			if v.openIdConnect != nil {
				echoCtx := oapimw.GetEchoContext(ctx)
				return v.openIdConnect(v.next)(echoCtx)
			}
		default:
			return fmt.Errorf("unknown security scheme type: %s", input.SecurityScheme.Type)
		}

		return openapi3filter.ErrAuthenticationServiceMissing
	}
}
