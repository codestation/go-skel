package cfg

import (
	"bytes"
	"slices"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

type testSettings struct {
	HexadecimalValue []byte          `mapstructure:"hexadecimal-value"`
	DecimalValue     decimal.Decimal `mapstructure:"decimal-value"`
	SliceValue       []string        `mapstructure:"slice-value"`
}

func TestUnmarshal(t *testing.T) {
	viper.Set("hexadecimal-value", "deadbeef")
	viper.Set("decimal-value", "123.45")
	viper.Set("duration-value", "2h34m56s")
	viper.Set("slice-value", "a,b,c")

	var settings testSettings
	err := viper.Unmarshal(&settings, unmarshalDecoder)
	if err != nil {
		t.Fatalf("failed to unmarshal settings: %v", err)
	}

	if !bytes.Equal(settings.HexadecimalValue, []byte{0xde, 0xad, 0xbe, 0xef}) {
		t.Errorf("failed to unmarshal hexadecimal value: %v", settings.HexadecimalValue)
	}

	if !settings.DecimalValue.Equal(decimal.NewFromFloat(123.45)) {
		t.Errorf("failed to unmarshal decimal value: %v", settings.DecimalValue)
	}

	if !slices.Equal(settings.SliceValue, []string{"a", "b", "c"}) {
		t.Errorf("failed to unmarshal slice value: %v", settings.SliceValue)
	}
}
