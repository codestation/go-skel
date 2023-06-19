// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package jwt

import "github.com/golang-jwt/jwt/v5"

type TokenError struct {
	Token *jwt.Token
	Err   error
}

func (e *TokenError) Error() string {
	return e.Err.Error()
}

func (e *TokenError) Unwrap() error {
	return e.Err
}
