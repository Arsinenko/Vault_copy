package models

import (
	"encoding/json"
	"time"
)

type Secret struct {
	ID           int       `db:"id"`
	SID          string    `db:"sid"`
	Data         []byte    `db:"data"`
	AppID        int64     `db:"app_id"`
	CreationDate time.Time `db:"creation_date"`
	Metadata     []byte    `db:"metadata"`
}

type App struct {
	ID           int       `db:"id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	OwnerID      int64     `db:"owner_id"`
	CreationDate time.Time `db:"creation_date"`
	Metadata     []byte    `db:"metadata"`
	APIPath      string    `db:"api_path"`
}

type User struct {
	ID           int             `db:"id"`
	FullName     string          `db:"full_name"`
	PhoneNumber  string          `db:"phone_number"`
	Email        string          `db:"email"`
	TwoFactorKey []byte          `db:"two_factor_key"`
	CreationDate time.Time       `db:"creation_date"`
	Password     string          `db:"password"`
	Metadata     json.RawMessage `db:"metadata"`
}

type APIToken struct {
	ID           int       `db:"id"`
	Token        []byte    `db:"token"`
	AppID        int64     `db:"app_id"`
	CreationDate time.Time `db:"creation_date"`
}

type AuditLog struct {
	ID        int       `db:"id"`
	Action    int16     `db:"action"`
	UserID    int64     `db:"user_id"`
	AppID     int64     `db:"app_id"`
	SecretID  int64     `db:"secret_id"`
	TokenHash []byte    `db:"token_hash"`
	Date      time.Time `db:"date"`
}

type Cert struct {
	ID           int       `db:"id"`
	Public       []byte    `db:"public"`
	Private      []byte    `db:"private"`
	CreationDate time.Time `db:"creation_date"`
	AppID        int64     `db:"app_id"`
	Metadata     []byte    `db:"metadata"`
}

type Policy struct {
	AppID    int64  `db:"app_id"`
	UserID   int64  `db:"user_id"`
	Rules    []byte `db:"rules"`
	Metadata []byte `db:"metadata"`
}
