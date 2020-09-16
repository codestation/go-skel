package types

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUUID_MarshalJSON(t *testing.T) {
	valid := uuid.Must(uuid.FromString("ec942096-1f47-40f5-a856-78a1d8aff59b"))
	t.Run("Valid", func(t *testing.T) {
		value := UUID{valid}
		result, err := value.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, string(result), "\""+valid.String()+"\"")
	})
}

func TestUUID_UnmarshalJSON(t *testing.T) {
	valid := uuid.Must(uuid.FromString("ec942096-1f47-40f5-a856-78a1d8aff59b"))
	t.Run("Valid", func(t *testing.T) {
		value := UUID{}
		err := value.UnmarshalJSON([]byte("\"" + valid.String() + "\""))
		assert.NoError(t, err)
		assert.Equal(t, valid.String(), value.String())
	})
}
