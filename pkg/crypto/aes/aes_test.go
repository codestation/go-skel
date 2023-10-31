// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package aes

import (
	"bytes"
	"crypto/rand"
	"testing"

	"megpoid.dev/go/go-skel/pkg/crypto"
)

func TestEncryptDecrypt(t *testing.T) {
	key, err := crypto.GenerateRandomKey(KeySize)
	if err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}

	plaintext := make([]byte, 1024)
	_, err = rand.Read(plaintext)
	if err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("Encryption failed: ciphertext is equal to plaintext")
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatal("Decryption failed: decrypted text does not match plaintext")
	}
}

func BenchmarkEncryptDecrypt(b *testing.B) {
	key, err := crypto.GenerateRandomKey(KeySize)
	if err != nil {
		b.Fatalf("Failed to generate random key: %v", err)
	}

	plaintext := make([]byte, 1024)
	_, err = rand.Read(plaintext)
	if err != nil {
		b.Fatalf("Failed to generate random data: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ciphertext, err := Encrypt(key, plaintext)
		if err != nil {
			b.Fatalf("Failed to encrypt: %v", err)
		}

		decryptedText, err := Decrypt(key, ciphertext)
		if err != nil {
			b.Fatalf("Failed to decrypt: %v", err)
		}

		if !bytes.Equal(plaintext, decryptedText) {
			b.Fatal("Decryption failed: decrypted text does not match plaintext")
		}
	}
}
