package cryptoOperation

import (
	"crypto/rand"
	"crypto/sha256"
)

// Generate random salt
func SALT(size int) []byte {
	salt := make([]byte, size)
	_, e := rand.Read(salt)
	if e != nil {
		panic(e)
	}
	return salt
}

// SHA256 HASH
func SHA256(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}