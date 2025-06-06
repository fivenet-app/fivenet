package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"golang.org/x/crypto/argon2"
)

const (
	saltLength  = 16 // 128-bit salt
	nonceLength = 12 // for AES-GCM
	keyLength   = 32 // for AES-256
)

type Crypt struct {
	key []byte
}

func New(cfg *config.Config) *Crypt {
	return &Crypt{
		key: []byte(cfg.Secret),
	}
}

func deriveKey(password, salt []byte) []byte {
	// Use Argon2id to derive a key from password+salt
	return argon2.IDKey(password, salt, 2, 64*1024, 4, keyLength)
}

func (c *Crypt) Encrypt(input string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := deriveKey(c.key, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, nonceLength)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ciphertext := aesgcm.Seal(nil, nonce, []byte(input), nil)

	// Store as salt + nonce + ciphertext, all base64-encoded for DB storage
	out := append(salt, nonce...)
	out = append(out, ciphertext...)
	return base64.StdEncoding.EncodeToString(out), nil
}

func (c *Crypt) Decrypt(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	if len(data) < saltLength+nonceLength {
		return "", errors.New("ciphertext too short")
	}
	salt := data[:saltLength]
	nonce := data[saltLength : saltLength+nonceLength]
	ciphertext := data[saltLength+nonceLength:]

	key := deriveKey(c.key, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (c *Crypt) DecryptPointerString(input *string) (*string, error) {
	if input == nil || *input == "" {
		return nil, nil
	}

	out, err := c.Decrypt(*input)
	return &out, err
}

func (c *Crypt) EncryptPointerString(input *string) (*string, error) {
	if input == nil || *input == "" {
		return nil, nil
	}

	out, err := c.Encrypt(*input)
	return &out, err
}
