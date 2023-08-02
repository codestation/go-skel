// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	oapimw "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const ClaimsContextKey = "jwt_claims"

var (
	ErrNoAuthHeader      = errors.New("authorization header is missing")
	ErrInvalidAuthHeader = errors.New("authorization header is malformed")
)

type ValidatorOption func(opt *oapimw.Options)

func WithSkipperFunc(skipFn middleware.Skipper) ValidatorOption {
	return func(o *oapimw.Options) {
		o.Skipper = skipFn
	}
}

func OapiValidator(spec *openapi3.T, secret []byte, opts ...ValidatorOption) echo.MiddlewareFunc {
	options := &oapimw.Options{
		SilenceServersWarning: true,
		Options: openapi3filter.Options{
			AuthenticationFunc: newAuthenticator(secret),
		},
	}

	for _, opt := range opts {
		opt(options)
	}

	return oapimw.OapiRequestValidatorWithOptions(spec, options)
}

func newAuthenticator(secret []byte) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		if input.SecuritySchemeName != "BearerAuth" {
			return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
		}

		tokenString, err := getJwtFromRequest(input.RequestValidationInput.Request)
		if err != nil {
			return fmt.Errorf("failed to get jwt from request: %w", err)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return secret, nil
		})
		if err != nil {
			return &TokenError{Token: token, Err: err}
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			echoCtx := oapimw.GetEchoContext(ctx)
			echoCtx.Set(ClaimsContextKey, token)
		}

		return nil
	}
}

func getJwtFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	// Check for the Authorization header.
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}
	// We expect a header value of the form "Bearer <token>", with 1 space after
	// Bearer, per spec.
	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}
