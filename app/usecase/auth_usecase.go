// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"megpoid.dev/go/go-skel/pkg/apperror"
)

// used to validate that the implementation matches the interface
var _ Auth = &AuthInteractor{}

type JwtCustomClaims struct {
	jwt.RegisteredClaims
	User string `json:"user"`
}

type AuthInteractor struct {
	common
	jwtSecret []byte
}

func NewAuth(jwtSecret []byte) *AuthInteractor {
	return &AuthInteractor{
		common:    newCommon(),
		jwtSecret: jwtSecret,
	}
}

func (uc *AuthInteractor) Login(ctx context.Context, username, password string) (string, error) {
	t := uc.printer(ctx)

	// TODO: really verify username/password
	if len(username) == 0 || len(password) == 0 {
		return "", apperror.NewAuthnError(t.Sprintf("Invalid username or password"), nil)
	}

	claims := JwtCustomClaims{
		User: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// create a jwt token that expires in 24 hours
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(uc.jwtSecret)
	if err != nil {
		return "", apperror.NewAppError(t.Sprintf("Failed to sign token"), err)
	}

	return s, nil
}
