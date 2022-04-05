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
	DefaultDataSourceName  = "postgres://goapp:secret@localhost/goapp?sslmode=disable&binary_parameters=yes"
	DefaultMaxIdleConns    = 10
	DefaultMaxOpenConns    = 100
	DefaultConnMaxLifetime = 1 * time.Hour
	DefaultConnMaxIdleTime = 5 * time.Minute
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
	Debug         bool   `mapstructure:"debug"`
	RunMigrations bool   `mapstructure:"run-migrations"`
	EncryptionKey []byte `mapstructure:"encryption-key"`
	JwtSecret     []byte `mapstructure:"jwt-secret"`
}

func (cfg *GeneralSettings) SetDefaults() {}

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
