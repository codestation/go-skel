// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package ecdsa

import (
	"bytes"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

func TestGenerateSecretKey(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}
	if key.Curve != elliptic.P256() {
		t.Errorf("ecdsa key curve is incorrect: expected %v, got %v", elliptic.P256(), key.Curve)
	}
}

func TestMarshalUnmarshalPrivateKey(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate ecdsa key: %v", err)
	}

	marshaledKey, err := MarshalPrivateKey(key)
	if err != nil {
		t.Fatalf("Failed to marshal ecdsa key: %v", err)
	}

	unmarshalledKey, err := UnmarshalPrivateKey(marshaledKey)
	if err != nil {
		t.Fatalf("Failed to unmarshal ecdsa key: %v", err)
	}

	if key.Curve != unmarshalledKey.Curve {
		t.Errorf("ecdsa key curve is incorrect: expected %v, got %v", key.Curve, unmarshalledKey.Curve)
	}

	if !key.Equal(unmarshalledKey) {
		t.Errorf("ecdsa keys are not equal")
	}
}

func TestSignVerify(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key: %v", err)
	}

	publicKey := &key.PublicKey
	plaintext := make([]byte, 32)

	_, err = rand.Read(plaintext)
	if err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}

	signature, err := Sign(key, plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	if !Verify(publicKey, plaintext, signature) {
		t.Fatalf("Failed to verify data: %v", err)
	}
}

func TestGenerateSharedKey(t *testing.T) {
	bob, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key: %v", err)
	}

	alice, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate ECDSA key: %v", err)
	}

	bobSharedKey, err := GenerateSharedKey(&alice.PublicKey, bob)
	if err != nil {
		t.Fatalf("Failed to generate shared key: %v", err)
	}

	aliceShredKey, err := GenerateSharedKey(&bob.PublicKey, alice)
	if err != nil {
		t.Fatalf("Failed to generate shared key: %v", err)
	}

	if !bytes.Equal(bobSharedKey, aliceShredKey) {
		t.Errorf("Shared keys are not equal")
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

		if !Verify(&key.PublicKey, plaintext, signature) {
			b.Fatalf("Failed to verify data: %v", err)
		}
	}
}
