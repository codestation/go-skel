// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cfg

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

type Configurable interface {
	SetDefaults()
	Validate() error
}

var typeOfBytes = reflect.TypeOf([]byte(nil))

var unmarshalDecoder = viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
	HexStringToByteArray(),
	StringToDecimal(),
	StringToTime(),
	mapstructure.StringToTimeDurationHookFunc(),
	mapstructure.StringToSliceHookFunc(","),
))

func ReadConfig(val Configurable) error {
	if err := viper.Unmarshal(val, unmarshalDecoder); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	val.SetDefaults()

	if err := val.Validate(); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}

func HexStringToByteArray() mapstructure.DecodeHookFuncType {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String || t != typeOfBytes {
			return data, nil
		}

		if strings.HasPrefix(data.(string), "s:") {
			return []byte(strings.Split(data.(string), ":")[1]), nil
		}

		return hex.DecodeString(data.(string))
	}
}

func StringToDecimal() mapstructure.DecodeHookFuncType {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String || t != reflect.TypeOf(decimal.Decimal{}) {
			return data, nil
		}

		return decimal.NewFromString(data.(string))
	}
}

func StringToTime() mapstructure.DecodeHookFuncType {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String || t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		return time.Parse(time.RFC3339, data.(string))
	}
}
