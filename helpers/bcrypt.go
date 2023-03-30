package helpers

import "golang.org/x/crypto/bcrypt"

const salt = 8

func HashPass(p string) (string, error) {
	password := []byte(p)
	hash, err := bcrypt.GenerateFromPassword(password, salt)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePass(h, p []byte) bool {
	err := bcrypt.CompareHashAndPassword(h, p)
	return err == nil
}
