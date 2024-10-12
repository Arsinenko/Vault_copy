package db_operations

import (
	"Vault_copy/db_operations/models"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

func DB_connection() (*gorm.DB, error) {
	username := "your_username"
	password := "your_password"
	dbName := "your_database"
	host := "localhost"
	port := 5432

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName)

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateUser(db *gorm.DB, fullname string, phone string, email string, password string, TwoFactorKey []byte) {
	currentTime := time.Now()
	var user = models.User{
		FullName:     fullname,
		PhoneNumber:  phone,
		Email:        email,
		Password:     password,
		CreationDate: currentTime,
		TwoFactorKey: TwoFactorKey,
		Metadata:     nil,
	}
	db.Create(&user)
}
