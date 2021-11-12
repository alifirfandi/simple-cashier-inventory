package helper

import "golang.org/x/crypto/bcrypt"

func CompareHash(hash, text string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))

	return err == nil
}
