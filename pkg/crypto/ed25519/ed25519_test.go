// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package ed25519

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"testing"
)

func TestSignVerify(t *testing.T) {
	publicKey, privateKey, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("hello world")
	signature := Sign(privateKey, data)
	if !Verify(publicKey, data, signature) {
		t.Fatal("signature verification failed")
	}
}

func TestCurve255519Conversion(t *testing.T) {
	bobPublicKey, bobPrivateKey, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	alicePublicKey, alicePrivateKey, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	var bobCurve25519PublicKey, bobCurve25519PrivateKey [32]byte
	if err = PublicKeyToCurve25519(&bobCurve25519PublicKey, bobPublicKey); err != nil {
		t.Fatal(err)
	}

	PrivateKeyToCurve25519(&bobCurve25519PrivateKey, bobPrivateKey)

	if bytes.Equal(bobCurve25519PublicKey[:], bobCurve25519PrivateKey[:]) {
		t.Fatal("curve25519 public and private keys are equal")
	}

	bobCurvePrivateKey, err := ecdh.X25519().NewPrivateKey(bobCurve25519PrivateKey[:])
	if err != nil {
		t.Fatal(err)
	}

	bobCurvePublicKey, err := ecdh.X25519().NewPublicKey(bobCurve25519PublicKey[:])
	if err != nil {
		t.Fatal(err)
	}

	var aliceCurve25519PublicKey, aliceCurve25519PrivateKey [32]byte
	if err = PublicKeyToCurve25519(&aliceCurve25519PublicKey, alicePublicKey); err != nil {
		t.Fatal(err)
	}

	PrivateKeyToCurve25519(&aliceCurve25519PrivateKey, alicePrivateKey)

	aliceCurvePrivateKey, err := ecdh.X25519().NewPrivateKey(aliceCurve25519PrivateKey[:])
	if err != nil {
		t.Fatal(err)
	}

	aliceCurvePublicKey, err := ecdh.X25519().NewPublicKey(aliceCurve25519PublicKey[:])
	if err != nil {
		t.Fatal(err)
	}

	bobSharedKey, err := bobCurvePrivateKey.ECDH(aliceCurvePublicKey)
	if err != nil {
		t.Fatal(err)
	}

	aliceSharedKey, err := aliceCurvePrivateKey.ECDH(bobCurvePublicKey)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(bobSharedKey, aliceSharedKey) {
		t.Fatal("shared keys are not equal")
	}
}

func TestGenerateSharedKey(t *testing.T) {
	bobPublicKey, bobPrivateKey, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	alicePublicKey, alicePrivateKey, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	bobSharedKey, err := GenerateSharedKey(bobPublicKey, alicePrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	aliceSharedKey, err := GenerateSharedKey(alicePublicKey, bobPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(bobSharedKey, aliceSharedKey) {
		t.Fatal("shared keys are not equal")
	}
}

func BenchmarkSignVerify(b *testing.B) {
	publicKey, privateKey, err := GenerateKey()
	if err != nil {
		b.Fatal(err)
	}

	plaintext := make([]byte, 32)
	_, err = rand.Read(plaintext)
	if err != nil {
		b.Fatalf("Failed to generate random data: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		signature := Sign(privateKey, plaintext)

		if ok := Verify(publicKey, plaintext, signature); !ok {
			b.Fatalf("Failed to verify signature")
		}
	}
}
