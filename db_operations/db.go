package db_operations

import (
	"Vault_copy/db_operations/cryptoOperation"
	"Vault_copy/db_operations/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strings"
)

func DB_connection() (*gorm.DB, error) {
	username := "vault"
	password := "test12344444"
	dbName := "vault"
	host := "localhost"
	port := 5432

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName)

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	db.LogMode(true) // Включает логирование SQL-запросов

	db.AutoMigrate(&models.User{})
	return db, nil
}

func CreateUser(db *gorm.DB, fullname string, phone string, email string, password string, TwoFactorKey []byte, Metadata json.RawMessage) {
	currentTime := time.Now()
	passwd, err := cryptoOperation.RegisterUser(password)
	if err != nil {
		fmt.Println(err)
	}

	var user = models.User{
		FullName:     fullname,
		PhoneNumber:  phone,
		Email:        email,
		Password:     string(passwd),
		CreationDate: currentTime,
		TwoFactorKey: TwoFactorKey,
		Metadata:     Metadata,
	}

	if err := db.Create(&user).Error; err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		return
	}
	fmt.Printf("Created user: %v\n", user)
}

func Auth(db *gorm.DB, identifier string, password byte) {
	if strings.Contains(identifier, "@") {
		var user models.User
		result := db.Model(&models.User{}).Where("email = ? AND password = ?", identifier, string(password)).First(&user)
		if result.RowsAffected == 1 {
			fmt.Printf("Found user: %v\n", user)
			//var token = db.Model(&models.securePassword{}).Where("token = ?", identifier).First(&user)
		}
		fmt.Printf("securePassword does not exist: %v\n", user)
		return
	}
	var user models.User
	result := db.Model(&models.User{}).Where("phone = ? AND password = ?", identifier, string(password)).First(&user)
	if result.RowsAffected == 1 {
		fmt.Printf("Found user: %v\n", user)
		return
	}
	fmt.Printf("securePassword does not exist: %v\n", user)
	return
}
