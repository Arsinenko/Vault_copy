package cryptoOperation

import (
	"crypto/rand"
	"crypto/sha256"
)

// Generate random salt
func SALT(size int) []byte {
	salt := make([]byte, size);
	_, e := rand.Read(salt);
	if e != nil { panic(e); }
	return salt;
}

// SHA256 HASH
func SHA256(data []byte) []byte {
	hasher := sha256.New();
	hasher.Write(data);
	return hasher.Sum(nil);
}

// Хеширование пароля с двумя солями
// func HashPassword(password string, salt1, salt2 []byte) []byte {
// 	hasher := sha256.New()
// 	hasher.Write(salt1)
// 	hasher.Write([]byte(password))
// 	hasher.Write(salt2)
// 	return hasher.Sum(nil)
// }

// Функция для регистрации пользователя
// func RegisterUser(password string) (string, error) {
// 	// Генерация двух солей
// 	salt1 := gen_salt(16)
// 
// 	salt2 := gen_salt(16)
// 
// 	// Хешируем пароль
// 	passHash := HashPassword(password, salt1, salt2)
// 
// 	// Объединяем все части (salt1 + passHash + salt2) и возвращаем в hex-формате
// 	result := append(salt1, passHash...)
// 	result = append(result, salt2...)
// 
// 	return hex.EncodeToString(result), nil
// }

// func checkPassword(storedHash, password string) (bool, error) {
// 	// Преобразуем хранимый хэш из hex в байты
// 	storedBytes, err := hex.DecodeString(storedHash)
// 	if err != nil {
// 		return false, err
// 	}
// 
// 	// Разделяем на соль 1, хэш пароля и соль 2
// 	if len(storedBytes) != 64 {
// 		return false, errors.New("invalid stored hash length")
// 	}
// 	salt1 := storedBytes[:16]
// 	passHash := storedBytes[16:48]
// 	salt2 := storedBytes[48:]
// 
// 	// Хешируем введённый пароль с теми же солями
// 	calculatedHash := HashPassword(password, salt1, salt2)
// 
// 	// Сравниваем хэши
// 	for i := 0; i < 32; i++ {
//  	if passHash[i] != calculatedHash[i] {
//			return false, nil
//  	}
// 	}
// 	return true, nil
// }