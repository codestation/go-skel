// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package rsa

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestGenerateSecretKey(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}
	if key.Size() != 256 {
		t.Errorf("RSA key size is incorrect: expected 256, got %d", key.Size())
	}
}

func TestMarshalUnmarshalPrivateKey(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	marshaledKey := MarshalPrivateKey(key)
	unmarshalledKey, err := UnmarshalPrivateKey(marshaledKey)
	if err != nil {
		t.Fatalf("Failed to unmarshal RSA key: %v", err)
	}

	if key.Size() != unmarshalledKey.Size() {
		t.Errorf("RSA key size is incorrect: expected %d, got %d", key.Size(), unmarshalledKey.Size())
	}

	if !key.Equal(unmarshalledKey) {
		t.Errorf("RSA keys are not equal")
	}
}

func TestMarshalUnmarshalPublicKey(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	marshaledKey := MarshalPublicKey(&key.PublicKey)
	unmarshalledKey, err := UnmarshalPublicKey(marshaledKey)
	if err != nil {
		t.Fatalf("Failed to unmarshal RSA key: %v", err)
	}

	if key.Size() != unmarshalledKey.Size() {
		t.Errorf("RSA key size is incorrect: expected %d, got %d", key.Size(), unmarshalledKey.Size())
	}

	if !key.PublicKey.Equal(unmarshalledKey) {
		t.Errorf("RSA keys are not equal")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	publicKey := &key.PublicKey
	plaintext := make([]byte, MaxPlaintextSize)

	_, err = rand.Read(plaintext)
	if err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}

	ciphertext, err := Encrypt(publicKey, plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Errorf("Encryption failed: ciphertext is equal to plaintext")
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decryption failed: decrypted text does not match plaintext")
	}
}

func TestSignVerify(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	publicKey := &key.PublicKey
	plaintext := make([]byte, MaxPlaintextSize)

	_, err = rand.Read(plaintext)
	if err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}

	signature, err := Sign(key, plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	err = Verify(publicKey, plaintext, signature)
	if err != nil {
		t.Fatalf("Failed to verify data: %v", err)
	}
}

func BenchmarkSignVerify(b *testing.B) {
	key, err := GenerateKey()
	if err != nil {
		b.Fatalf("Failed to generate ECDSA key: %v", err)
	}

	plaintext := make([]byte, 32)
	_, err = rand.Read(plaintext)
	if err != nil {
		b.Fatalf("Failed to generate random data: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		signature, err := Sign(key, plaintext)
		if err != nil {
			b.Fatalf("Failed to encrypt data: %v", err)
		}

		err = Verify(&key.PublicKey, plaintext, signature)
		if err != nil {
			b.Fatalf("Failed to verify data: %v", err)
		}
	}
}

func BenchmarkEncryptDecrypt(b *testing.B) {
	key, err := GenerateKey()
	if err != nil {
		b.Fatalf("Failed to generate ECDSA key: %v", err)
	}

	plaintext := make([]byte, 32)
	_, err = rand.Read(plaintext)
	if err != nil {
		b.Fatalf("Failed to generate random data: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ciphertext, err := Encrypt(&key.PublicKey, plaintext)
		if err != nil {
			b.Fatalf("Failed to encrypt data: %v", err)
		}

		_, err = Decrypt(key, ciphertext)
		if err != nil {
			b.Fatalf("Failed to decrypt data: %v", err)
		}
	}
}
