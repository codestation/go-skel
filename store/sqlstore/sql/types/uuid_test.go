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
