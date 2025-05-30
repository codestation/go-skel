// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package config

import (
	"time"

	"github.com/spf13/pflag"
)

type OIDCSettings struct {
	ProviderTimeout time.Duration `mapstructure:"oidc-provider-timeout"`
	IssuerURL       string        `mapstructure:"oidc-issuer-url"`
	ClientID        string        `mapstructure:"oidc-client-id"`
	ClientSecret    string        `mapstructure:"oidc-client-secret"`
	RedirectURL     string        `mapstructure:"oidc-redirect-url"`
	Scopes          []string      `mapstructure:"oidc-scopes"`
}

func (cfg *OIDCSettings) SetDefaults() {
}

func (cfg *OIDCSettings) Validate() error {
	return nil
}

func LoadOIDCFlags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.Duration("oidc-provider-timeout", 5*time.Second, "OIDC provider timeout")
	fs.String("oidc-issuer-url", "", "OIDC issuer URL")
	fs.String("oidc-client-id", "", "OIDC client ID")
	fs.String("oidc-client-secret", "", "OIDC client secret")
	fs.String("oidc-redirect-url", "", "OIDC redirect URL")
	fs.StringSlice("oidc-scopes", []string{"email", "profile"}, "OIDC scopes")

	return fs
}
