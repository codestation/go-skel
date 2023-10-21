package crypto

import "fmt"

type CipherSigner struct {
	cipher Cipher
	signer Signer
}

func NewCipherSigner(cipher Cipher, signer Signer) *CipherSigner {
	return &CipherSigner{cipher: cipher, signer: signer}
}

func (cs *CipherSigner) EncryptAndSign(plaintext []byte) ([]byte, error) {
	ciphertext, err := cs.cipher.Encrypt(plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt: %w", err)
	}

	signature, err := cs.signer.Sign(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	return append(signature, ciphertext...), nil
}

func (cs *CipherSigner) VerifyAndDecrypt(ciphertext []byte) ([]byte, error) {
	signature := ciphertext[:cs.signer.SignatureSize()]
	ciphertext = ciphertext[cs.signer.SignatureSize():]

	err := cs.signer.Verify(ciphertext, signature)
	if err != nil {
		return nil, fmt.Errorf("failed to verify: %w", err)
	}

	plaintext, err := cs.cipher.Decrypt(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

func (cs *CipherSigner) SignAndEncrypt(plaintext []byte) ([]byte, error) {
	signature, err := cs.signer.Sign(plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	ciphertext, err := cs.cipher.Encrypt(append(signature, plaintext...))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt: %w", err)
	}

	return ciphertext, nil
}

func (cs *CipherSigner) DecryptAndVerify(ciphertext []byte) ([]byte, error) {
	plaintext, err := cs.cipher.Decrypt(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	signature := plaintext[:cs.signer.SignatureSize()]
	plaintext = plaintext[cs.signer.SignatureSize():]

	err = cs.signer.Verify(plaintext, signature)
	if err != nil {
		return nil, fmt.Errorf("failed to verify: %w", err)
	}

	return plaintext, nil
}
