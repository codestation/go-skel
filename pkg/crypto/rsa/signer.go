// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package rsa

import "crypto/rsa"

type Signer struct {
	privateKey *rsa.PrivateKey
}

func NewSigner(privateKey *rsa.PrivateKey) *Signer {
	return &Signer{privateKey: privateKey}
}

func (s *Signer) Sign(data []byte) ([]byte, error) {
	return Sign(s.privateKey, data)
}

func (s *Signer) Verify(data, signature []byte) error {
	return Verify(&s.privateKey.PublicKey, data, signature)
}

func (s *Signer) SignatureSize() int {
	return SignatureSize
}

type Cipher struct {
	privateKey *rsa.PrivateKey
}

func NewCipher(privateKey *rsa.PrivateKey) *Cipher {
	return &Cipher{privateKey: privateKey}
}

func (c *Cipher) Encrypt(plaintext []byte) ([]byte, error) {
	return Encrypt(&c.privateKey.PublicKey, plaintext)
}

func (c *Cipher) Decrypt(ciphertext []byte) ([]byte, error) {
	return Decrypt(c.privateKey, ciphertext)
}
