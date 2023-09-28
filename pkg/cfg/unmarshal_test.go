package cfg

import (
	"bytes"
	"slices"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

type testSettings struct {
	HexadecimalValue []byte          `mapstructure:"hexadecimal-value"`
	DecimalValue     decimal.Decimal `mapstructure:"decimal-value"`
	DurationValue    time.Duration   `mapstructure:"duration-value"`
	TimeValue        time.Time       `mapstructure:"time-value"`
	SliceValue       []string        `mapstructure:"slice-value"`
}

func TestUnmarshal(t *testing.T) {
	viper.Set("hexadecimal-value", "deadbeef")
	viper.Set("decimal-value", "123.45")
	viper.Set("duration-value", "2h45m")
	viper.Set("time-value", "2023-09-28T12:34:56Z")
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

	if settings.DurationValue != 2*time.Hour+45*time.Minute {
		t.Errorf("failed to unmarshal duration value: %v", settings.DurationValue)
	}

	if !settings.TimeValue.Equal(time.Date(2023, 9, 28, 12, 34, 56, 0, time.UTC)) {
		t.Errorf("failed to unmarshal time value: %v", settings.TimeValue)
	}

	if !slices.Equal(settings.SliceValue, []string{"a", "b", "c"}) {
		t.Errorf("failed to unmarshal slice value: %v", settings.SliceValue)
	}
}
