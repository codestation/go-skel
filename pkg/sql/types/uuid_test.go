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
