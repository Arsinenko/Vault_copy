package cryptoOperation

import (
	"golang.org/x/crypto/bcrypt"
)

// TODO check if it's necessary
// HashPassword хеширует пароль с использованием bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash сравнивает пароль с хешем
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
