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

// saltLength is the length of the random salt in bytes (128 bits).
// nonceLength is the length of the nonce for AES-GCM.
// keyLength is the length of the derived key for AES-256.
const (
	saltLength  = 16 // 128-bit salt
	nonceLength = 12 // for AES-GCM
	keyLength   = 32 // for AES-256
)

// Crypt provides encryption and decryption using AES-GCM with Argon2id key derivation.
type Crypt struct {
	// key is the base secret used for key derivation
	key []byte
}

// New creates a new Crypt instance using the application's secret from config.
func New(cfg *config.Config) *Crypt {
	return &Crypt{
		key: []byte(cfg.Secret),
	}
}

// deriveKey uses Argon2id to derive a key from the given password and salt.
func deriveKey(password, salt []byte) []byte {
	// Use Argon2id to derive a key from password+salt
	return argon2.IDKey(password, salt, 2, 64*1024, 4, keyLength)
}

// Encrypt encrypts the input string using AES-GCM with a random salt and nonce.
// The output is base64-encoded and includes salt + nonce + ciphertext.
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
	out := make([]byte, 0, saltLength+nonceLength+len(ciphertext))
	out = append(out, salt...)
	out = append(out, nonce...)
	out = append(out, ciphertext...)
	return base64.StdEncoding.EncodeToString(out), nil
}

// Decrypt decrypts a base64-encoded string produced by Encrypt, returning the original plaintext.
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

// DecryptPointerString decrypts a pointer to a string, returning a pointer to the plaintext or nil if input is nil/empty.
func (c *Crypt) DecryptPointerString(input *string) (*string, error) {
	if input == nil || *input == "" {
		return nil, nil
	}

	out, err := c.Decrypt(*input)
	return &out, err
}

// EncryptPointerString encrypts a pointer to a string, returning a pointer to the ciphertext or nil if input is nil/empty.
func (c *Crypt) EncryptPointerString(input *string) (*string, error) {
	if input == nil || *input == "" {
		return nil, nil
	}

	out, err := c.Encrypt(*input)
	return &out, err
}
