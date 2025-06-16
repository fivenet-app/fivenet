package utils

import (
	"crypto/rand"
	"math/big"
)

// letters contains all alphanumeric characters for random string generation.
const (
	letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	numbers = "0123456789"
)

// GenerateRandomString returns a securely generated random string of length n using letters and digits.
// Taken from https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb
// It will return an error if the system's secure random number generator fails to function correctly.
func GenerateRandomString(n int) (string, error) {
	return GenerateRandomStringFromLetters(n, letters)
}

// GenerateRandomNumberString returns a securely generated random string of digits of length n.
func GenerateRandomNumberString(n int) (string, error) {
	return GenerateRandomStringFromLetters(n, numbers)
}

// GenerateRandomStringFromLetters returns a securely generated random string of length n using the provided letters set.
// Returns an error if the system's secure random number generator fails.
func GenerateRandomStringFromLetters(n int, letters string) (string, error) {
	ret := make([]byte, n)
	for i := range n {
		// Securely generate a random index for the letters set
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
