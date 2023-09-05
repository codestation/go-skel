// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package apikey

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	oapimw "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	ErrNoAuthHeader = errors.New("token header is missing")
	ErrInvalidToken = errors.New("token is invalid")
)

type ValidatorOption func(opt *oapimw.Options)

func WithSkipperFunc(skipFn middleware.Skipper) ValidatorOption {
	return func(o *oapimw.Options) {
		o.Skipper = skipFn
	}
}

func SharedTokenValidator(spec *openapi3.T, token string, opts ...ValidatorOption) echo.MiddlewareFunc {
	options := &oapimw.Options{
		SilenceServersWarning: true,
		Options: openapi3filter.Options{
			AuthenticationFunc: newAuthenticator(token),
		},
	}

	for _, opt := range opts {
		opt(options)
	}

	return oapimw.OapiRequestValidatorWithOptions(spec, options)
}

func newAuthenticator(token string) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		if input.SecuritySchemeName != "ApiKeyAuth" {
			return fmt.Errorf("security scheme %s != 'ApiKeyAuth'", input.SecuritySchemeName)
		}

		tokenRequest, err := getTokenFromRequest(input.RequestValidationInput.Request)
		if err != nil {
			return fmt.Errorf("failed to get token from request: %w", err)
		}

		if tokenRequest != token {
			return ErrInvalidToken
		}

		return nil
	}
}

func getTokenFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("X-API-KEY")
	// Check for the token header.
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}

	return authHdr, nil
}
