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
