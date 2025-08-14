package crypt

import (
	"testing"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
)

func TestEncrypt(t *testing.T) {
	cfg := &config.Config{Secret: "test-secret"}
	crypt := New(cfg)

	plaintext := "Hello, World!"
	encrypted, err := crypt.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if encrypted == "" {
		t.Fatal("Encrypt returned an empty string")
	}

	// Ensure the encrypted string is not the same as the plaintext
	if encrypted == plaintext {
		t.Fatal("Encrypted string should not match the plaintext")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	cfg := &config.Config{Secret: "test-secret"}
	crypt := New(cfg)

	plaintext := "Hello, World!"
	encrypted, err := crypt.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := crypt.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Fatalf(
			"Decrypted text does not match original plaintext. Got: %s, Want: %s",
			decrypted,
			plaintext,
		)
	}
}
