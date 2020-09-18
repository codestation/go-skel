/*
Copyright © 2020 codestation <codestation404@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import "fmt"

// Config holds all the program arguments and environment variables passed to the application.
// The variables set will depend on the command used.
type Config struct {
	Debug      bool
	Addr       string `mapstructure:"listen"`
	DSN        string
	DBHost     string `mapstructure:"db-host"`
	DBPort     string `mapstructure:"db-port"`
	DBUser     string `mapstructure:"db-user"`
	DBPassword string `mapstructure:"db-password"`
	DBName     string `mapstructure:"db-name"`
	DBSSL      bool   `mapstructure:"db-ssl"`
	MasterKey  []byte `mapstructure:"master-key"`
	JWTSecret  []byte `mapstructure:"jwt-secret"`
}

func (c *Config) GetDSN() string {
	if c.DSN != "" {
		return c.DSN
	}

	var sslMode string
	if c.DBSSL {
		sslMode = "enable"
	} else {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		sslMode,
	)

}
