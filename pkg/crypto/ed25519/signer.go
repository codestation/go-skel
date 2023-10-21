// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package ed25519

import (
	"crypto/ed25519"
	"fmt"
)

type Signer struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

func NewSigner(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) *Signer {
	return &Signer{privateKey: privateKey, publicKey: publicKey}
}

func (s *Signer) Sign(data []byte) ([]byte, error) {
	return Sign(s.privateKey, data), nil
}

func (s *Signer) Verify(data, signature []byte) error {
	if !Verify(s.publicKey, data, signature) {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

func (s *Signer) SignatureSize() int {
	return ed25519.SignatureSize
}
