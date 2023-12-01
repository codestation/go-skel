// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package middleware

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"go.megpoid.dev/go-skel/pkg/crypto"
	"golang.org/x/oauth2"
)

const (
	OauthCookieState         = "oauth_state"
	OauthCookieNonce         = "oauth_nonce"
	OauthCookieAccessToken   = "oauth_access_token"
	OauthCookieRefreshToken  = "oauth_refresh_token"
	OauthRefreshTokenTimeout = 30 * 24 * time.Hour
)

var UserClaims = "oidcClaims"

type Claims map[string]any

type Config struct {
	IssuerURL    string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

type Auth struct {
	Provider    *oidc.Provider
	Verifier    *oidc.IDTokenVerifier
	OAuthConfig *oauth2.Config
}

func NewOIDCAuth(ctx context.Context, config *Config) (*Auth, error) {
	provider, err := oidc.NewProvider(ctx, config.IssuerURL)
	if err != nil {
		return nil, err
	}

	scopes := make([]string, 0)
	scopes = append(scopes, config.Scopes...)

	if !slices.Contains(scopes, "profile") {
		scopes = append(scopes, "profile")
	}
	if !slices.Contains(scopes, "email") {
		scopes = append(scopes, "email")
	}

	oauth2Config := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: append([]string{oidc.ScopeOpenID}, scopes...),
	}

	oidcVerifier := provider.Verifier(&oidc.Config{
		ClientID: config.ClientID,
	})

	auth := &Auth{
		Provider:    provider,
		Verifier:    oidcVerifier,
		OAuthConfig: &oauth2Config,
	}

	return auth, nil
}

func GetClaims(c echo.Context) (Claims, error) {
	claims, ok := c.Get(UserClaims).(Claims)
	if ok {
		return claims, nil
	}

	return Claims{}, errors.New("failed to get user claims")
}

func (auth *Auth) RedirectHandler(e echo.Context) error {
	state, err := crypto.GenerateRandomKey(16)
	if err != nil {
		return fmt.Errorf("failed to generate state cookie: %w", err)
	}

	stateValue := base64.RawURLEncoding.EncodeToString(state)

	nonce, err := crypto.GenerateRandomKey(16)
	if err != nil {
		return fmt.Errorf("failed to generate nonce cookie: %w", err)
	}

	nonceValue := base64.RawURLEncoding.EncodeToString(nonce)

	e.SetCookie(&http.Cookie{
		Name:     OauthCookieState,
		Value:    stateValue,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   e.Request().URL.Scheme == "https",
		HttpOnly: true,
	})

	e.SetCookie(&http.Cookie{
		Name:     OauthCookieNonce,
		Value:    nonceValue,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   e.Request().URL.Scheme == "https",
		HttpOnly: true,
	})

	return e.Redirect(http.StatusFound, auth.OAuthConfig.AuthCodeURL(stateValue, oidc.Nonce(nonceValue)))
}

func (auth *Auth) CallbackHandler(c echo.Context) error {
	state, err := c.Cookie(OauthCookieState)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to get state cookie: %v", err))
	}

	req := c.Request()

	if req.URL.Query().Get("state") != state.Value {
		return echo.NewHTTPError(http.StatusBadRequest, "state did not match")
	}

	oauth2Token, err := auth.OAuthConfig.Exchange(req.Context(), req.URL.Query().Get("code"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to exchange oauth token: %v", err))
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to extract id_token from oauth2 token")
	}

	// Parse and verify ID Token payload.
	idToken, err := auth.Verifier.Verify(req.Context(), rawIDToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to verify ID Token: %v", err))
	}

	nonce, err := c.Cookie(OauthCookieNonce)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to get nonce cookie: %v", err))
	}

	if idToken.Nonce != nonce.Value {
		return echo.NewHTTPError(http.StatusBadRequest, "nonce did not match")
	}

	mapClaims := make(Claims)
	if err := idToken.Claims(&mapClaims); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to read claim from ID Token: %v", err))
	}

	c.Set(UserClaims, mapClaims)

	// Get access token expiry
	accessTokenExpire := time.Now().Add(time.Duration(oauth2Token.Expiry.Unix()-time.Now().Unix()) * time.Second)
	// Get refresh token expiry
	refreshTokenExpire := time.Now().Add(OauthRefreshTokenTimeout)

	c.SetCookie(&http.Cookie{
		Name:     OauthCookieAccessToken,
		Value:    oauth2Token.AccessToken,
		MaxAge:   accessTokenExpire.Second(),
		Secure:   c.Request().URL.Scheme == "https",
		HttpOnly: true,
	})

	c.SetCookie(&http.Cookie{
		Name:     OauthCookieRefreshToken,
		Value:    oauth2Token.RefreshToken,
		MaxAge:   refreshTokenExpire.Second(),
		Secure:   c.Request().URL.Scheme == "https",
		HttpOnly: true,
	})

	return nil
}

func (auth *Auth) RefreshHandler(c echo.Context) error {
	refreshCookie, err := c.Cookie(OauthCookieRefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("failed to get refresh token cookie: %v", err))
	}

	token := &oauth2.Token{
		RefreshToken: refreshCookie.Value,
	}

	newToken, err := auth.OAuthConfig.TokenSource(c.Request().Context(), token).Token()
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("failed to refresh token: %v", err))
	}

	accessTokenExpire := time.Now().Add(time.Duration(newToken.Expiry.Unix()-time.Now().Unix()) * time.Second)
	refreshTokenExpire := time.Now().Add(OauthRefreshTokenTimeout)

	c.SetCookie(&http.Cookie{
		Name:     OauthCookieAccessToken,
		Value:    newToken.AccessToken,
		MaxAge:   accessTokenExpire.Second(),
		Secure:   c.Request().URL.Scheme == "https",
		HttpOnly: true,
	})

	c.SetCookie(&http.Cookie{
		Name:     OauthCookieRefreshToken,
		Value:    newToken.RefreshToken,
		MaxAge:   refreshTokenExpire.Second(),
		Secure:   c.Request().URL.Scheme == "https",
		HttpOnly: true,
	})

	return nil
}

func (auth *Auth) LogoutHandler(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:    OauthCookieAccessToken,
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})

	c.SetCookie(&http.Cookie{
		Name:    OauthCookieRefreshToken,
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})

	c.SetCookie(&http.Cookie{
		Name:    OauthCookieState,
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})

	c.SetCookie(&http.Cookie{
		Name:    OauthCookieNonce,
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})

	return nil
}

func NewOIDC(auth *Auth) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := auth.CallbackHandler(c); err != nil {
				return err
			}
			return next(c)
		}
	}
}
