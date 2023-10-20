// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"fmt"
)

// GenerateKey generates a new ECDSA private key for the P-256 curve
func GenerateKey() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("error generating private key: %w", err)
	}

	return privateKey, nil
}

// MarshalPrivateKey marshals a private key to SEC 1, ASN.1 DER form
func MarshalPrivateKey(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	data, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("error marshaling private key: %w", err)
	}
	return data, nil
}

// MarshalPublicKey marshals a public key to SEC 1, ASN.1 DER form
func MarshalPublicKey(publicKey *ecdsa.PublicKey) ([]byte, error) {
	x509EncodedPub, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("error marshaling public key: %w", err)
	}

	return x509EncodedPub, nil
}

// UnmarshalPrivateKey unmarshal a private key from SEC 1, ASN.1 DER form
func UnmarshalPrivateKey(data []byte) (*ecdsa.PrivateKey, error) {
	return x509.ParseECPrivateKey(data)
}

// UnmarshalPublicKey unmarshal a public key from SEC 1, ASN.1 DER form
func UnmarshalPublicKey(data []byte) (*ecdsa.PublicKey, error) {
	genericPublicKey, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %w", err)
	}

	publicKey, ok := genericPublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("is not a ecdsa public key")
	}

	return publicKey, nil
}

// Sign signs a message using ECDSA and returns an ASN.1 encoded signature
func Sign(privateKey *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	hashed := sha256.Sum256(data)
	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, hashed[:])
	if err != nil {
		return nil, fmt.Errorf("error signing data: %w", err)
	}

	return signature, nil
}

// Verify verifies a message using ECDSA with an ASN.1 encoded signature
func Verify(publicKey *ecdsa.PublicKey, data, signature []byte) bool {
	hashed := sha256.Sum256(data)
	return ecdsa.VerifyASN1(publicKey, hashed[:], signature)
}

// GenerateSharedKey generates a shared key using ECDH exchange between a public and private key
func GenerateSharedKey(publicKey *ecdsa.PublicKey, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	ecdhPrivateKey, err := privateKey.ECDH()
	if err != nil {
		return nil, fmt.Errorf("error generating ECDH private key: %w", err)
	}

	ecdhPublicKey, err := publicKey.ECDH()
	if err != nil {
		return nil, fmt.Errorf("error generating ECDH public key: %w", err)
	}

	shared, err := ecdhPrivateKey.ECDH(ecdhPublicKey)
	if err != nil {
		return nil, fmt.Errorf("error generating shared key: %w", err)
	}

	return shared, nil
}
