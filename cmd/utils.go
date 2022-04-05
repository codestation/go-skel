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
	"encoding/hex"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var typeOfBytes = reflect.TypeOf([]byte(nil))

var unmarshalDecoders = []viper.DecoderConfigOption{
	viper.DecodeHook(HexStringToByteArray()),
	viper.DecodeHook(mapstructure.StringToTimeDurationHookFunc()),
}

func HexStringToByteArray() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t != typeOfBytes {
			return data, nil
		}

		if strings.HasPrefix(data.(string), "s:") {
			return []byte(strings.Split(data.(string), ":")[1]), nil
		}

		return hex.DecodeString(data.(string))
	}
}
