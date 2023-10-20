// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package hasher

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	hashIdentifier = "argon2id-hkdf"
)

// Hash represents a password hash derived from `Argon2id` and HKDF
type Hash struct {
	Alg        string
	TimeCost   uint32
	MemoryCost uint32
	Threads    uint8
	Info       string
	Hash       []byte
	Salt       []byte
}

// NewFromHash created a new hash from an existing password hash
func NewFromHash(hash string) (*Hash, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return nil, fmt.Errorf("invalid hash format")
	}

	identifier := parts[1]
	if identifier != hashIdentifier {
		return nil, fmt.Errorf("invalid hash identifier")
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return nil, fmt.Errorf("invalid hash version format")
	}

	if version != argon2.Version {
		return nil, fmt.Errorf("invalid hash version")
	}

	var (
		hashMemoryCost, hashTimeCost uint32
		hashThreads                  uint8
		hashInfo                     string
	)

	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d,i=%s", &hashMemoryCost, &hashTimeCost, &hashThreads, &hashInfo)
	if err != nil {
		return nil, fmt.Errorf("invalid hash parameters format")
	}

	if hashInfo != InfoAuth {
		return nil, fmt.Errorf("invalid hash info")
	}

	var salt, authKey []byte
	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, fmt.Errorf("invalid hash salt format")
	}
	authKey, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, fmt.Errorf("invalid hash auth key format")
	}

	h := &Hash{
		Alg:        identifier,
		TimeCost:   hashTimeCost,
		MemoryCost: hashMemoryCost,
		Threads:    hashThreads,
		Hash:       authKey,
		Salt:       salt,
		Info:       hashInfo,
	}

	return h, nil
}

// String return an encoded password hash
func (h *Hash) String() string {
	b64Salt := base64.RawStdEncoding.EncodeToString(h.Salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(h.Hash)

	encodedHash := fmt.Sprintf("$%s$v=%d$m=%d,t=%d,p=%d,i=%s$%s$%s", h.Alg, argon2.Version,
		h.MemoryCost, h.TimeCost, h.Threads, h.Info, b64Salt, b64Hash)

	return encodedHash
}

// Verify if a password hash matches the provided password
func (h *Hash) Verify(password string) bool {
	key := argon2.IDKey([]byte(password), h.Salt, h.TimeCost, h.MemoryCost, h.Threads, secretKeySize)
	authKey, err := DeriveKeyFromHash(key, h.Info)
	if err != nil {
		return false
	}

	// always use a constant time function to compare hashes
	return subtle.ConstantTimeCompare(authKey, h.Hash) == 1
}
