// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package ecdsa

import (
	"crypto/ecdsa"
	"fmt"
)

type Signer struct {
	privateKey *ecdsa.PrivateKey
	Verifier
}

func NewSigner(privateKey *ecdsa.PrivateKey) *Signer {
	return &Signer{privateKey: privateKey}
}

func (s *Signer) Sign(data []byte) ([]byte, error) {
	return Sign(s.privateKey, data)
}

func (s *Signer) Verify(data, signature []byte) error {
	if !Verify(&s.privateKey.PublicKey, data, signature) {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

func (s *Signer) SignatureSize() int {
	return SignatureSize
}

type Verifier struct {
	publicKey *ecdsa.PublicKey
}

func NewVerifier(publicKey *ecdsa.PublicKey) *Verifier {
	return &Verifier{publicKey: publicKey}
}

func (v *Verifier) Verify(data, signature []byte) error {
	if !Verify(v.publicKey, data, signature) {
		return fmt.Errorf("invalid signature")
	}

	return nil
}
