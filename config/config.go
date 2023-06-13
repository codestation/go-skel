// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

const (
	DefaultWorkers       = 5
	DefaultRedisAddr     = "127.0.0.1:6379"
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
	DatabaseSettings  DatabaseSettings
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
		return fmt.Errorf("failed to read server settings: %w", err)
	}
	if err = fn(&cfg.DatabaseSettings); err != nil {
		return fmt.Errorf("failed to read database settings: %w", err)
	}
	if err = fn(&cfg.MigrationSettings); err != nil {
		return fmt.Errorf("failed to read migration settings: %w", err)
	}
	return nil
}

func (cfg *Config) SetDefaults() {
	cfg.ServerSettings.SetDefaults()
	cfg.DatabaseSettings.SetDefaults()
	cfg.MigrationSettings.SetDefaults()
}

func (cfg *Config) Validate() error {
	if err := cfg.GeneralSettings.Validate(); err != nil {
		return err
	}
	if err := cfg.ServerSettings.Validate(); err != nil {
		return err
	}
	if err := cfg.DatabaseSettings.Validate(); err != nil {
		return err
	}
	return nil
}

type GeneralSettings struct {
	Debug         bool   `mapstructure:"debug"`
	RunMigrations bool   `mapstructure:"run-migrations"`
	EncryptionKey []byte `mapstructure:"encryption-key"`
	RedisAddr     string `mapstructure:"redis-addr"`
	Workers       int    `mapstructure:"workers"`
}

func (cfg *GeneralSettings) Validate() error {
	if len(cfg.EncryptionKey) > 0 && len(cfg.EncryptionKey) < 32 {
		return errors.New("GeneralSettings: encryption key must have at least 32 bytes")
	}

	return nil
}

type ServerSettings struct {
	ListenAddress    string        `mapstructure:"listen"`
	Timeout          time.Duration `mapstructure:"timeout"`
	ReadTimeout      time.Duration `mapstructure:"read-timeout"`
	WriteTimeout     time.Duration `mapstructure:"write-timeout"`
	IdleTimeout      time.Duration `mapstructure:"idle-timeout"`
	BodyLimit        string        `mapstore:"body-limit"`
	CorsAllowOrigins []string      `mapstructure:"cors-allow-origin"`
	JwtSecret        []byte        `mapstructure:"jwt-secret"`
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

	if len(cfg.CorsAllowOrigins) == 0 {
		cfg.CorsAllowOrigins = append(cfg.CorsAllowOrigins, "*")
	}
}

func (cfg *ServerSettings) Validate() error {
	if len(cfg.JwtSecret) > 0 && len(cfg.JwtSecret) < 32 {
		return errors.New("GeneralSettings: jwt secret must have at least 32 bytes")
	}

	return nil
}

type DatabaseSettings struct {
	DataSourceName  string        `mapstructure:"dsn"`
	MaxIdleConns    int           `mapstructure:"max-idle-conns"`
	MaxOpenConns    int           `mapstructure:"max-open-conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn-max-lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn-max-idle-time"`
	QueryLimit      uint          `mapstructure:"query-limit"`
}

func (cfg *DatabaseSettings) SetDefaults() {
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

func (cfg *DatabaseSettings) Validate() error {
	return nil
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

func LoadDatabaseFlags() *pflag.FlagSet {
	fs := pflag.FlagSet{}
	fs.String("dsn", DefaultDataSourceName, "Database connection string")
	fs.String("driver", DefaultDriverName, "Database driver")
	fs.Int("max-open-conns", DefaultMaxOpenConns, "Max open connections")
	fs.Int("max-idle-conns", DefaultMaxIdleConns, "Max idle connections")
	fs.Duration("conn-max-lifetime", DefaultConnMaxLifetime, "Max lifetime of the connection")
	fs.Duration("conn-max-idle-time", DefaultConnMaxIdleTime, "Max idle time of the connection")
	fs.Int("query-limit", DefaultQueryLimit, "Max results per query")

	return &fs
}

func LoadGeneralFlags() *pflag.FlagSet {
	fs := pflag.FlagSet{}
	fs.String("encryption-key", "", "Application encryption key")
	fs.Int("workers", DefaultWorkers, "Workers")
	fs.String("redis-addr", DefaultRedisAddr, "Redis address")

	return &fs
}

func LoadServerFlags() *pflag.FlagSet {
	fs := pflag.FlagSet{}
	fs.StringP("listen", "l", DefaultListenAddress, "Listen address")
	fs.DurationP("timeout", "t", DefaultReadTimeout, "Request timeout")
	fs.Duration("read-timeout", 0, "Request read timeout")
	fs.Duration("write-timeout", 0, "Request write timeout")
	fs.Duration("idle-timeout", 0, "Request idle timeout")
	fs.String("body-limit", DefaultBodyLimit, "Max body size for http requests")
	fs.StringSlice("cors-allow-origin", []string{}, "CORS Allowed origins")
	fs.String("jwt-secret", "", "JWT secret key")

	return &fs
}

func LoadMigrateFlags() *pflag.FlagSet {
	fs := pflag.FlagSet{}
	fs.Bool("rollback", false, "Rollback last migration")
	fs.Bool("redo", false, "Rollback last migration then migrate again")
	fs.Bool("reset", false, "Drop all tables and run migration")
	fs.Int("step", 1, "Steps to rollback/redo")

	return &fs
}
