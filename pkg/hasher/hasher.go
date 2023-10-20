// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package hasher

import (
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/hkdf"
	"megpoid.dev/go/go-skel/pkg/crypto"
)

const (
	timeCost      = 1
	memoryCost    = 64 * 1024
	threads       = 2
	secretKeySize = 32
	saltSize      = 16
	InfoAuth      = "auth"
)

func DeriveKeyFromHash(secret []byte, info string) ([]byte, error) {
	reader := hkdf.New(sha256.New, secret, nil, []byte(info))
	secretKey := make([]byte, secretKeySize)

	if _, err := io.ReadFull(reader, secretKey); err != nil {
		return nil, fmt.Errorf("failed to derive secret key: %w", err)
	}

	return secretKey, nil
}

type Option func(*Hasher)

func WithDefaultOpts() Option {
	return func(o *Hasher) {
		o.TimeCost = timeCost
		o.MemoryCost = memoryCost
		o.Threads = threads
	}
}

// WithParameters allow to change the default parameters for the hash.
func WithParameters(timeCost uint32, memoryCost uint32, threads uint8) Option {
	return func(o *Hasher) {
		o.TimeCost = timeCost
		o.MemoryCost = memoryCost
		o.Threads = threads
	}
}

// WithSalt sets the salt to use for the hash.
func WithSalt(salt []byte) Option {
	return func(o *Hasher) {
		o.Salt = salt
	}
}

type Hasher struct {
	TimeCost   uint32
	MemoryCost uint32
	Threads    uint8
	Info       string
	Salt       []byte
	key        []byte
}

// NewHasher returns a new Hasher with the given password and options.
// The hasher can be used to either generate a new hash (for auth) or derive a secret key.
func NewHasher(password string, opts ...Option) (*Hasher, error) {
	s := &Hasher{Info: InfoAuth}
	for _, opt := range opts {
		opt(s)
	}

	if s.TimeCost == 0 || s.Threads == 0 {
		WithDefaultOpts()(s)
	}

	if s.Salt == nil {
		salt, err := crypto.GenerateRandomBytes(saltSize)
		if err != nil {
			return nil, fmt.Errorf("failed to generate salt: %w", err)
		}

		s.Salt = salt
	} else if len(s.Salt) != saltSize {
		return nil, fmt.Errorf("salt must be %d bytes", saltSize)
	}

	s.key = argon2.IDKey([]byte(password), s.Salt, s.TimeCost, s.MemoryCost, s.Threads, secretKeySize)

	return s, nil
}

// Hash generates a new hash that can be used for authentication.
func (h *Hasher) Hash() (*Hash, error) {
	authKey, err := DeriveKeyFromHash(h.key, InfoAuth)
	if err != nil {
		return nil, fmt.Errorf("failed to derive auth key: %w", err)
	}

	hash := &Hash{
		Alg:        hashIdentifier,
		TimeCost:   h.TimeCost,
		MemoryCost: h.MemoryCost,
		Threads:    h.Threads,
		Salt:       h.Salt,
		Hash:       authKey,
		Info:       h.Info,
	}

	return hash, nil
}

// DeriveKey derives a new key from the password hash.
func (h *Hasher) DeriveKey(info string) ([]byte, error) {
	if info == InfoAuth {
		return nil, fmt.Errorf("cannot use %q as info for derived key", InfoAuth)
	}

	key, err := DeriveKeyFromHash(h.key, info)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	return key, nil
}
