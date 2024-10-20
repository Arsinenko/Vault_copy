package db_operations

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)


var db_active *gorm.DB;
var db_init = false

func InitDB() (*gorm.DB, error) {
	if db_init { return db_active, nil; }

	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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