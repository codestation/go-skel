// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package config

import (
	"errors"
	"fmt"
	"time"
)

const (
	DefaultListenAddress = ":8000"
	DefaultReadTimeout   = 1 * time.Minute
	DefaultWriteTimeout  = 1 * time.Minute
	DefaultIdleTimeout   = 1 * time.Minute
	DefaultBodyLimit     = "10MB"

	DefaultDriverName      = "postgres"
	DefaultDataSourceName  = "postgres://goapp:secret@localhost/goapp?sslmode=disable"
	DefaultMaxIdleConns    = 10
	DefaultMaxOpenConns    = 100
	DefaultConnMaxLifetime = 1 * time.Hour
	DefaultConnMaxIdleTime = 5 * time.Minute
	DefaultQueryLimit      = 1000
)

type Config struct {
	GeneralSettings   GeneralSettings
	ServerSettings    ServerSettings
	SqlSettings       SqlSettings
	MigrationSettings MigrationSettings
}

type Option func(c *Config) error

func WithUnmarshal(fn func(val any) error) Option {
	return func(c *Config) error {
		return c.Unmarshal(fn)
	}
}

func NewConfig(opts ...Option) (*Config, error) {
	cfg := &Config{}
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}
	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) Unmarshal(fn func(val any) error) error {
	var err error
	if err = fn(&cfg.GeneralSettings); err != nil {
		return fmt.Errorf("failed to read general settings: %w", err)
	}
	if err = fn(&cfg.ServerSettings); err != nil {
		return fmt.Errorf("failed to read general settings: %w", err)
	}
	if err = fn(&cfg.SqlSettings); err != nil {
		return fmt.Errorf("failed to read general settings: %w", err)
	}
	if err = fn(&cfg.MigrationSettings); err != nil {
		return fmt.Errorf("failed to read general settings: %w", err)
	}
	return nil
}

func (cfg *Config) SetDefaults() {
	cfg.GeneralSettings.SetDefaults()
	cfg.ServerSettings.SetDefaults()
	cfg.SqlSettings.SetDefaults()
	cfg.MigrationSettings.SetDefaults()
}

func (cfg *Config) Validate() error {
	if err := cfg.GeneralSettings.Validate(); err != nil {
		return err
	}
	return nil
}

type GeneralSettings struct {
	Debug            bool     `mapstructure:"debug"`
	RunMigrations    bool     `mapstructure:"run-migrations"`
	EncryptionKey    []byte   `mapstructure:"encryption-key"`
	JwtSecret        []byte   `mapstructure:"jwt-secret"`
	CorsAllowOrigins []string `mapstructure:"cors-allow-origin"`
}

func (cfg *GeneralSettings) SetDefaults() {
	if len(cfg.CorsAllowOrigins) == 0 {
		cfg.CorsAllowOrigins = append(cfg.CorsAllowOrigins, "*")
	}
}

func (cfg *GeneralSettings) Validate() error {
	if len(cfg.EncryptionKey) > 0 && len(cfg.EncryptionKey) < 32 {
		return errors.New("GeneralSettings: encryption key must have at least 32 bytes")
	}
	if len(cfg.JwtSecret) > 0 && len(cfg.JwtSecret) < 32 {
		return errors.New("GeneralSettings: jwt secret must have at least 32 bytes")
	}
	return nil
}

type ServerSettings struct {
	ListenAddress string        `mapstructure:"listen"`
	Timeout       time.Duration `mapstructure:"timeout"`
	ReadTimeout   time.Duration `mapstructure:"read-timeout"`
	WriteTimeout  time.Duration `mapstructure:"write-timeout"`
	IdleTimeout   time.Duration `mapstructure:"idle-timeout"`
	BodyLimit     string        `mapstore:"body-limit"`
}

func (cfg *ServerSettings) SetDefaults() {
	if cfg.ListenAddress == "" {
		cfg.ListenAddress = DefaultListenAddress
	}
	if cfg.ReadTimeout == 0 {
		if cfg.Timeout != 0 {
			cfg.ReadTimeout = cfg.Timeout
		} else {
			cfg.ReadTimeout = DefaultReadTimeout
		}
	}
	if cfg.WriteTimeout == 0 {
		if cfg.Timeout != 0 {
			cfg.WriteTimeout = cfg.Timeout
		} else {
			cfg.WriteTimeout = DefaultWriteTimeout
		}
	}
	if cfg.IdleTimeout == 0 {
		if cfg.Timeout != 0 {
			cfg.IdleTimeout = cfg.Timeout
		} else {
			cfg.IdleTimeout = DefaultIdleTimeout
		}
	}

	if cfg.BodyLimit == "" {
		cfg.BodyLimit = DefaultBodyLimit
	}
}

type SqlSettings struct {
	DriverName      string        `mapstructure:"driver"`
	DataSourceName  string        `mapstructure:"dsn"`
	MaxIdleConns    int           `mapstructure:"max-idle-conns"`
	MaxOpenConns    int           `mapstructure:"max-open-conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn-max-lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn-max-idle-time"`
	QueryLimit      uint          `mapstructure:"query-limit"`
}

func (cfg *SqlSettings) SetDefaults() {
	if cfg.DriverName == "" {
		cfg.DriverName = DefaultDriverName
	}
	if cfg.DataSourceName == "" {
		cfg.DataSourceName = DefaultDataSourceName
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = DefaultMaxIdleConns
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = DefaultMaxOpenConns
	}
	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = DefaultConnMaxLifetime
	}
	if cfg.ConnMaxIdleTime == 0 {
		cfg.ConnMaxIdleTime = DefaultConnMaxIdleTime
	}
	if cfg.QueryLimit == 0 {
		cfg.QueryLimit = DefaultQueryLimit
	}
}

type MigrationSettings struct {
	Redo     bool
	Rollback bool
	Reset    bool
	Seed     bool
	Step     int
}

func (cfg *MigrationSettings) SetDefaults() {
	if cfg.Step == 0 && (cfg.Rollback || cfg.Redo) {
		cfg.Step = 1
	}
}