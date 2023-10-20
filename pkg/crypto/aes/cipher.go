// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package aes

type Cipher struct {
	key []byte
}

func NewCipher(key []byte) *Cipher {
	return &Cipher{key: key}
}

func (c *Cipher) Encrypt(plaintext []byte) ([]byte, error) {
	return Encrypt(c.key, plaintext)
}

func (c *Cipher) Decrypt(ciphertext []byte) ([]byte, error) {
	return Decrypt(c.key, ciphertext)
}
