// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package crypto

import (
	"bytes"
	"testing"

	"megpoid.dev/go/go-skel/pkg/crypto/aes"
	"megpoid.dev/go/go-skel/pkg/crypto/ecdsa"
	"megpoid.dev/go/go-skel/pkg/crypto/ed25519"
	"megpoid.dev/go/go-skel/pkg/crypto/rsa"
)

func TestCipherSignerEd25519(t *testing.T) {
	public, private, err := ed25519.GenerateKey()
	if err != nil {
		t.Fatalf("cannot generate private key: %v", err)
	}

	secretKey, err := GenerateRandomBytes(aes.KeySize)
	if err != nil {
		t.Fatalf("cannot generate key: %v", err)
	}

	data := []byte("some data")

	signer := ed25519.NewSigner(private, public)
	cipher := aes.NewCipher(secretKey)

	cs := NewCipherSigner(cipher, signer)
	ciphertext, err := cs.EncryptAndSign(data)
	if err != nil {
		t.Fatalf("cannot encrypt and sign data: %v", err)
	}

	plaintext, err := cs.VerifyAndDecrypt(ciphertext)
	if err != nil {
		t.Fatalf("cannot verify and decrypt ciphertext: %v", err)
	}

	if !bytes.Equal(data, plaintext) {
		t.Fatalf("data and plaintext are not equal")
	}

	ciphertext, err = cs.SignAndEncrypt(data)
	if err != nil {
		t.Fatalf("cannot sign and encrypt data: %v", err)
	}

	plaintext, err = cs.DecryptAndVerify(ciphertext)
	if err != nil {
		t.Fatalf("cannot decrypt and verify ciphertext: %v", err)
	}

	if !bytes.Equal(data, plaintext) {
		t.Fatalf("data and plaintext are not equal")
	}
}

func TestCipherSignerECDSA(t *testing.T) {
	private, err := ecdsa.GenerateKey()
	if err != nil {
		t.Fatalf("cannot generate private key: %v", err)
	}

	secretKey, err := GenerateRandomBytes(aes.KeySize)
	if err != nil {
		t.Fatalf("cannot generate key: %v", err)
	}

	data := []byte("some data")

	signer := ecdsa.NewSigner(private)
	cipher := aes.NewCipher(secretKey)

	cs := NewCipherSigner(cipher, signer)
	ciphertext, err := cs.EncryptAndSign(data)
	if err != nil {
		t.Fatalf("cannot encrypt and sign data: %v", err)
	}

	plaintext, err := cs.VerifyAndDecrypt(ciphertext)
	if err != nil {
		t.Fatalf("cannot verify and decrypt ciphertext: %v", err)
	}

	if !bytes.Equal(data, plaintext) {
		t.Fatalf("data and plaintext are not equal")
	}

	ciphertext, err = cs.SignAndEncrypt(data)
	if err != nil {
		t.Fatalf("cannot sign and encrypt data: %v", err)
	}

	plaintext, err = cs.DecryptAndVerify(ciphertext)
	if err != nil {
		t.Fatalf("cannot decrypt and verify ciphertext: %v", err)
	}

	if !bytes.Equal(data, plaintext) {
		t.Fatalf("data and plaintext are not equal")
	}
}

func TestCipherSignerRSA(t *testing.T) {
	private, err := rsa.GenerateKey()
	if err != nil {
		t.Fatalf("cannot generate private key: %v", err)
	}

	secretKey, err := GenerateRandomBytes(aes.KeySize)
	if err != nil {
		t.Fatalf("cannot generate key: %v", err)
	}

	data := []byte("some data")

	signer := rsa.NewSigner(private)
	cipher := aes.NewCipher(secretKey)

	cs := NewCipherSigner(cipher, signer)
	ciphertext, err := cs.EncryptAndSign(data)
	if err != nil {
		t.Fatalf("cannot encrypt and sign data: %v", err)
	}

	plaintext, err := cs.VerifyAndDecrypt(ciphertext)
	if err != nil {
		t.Fatalf("cannot verify and decrypt ciphertext: %v", err)
	}

	if !bytes.Equal(data, plaintext) {
		t.Fatalf("data and plaintext are not equal")
	}

	ciphertext, err = cs.SignAndEncrypt(data)
	if err != nil {
		t.Fatalf("cannot sign and encrypt data: %v", err)
	}

	plaintext, err = cs.DecryptAndVerify(ciphertext)
	if err != nil {
		t.Fatalf("cannot decrypt and verify ciphertext: %v", err)
	}

	if !bytes.Equal(data, plaintext) {
		t.Fatalf("data and plaintext are not equal")
	}
}
