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
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// ApplicationKeySize bust be 32 bytes to be used as a key by AES-256
const ApplicationKeySize = 32

// genkeyCmd represents the genkey command
var genkeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		key := make([]byte, ApplicationKeySize)
		_, err := rand.Read(key)
		if err != nil {
			return err
		}

		hexkey := hex.EncodeToString(key)

		outputFile := viper.GetString("output")
		if outputFile == "" {
			fmt.Printf("Generated key: %s\n", hexkey)
		} else {
			err = ioutil.WriteFile(outputFile, []byte(hexkey+"\n"), 0600)
			if err != nil {
				return err
			}
			fmt.Printf("Saved generated key to file %s\n", outputFile)
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(genkeyCmd)
	genkeyCmd.Flags().StringP("output", "o", "", "Save generated key to file")
	err := viper.BindPFlags(genkeyCmd.Flags())
	if err != nil {
		panic(err)
	}
}
