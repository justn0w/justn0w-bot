package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(str string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
