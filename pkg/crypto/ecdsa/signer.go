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
