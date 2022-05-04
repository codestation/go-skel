// Copyright (c) 2022 codestation
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package model

import (
	"errors"
	"time"
)

const (
	DefaultListenAddress = ":8000"
	DefaultReadTimeout   = 1 * time.Minute
	DefaultWriteTimeout  = 1 * time.Minute
	DefaultIdleTimeout   = 1 * time.Minute

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
	Step     int
}

func (cfg *MigrationSettings) SetDefaults() {
	if cfg.Step == 0 && (cfg.Rollback || cfg.Redo) {
		cfg.Step = 1
	}
}
