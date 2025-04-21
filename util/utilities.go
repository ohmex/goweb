package util

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// function to check given string is in array or not
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func GenerateAPIKeySecret() (string, string, string) {
	key := make([]byte, 16)
	secret := make([]byte, 32)
	rand.Read(key)
	rand.Read(secret)

	apiKey := hex.EncodeToString(key)
	apiSecret := hex.EncodeToString(secret)
	hashedSecret, _ := bcrypt.GenerateFromPassword([]byte(apiSecret), bcrypt.DefaultCost)

	return apiKey, apiSecret, string(hashedSecret)
}
