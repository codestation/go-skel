// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
)

const (
	// MaxPlaintextSize indicates the maximum plaintext size for RSA 2048 and SHA-256 hash
	MaxPlaintextSize = 190
	SignatureSize    = 256
)

// GenerateKey generates a random RSA private key
func GenerateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("error generating private key: %w", err)
	}

	return privateKey, nil
}

// MarshalPrivateKey marshals a private key to PKCS#1, ASN.1 DER form
func MarshalPrivateKey(privateKey *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(privateKey)
}

func MarshalPublicKey(publicKey *rsa.PublicKey) []byte {
	return x509.MarshalPKCS1PublicKey(publicKey)
}

// UnmarshalPrivateKey unmarshal a private key from PKCS#1, ASN.1 DER form
func UnmarshalPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	privateKey, err := x509.ParsePKCS1PrivateKey(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %w", err)
	}

	return privateKey, nil
}

// UnmarshalPublicKey unmarshal a public key from PKCS#1, ASN.1 DER form
func UnmarshalPublicKey(data []byte) (*rsa.PublicKey, error) {
	publicKey, err := x509.ParsePKCS1PublicKey(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %w", err)
	}

	return publicKey, nil
}

// Sign signs a message using RSA PKCS#1 v1.5
func Sign(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	hashed := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, fmt.Errorf("error signing data: %w", err)
	}

	return signature, nil
}

// Verify verifies a message using RSA PKCS#1 v1.5
func Verify(publicKey *rsa.PublicKey, data, signature []byte) error {
	hashed := sha256.Sum256(data)
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return fmt.Errorf("error verifying data: %w", err)
	}

	return nil
}

// Encrypt encrypts a plaintext using RSA-OAEP
func Encrypt(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	if len(data) > MaxPlaintextSize {
		return nil, fmt.Errorf("data is too large: %d bytes, max %d bytes", len(data), MaxPlaintextSize)
	}

	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, data, nil)
	if err != nil {
		return nil, fmt.Errorf("error encrypting data: %w", err)
	}

	return encryptedData, nil
}

// Decrypt decrypts a ciphertext using RSA-OAEP
func Decrypt(privateKey *rsa.PrivateKey, encryptedData []byte) ([]byte, error) {
	data, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %w", err)
	}

	return data, nil
}
