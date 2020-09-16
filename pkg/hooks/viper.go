package hooks

import (
	"encoding/hex"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

var typeOfBytes = reflect.TypeOf([]byte(nil))

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
