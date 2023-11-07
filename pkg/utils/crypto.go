package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomHexBytes(length int) (string, error) {
	randomBytes, err := GenerateRandomBytes(length)
	hexString := hex.EncodeToString(randomBytes)
	return hexString, err
}
func GenerateRandomBytes(length int) ([]byte, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	return randomBytes, nil
}
