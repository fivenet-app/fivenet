package crypt

import (
	"testing"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{Secret: "test-secret"}
	crypt := New(cfg)

	plaintext := "Hello, World!"
	encrypted, err := crypt.Encrypt(plaintext)
	require.NoError(t, err, "Encrypt failed")

	assert.NotEmpty(t, encrypted, "Encrypt returned an empty string")

	// Ensure the encrypted string is not the same as the plaintext
	assert.NotEqual(t, plaintext, encrypted, "Encrypted string should not match the plaintext")
}

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{Secret: "test-secret"}
	crypt := New(cfg)

	require := require.New(t)

	plaintext := "Hello, World!"
	encrypted, err := crypt.Encrypt(plaintext)
	require.NoError(err, "Encrypt failed")

	decrypted, err := crypt.Decrypt(encrypted)
	require.NoError(err, "Decrypt failed")

	assert.Equal(t, plaintext, decrypted, "Decrypted text does not match original plaintext")
}
