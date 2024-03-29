// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package crypto

import (
	"crypto/rand"
	"fmt"
)

type Cipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

type Signer interface {
	Sign(plaintext []byte) ([]byte, error)
	Verifier
}

type Verifier interface {
	Verify(plaintext, signature []byte) error
	SignatureSize() int
}

// GenerateRandomKey generates a random key of a given size
func GenerateRandomKey(size int) ([]byte, error) {
	randomBytes := make([]byte, size)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}

	return randomBytes, nil
}
