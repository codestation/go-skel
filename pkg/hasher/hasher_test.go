// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package hasher

import (
	"bytes"
	"crypto/rand"
	"testing"
)

const infoEncrypt = "encrypt"

func TestHasherPassword(t *testing.T) {
	password := "password"
	hasher, err := NewHasher(password)
	if err != nil {
		t.Fatalf("failed to create hasher: %v", err)
	}

	hash, err := hasher.Hash()
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if !hash.Verify(password) {
		t.Fatal("failed to compare password")
	}

	hashValue := hash.String()
	if err != nil {
		t.Fatalf("failed to get password hash: %v", err)
	}

	hashFromValue, err := NewFromHash(hashValue)
	if err != nil {
		t.Fatalf("failed to get hasher from hash: %v", err)
	}

	if !hashFromValue.Verify(password) {
		t.Fatal("failed to compare password")
	}
}

func TestSecretKeys(t *testing.T) {
	password := "password"
	hasher, err := NewHasher(password)
	if err != nil {
		t.Fatalf("failed to create hasher: %v", err)
	}

	secretKey, err := hasher.DeriveKey(infoEncrypt)
	if err != nil {
		t.Fatalf("failed to get secret key: %v", err)
	}

	anotherHasher, err := NewHasher(password, WithSalt(hasher.Salt))
	if err != nil {
		t.Fatalf("failed to create another hasher: %v", err)
	}

	anotherSecretKey, err := anotherHasher.DeriveKey(infoEncrypt)
	if err != nil {
		t.Fatalf("failed to get another secret key: %v", err)
	}

	if !bytes.Equal(secretKey, anotherSecretKey) {
		t.Fatal("secret keys are not equal")
	}
}

func TestDifferentSecretKeys(t *testing.T) {
	password := "password"
	hasher, err := NewHasher(password)
	if err != nil {
		t.Fatalf("failed to create hasher: %v", err)
	}

	secretKey, err := hasher.DeriveKey(infoEncrypt)
	if err != nil {
		t.Fatalf("failed to get secret key: %v", err)
	}

	anotherHasher, err := NewHasher(password, WithSalt(hasher.Salt))
	if err != nil {
		t.Fatalf("failed to create another hasher: %v", err)
	}

	anotherSecretKey, err := anotherHasher.DeriveKey("another_tag")
	if err != nil {
		t.Fatalf("failed to get another secret key: %v", err)
	}

	if bytes.Equal(secretKey, anotherSecretKey) {
		t.Fatal("secret keys must not be equal")
	}
}

func TestSalts(t *testing.T) {
	password := "password"
	hasher, err := NewHasher(password)
	if err != nil {
		t.Fatalf("failed to create hasher: %v", err)
	}

	anotherHasher, err := NewHasher(password)
	if err != nil {
		t.Fatalf("failed to create another hasher: %v", err)
	}

	if bytes.Equal(hasher.Salt, anotherHasher.Salt) {
		t.Fatal("salts must not be equal")
	}
}

func TestParameters(t *testing.T) {
	password := "password"
	hasher, err := NewHasher(password, WithParameters(1, 16384, 2))
	if err != nil {
		t.Fatalf("failed to create hasher: %v", err)
	}

	hash, err := hasher.Hash()
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	hashString := hash.String()
	if err != nil {
		t.Fatalf("failed to get password hash: %v", err)
	}

	newHash, err := NewFromHash(hashString)
	if err != nil {
		t.Fatalf("failed to get hasher from hash: %v", err)
	}

	if newHash.Threads != 2 {
		t.Fatal("threads must be 2")
	}

	if newHash.TimeCost != 1 {
		t.Fatal("time cost must be 1")
	}

	if newHash.MemoryCost != 16384 {
		t.Fatal("memory cost must be 16384")
	}
}

func BenchmarkHasher(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewHasher("password")
		if err != nil {
			b.Fatalf("failed to create hasher: %v", err)
		}
	}
}

func BenchmarkDeriveKeyFromHash(b *testing.B) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		b.Fatalf("Failed to generate random data: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DeriveKeyFromHash(key, "tag")
		if err != nil {
			b.Fatalf("failed to hash key: %v", err)
		}
	}
}
