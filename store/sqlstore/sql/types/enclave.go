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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Enclave struct {
	Nonce []byte
	Data  []byte
	Valid bool
	key   []byte
}

const (
	gcmStandardNonceSize = 12
)

func (c *Enclave) SetKey(key []byte) {
	c.key = key
}

func (c *Enclave) Open() ([]byte, error) {
	return c.OpenKey(c.key)
}

func (c *Enclave) Seal(data []byte) error {
	return c.SealKey(c.key, data)
}

func (c *Enclave) OpenKey(key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	data, err := aesgcm.Open(nil, c.Nonce, c.Data, nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Enclave) SealKey(key, data []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	c.Nonce = make([]byte, gcmStandardNonceSize)
	_, err = rand.Read(c.Nonce)
	if err != nil {
		return err
	}

	c.Valid = true
	c.Data = aesgcm.Seal(nil, c.Nonce, data, nil)
	return nil
}

// Scan implements the Scanner interface.
func (c *Enclave) Scan(value interface{}) error {
	if value == nil {
		*c = Enclave{Valid: false}
		return nil
	}

	switch src := value.(type) {
	case []byte:
		buf := make([]byte, len(src))
		copy(buf, src)
		c.Valid = true
		c.Nonce = buf[:gcmStandardNonceSize]
		c.Data = buf[gcmStandardNonceSize:]
		return nil
	}

	return fmt.Errorf("cannot scan %T", value)
}

func (c Enclave) Value() (driver.Value, error) {
	if c.Valid {
		return append(c.Nonce, c.Data...), nil
	}

	return nil, nil
}

func (c *Enclave) UnmarshalJSON(data []byte) error {
	var plaintext string
	if err := json.Unmarshal(data, &plaintext); err != nil {
		return err
	}

	return c.SealKey(c.key, []byte(plaintext))
}

func (c Enclave) MarshalJSON() ([]byte, error) {
	if !c.Valid {
		return json.Marshal(nil)
	}
	plaintext, err := c.OpenKey(c.key)
	if err != nil {
		return nil, err
	}

	return json.Marshal(plaintext)
}
