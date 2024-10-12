package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DBConfig представляет собой конфигурацию подключения к базе данных
type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

// NewDBConfig возвращает новую конфигурацию подключения к базе данных
func NewDBConfig(host, port, username, password, dbName string) *DBConfig {
	return &DBConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		DBName:   dbName,
	}
}

// ConnectDB подключается к базе данных с использованием предоставленной конфигурации
func ConnectDB(cfg *DBConfig) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	cfg := NewDBConfig("localhost", "5432", "username", "password", "dbname")
	db, err := ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Используйте подключение к базе данных
}
