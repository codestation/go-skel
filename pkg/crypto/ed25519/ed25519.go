// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package ed25519

import (
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"

	"filippo.io/edwards25519"
)

// GenerateKey generates a public/private key pair for Ed25519
func GenerateKey() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating key pair: %w", err)
	}

	return publicKey, privateKey, nil
}

// MarshalPrivateKey marshals a private key. For Ed25519, this is a no-op
func MarshalPrivateKey(privateKey ed25519.PrivateKey) ([]byte, error) {
	return privateKey, nil
}

// MarshalPublicKey marshals a public key. For Ed25519, this is a no-op
func MarshalPublicKey(publicKey ed25519.PublicKey) ([]byte, error) {
	return publicKey, nil
}

// UnmarshalPrivateKey unmarshal a private key. For Ed25519, this is a no-op
func UnmarshalPrivateKey(data []byte) (ed25519.PrivateKey, error) {
	return data, nil
}

// UnmarshalPublicKey unmarshal a public key. For Ed25519, this is a no-op
func UnmarshalPublicKey(data []byte) (ed25519.PublicKey, error) {
	return data, nil
}

// Sign signs a message using Ed25519
func Sign(privateKey ed25519.PrivateKey, data []byte) []byte {
	hashed := sha256.Sum256(data)
	return ed25519.Sign(privateKey, hashed[:])
}

// Verify verifies a message using Ed25519
func Verify(publicKey ed25519.PublicKey, data, signature []byte) bool {
	hashed := sha256.Sum256(data)
	return ed25519.Verify(publicKey, hashed[:], signature)
}

// PublicKeyToCurve25519 converts an Ed25519 public key into the curve25519
// public key that would be generated from the same private key.
func PublicKeyToCurve25519(ret *[32]byte, publicKey ed25519.PublicKey) error {
	point, err := edwards25519.NewGeneratorPoint().SetBytes(publicKey)
	if err != nil {
		return fmt.Errorf("unable to generate point from publicKey: %w", err)
	}

	copy(ret[:], point.BytesMontgomery())
	return nil
}

// PrivateKeyToCurve25519 converts an ed25519 private key into a corresponding
// curve25519 private key such that the resulting curve25519 public key will
// equal the result from PublicKeyToCurve25519.
func PrivateKeyToCurve25519(ret *[32]byte, privateKey ed25519.PrivateKey) {
	h := sha512.New()
	h.Write(privateKey.Seed())
	copy(ret[:], h.Sum(nil))

	ret[0] &= 248
	ret[31] &= 127
	ret[31] |= 64
}

// GenerateSharedKey generates a shared key using ECDH exchange between a public and private key
func GenerateSharedKey(publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) ([]byte, error) {
	var curve25519PublicKey, curve25519PrivateKey [32]byte
	if err := PublicKeyToCurve25519(&curve25519PublicKey, publicKey); err != nil {
		return nil, fmt.Errorf("failed to convert public key: %w", err)
	}

	PrivateKeyToCurve25519(&curve25519PrivateKey, privateKey)

	curvePrivateKey, err := ecdh.X25519().NewPrivateKey(curve25519PrivateKey[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create private key: %w", err)
	}

	curvePublicKey, err := ecdh.X25519().NewPublicKey(curve25519PublicKey[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create public key: %w", err)
	}

	sharedKey, err := curvePrivateKey.ECDH(curvePublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate shared key: %w", err)
	}

	return sharedKey, nil
}
