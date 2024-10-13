package db_operations

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)


var db_active *gorm.DB;
var db_init = false

func InitDB() (*gorm.DB, error) {
	if db_init { return db_active, nil; }

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

	db_active = db;
	db_init = true;
	db.LogMode(true) // Включает логирование SQL-запросов

	//db.AutoMigrate(&models.User{})
	//db.AutoMigrate(&models.APIToken{})
	//db.AutoMigrate(&models.App{})
	//db.AutoMigrate(&models.AuditLog{})
	//db.AutoMigrate(&models.Cert{})
	//db.AutoMigrate(&models.Policy{})
	//db.AutoMigrate(&models.Secret{})

	return db, nil
}