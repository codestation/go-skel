// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
)

const (
	KeySize              = 32
	gcmBlockSize         = 16
	gcmStandardNonceSize = 12
)

// Encrypt encrypts a plaintext using AES-GCM
func Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating GCM: %w", err)
	}

	nonce := make([]byte, gcmStandardNonceSize)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("error generating nonce: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	return append(nonce, ciphertext...), nil
}

// Decrypt decrypts a ciphertext using AES-GCM
func Decrypt(key, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < gcmStandardNonceSize+gcmBlockSize {
		return nil, errors.New("invalid ciphertext length")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating GCM: %w", err)
	}

	nonce := ciphertext[:gcmStandardNonceSize]
	data := ciphertext[gcmStandardNonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %w", err)
	}

	return plaintext, nil
}
