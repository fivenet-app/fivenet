package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	numbers = "0123456789"
)

// Taken from https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb
//
// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {
	return GenerateRandomStringFromLetters(n, letters)
}

func GenerateRandomNumberString(n int) (string, error) {
	return GenerateRandomStringFromLetters(n, numbers)
}

func GenerateRandomStringFromLetters(n int, letters string) (string, error) {
	ret := make([]byte, n)
	for i := range n {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
