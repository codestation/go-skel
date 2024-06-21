// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package config

import (
	"errors"

	"github.com/spf13/pflag"
)

const (
	DefaultWorkers   = 5
	DefaultRedisAddr = "127.0.0.1:6379"
)

type GeneralSettings struct {
	Debug         bool   `mapstructure:"debug"`
	LogFormat     string `mapstructure:"log-format"`
	RunMigrations bool   `mapstructure:"run-migrations"`
	EncryptionKey []byte `mapstructure:"encryption-key"`
	RedisAddr     string `mapstructure:"redis-addr"`
	Workers       int    `mapstructure:"workers"`
}

func (cfg *GeneralSettings) Validate() error {
	if len(cfg.EncryptionKey) > 0 && len(cfg.EncryptionKey) < 32 {
		return errors.New("GeneralSettings: encryption key must have at least 32 bytes")
	}

	if cfg.LogFormat != "" && cfg.LogFormat != "json" && cfg.LogFormat != "logfmt" {
		return errors.New("GeneralSettings: log format must be either json or text")
	}

	return nil
}

func (cfg *GeneralSettings) SetDefaults() {
	if cfg.RedisAddr == "" {
		cfg.RedisAddr = DefaultRedisAddr
	}
}

func LoadGeneralFlags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.String("encryption-key", "", "Application encryption key")
	fs.Int("workers", DefaultWorkers, "Workers")
	fs.String("redis-addr", DefaultRedisAddr, "Redis address")

	return fs
}
