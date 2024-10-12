package main

import (
	"crypto/sha256"
)

func encryption(data []byte) []byte {
	// Создайте новый объект хеша SHA256
	hash := sha256.New()
	hash.Write(data)
	hashValue := hash.Sum(nil)
	return hashValue
}


func main() {
	hashValue :=

}
