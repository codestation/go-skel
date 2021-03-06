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

	"github.com/stretchr/testify/assert"
)

var key = []byte{
	0xfe, 0xd6, 0xff, 0x95, 0x7b, 0x58, 0xfa, 0x7c, 0x41, 0x30, 0x27, 0x0e, 0x07, 0xe4, 0x05, 0x16,
	0x37, 0xca, 0x55, 0x64, 0xa9, 0x7c, 0xae, 0x05, 0x21, 0x4e, 0x3c, 0xec, 0x6c, 0xc2, 0xd7, 0xb5,
}

var nonce = []byte{
	0x0a, 0x65, 0x64, 0x0f, 0xd2, 0x32, 0xd5, 0x21, 0x7b, 0x50, 0x2d, 0x78,
}

var ciphertext = []byte{
	0x69, 0x43, 0x1f, 0x76, 0x39, 0x25, 0x4d, 0x02, 0xbb, 0x69, 0x19, 0x00, 0x09, 0x41, 0x84, 0xf4,
	0x1b, 0x2f, 0x90, 0xbf, 0x5f, 0x27, 0x34, 0xb4, 0x74, 0xc0, 0x29, 0x6f, 0x05,
}

var cleartext = "Hello, World!"

func TestContainer_Seal(t *testing.T) {
	c := Enclave{}
	err := c.SealKey(key, []byte(cleartext))
	assert.NoError(t, err)
	assert.True(t, c.Valid)
}

func TestContainer_Open(t *testing.T) {
	c := Enclave{Nonce: nonce, Data: ciphertext, Valid: true}
	clear, err := c.OpenKey(key)
	assert.NoError(t, err)
	assert.Equal(t, cleartext, string(clear))
}

func TestContainer_Value(t *testing.T) {
	c := Enclave{Nonce: nonce, Data: ciphertext, Valid: true}
	data, err := c.Value()
	assert.NoError(t, err)
	assert.Equal(t, append(nonce, ciphertext...), data)
}

func TestContainer_Scan(t *testing.T) {
	c := Enclave{}
	err := c.Scan(append(nonce, ciphertext...))
	assert.NoError(t, err)
	assert.Equal(t, nonce, c.Nonce)
	assert.Equal(t, ciphertext, c.Data)
	assert.True(t, c.Valid)
}
