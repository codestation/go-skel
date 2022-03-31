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

package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ApplicationKeySize must be 32 bytes to be used as a key by AES-256
const ApplicationKeySize = 32

// genkeyCmd represents the genkey command
var genkeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "Generate a random key",
	Long:  `Generate a random key to be used as a secret for other configurations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		length := viper.GetInt("length")
		if length < 8 || length > 8192 {
			return errors.New("key length must be between 8 and 8192")
		}

		key := make([]byte, length)
		_, err := rand.Read(key)
		if err != nil {
			return err
		}

		hexKey := hex.EncodeToString(key)

		outputFile := viper.GetString("output")
		if outputFile == "" {
			if viper.GetBool("quiet") {
				fmt.Printf("%s\n", hexKey)
			} else {
				fmt.Printf("Generated key: %s\n", hexKey)
			}
		} else {
			err = ioutil.WriteFile(outputFile, []byte(hexKey+"\n"), 0600)
			if err != nil {
				return err
			}
			if !viper.GetBool("quiet") {
				fmt.Printf("Saved generated key to file %s\n", outputFile)
			}
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(genkeyCmd)
	genkeyCmd.Flags().StringP("output", "o", "", "Save generated key to file")
	genkeyCmd.Flags().BoolP("quiet", "q", false, "Do not print extra messages")
	genkeyCmd.Flags().IntP("length", "l", ApplicationKeySize, "Use an specific key length")
	err := viper.BindPFlags(genkeyCmd.Flags())
	cobra.CheckErr(err)
}
