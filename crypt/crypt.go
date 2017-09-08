package crypt

import (
	"golang.org/x/crypto/bcrypt"
)

var (
	// DefaultHashComplexity defines the default complexity to be using for hashing
	// using the dccrypt hashing method.
	DefaultHashComplexity = 10
)

// BcryptAuthenticate attempts to validate expected value which is already encrypted using
// the bcrypt hash.The expected value is the already hashed value and the provided
// wukk be hashed and compared to validate if its a valid value.
func BcryptAuthenticate(expected, provided []byte) error {
	return bcrypt.CompareHashAndPassword(expected, provided)
}

// BcryptGenerate returns a value encrypted using the bcrypt hashing algorithmn, it takes
// all provided values to generate final output.
func BcryptGenerate(content []byte, hashComplexity int) ([]byte, error) {
	if hashComplexity <= 0 {
		hashComplexity = DefaultHashComplexity
	}

	return bcrypt.GenerateFromPassword(content, hashComplexity)
}
