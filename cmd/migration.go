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

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
	"time"
)

var matchFirstCap = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
var matchAllCap = regexp.MustCompile(`([a-z\d])([A-Z])`)

var templateContent = `-- +migrate Up

-- +migrate StatementBegin
-- Use this block to create functions
-- +migrate StatementEnd

-- Migrate Down
`

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// migrationCmd represents the migrate command
var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "Create a database migration file",
	Long:  `Create a database migration file`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		timestamp := time.Now().Format("20060102150405")
		name := toSnakeCase(args[0])
		filename := path.Join("db/migrations", timestamp+"_"+name+".sql")
		err := ioutil.WriteFile(filename, []byte(templateContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to create migration file: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrationCmd)
	err := viper.BindPFlags(migrationCmd.Flags())
	cobra.CheckErr(err)
}
