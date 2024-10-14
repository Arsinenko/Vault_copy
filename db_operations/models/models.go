package models

import (
	"time"

	"github.com/jackc/pgx/pgtype"
)

type Secret struct {
	ID           int64           		`db:"id"`
	SID          string          		`db:"sid"`
	Data         []byte          		`db:"data"`
	AppID        int32           		`db:"app_id"`
	CreationDate time.Time       		`db:"creation_date"`
	Metadata     pgtype.JSONB   		`db:"metadata"`
}
func (Secret) TableName() string {
	return "secret"
}

type App struct {
	ID           int32           		`db:"id"`
	Name         string          		`db:"name"`
	Description  string          		`db:"description"`
	OwnerID      int32           		`db:"owner_id"`
	CreationDate time.Time       		`db:"creation_date"`
	Metadata     pgtype.JSONB   		`db:"metadata"`
	APIPath      string          		`db:"api_path"`
}
func (App) TableName() string {
	return "app"
}

type User struct {
	ID           int32           		`db:"id"`
	FullName     string          		`db:"full_name"`
	PhoneNumber  string          		`db:"phone_number"`
	Email        string          		`db:"email"`
	TwoFactorKey string          		`db:"two_factor_key"`
	CreationDate time.Time       		`db:"creation_date"`
	Password     string          		`db:"password"`
	Metadata     pgtype.JSONB    		`db:"metadata"`
}
func (User) TableName() string {
	return "user"
}

type APIToken struct {
	ID           int       			 		`db:"id"`
	Token        string    			 		`db:"token"`
	AppID        int32     			 		`db:"app_id"`
	CreationDate time.Time 			 		`db:"creation_date"`
}
func (APIToken) TableName() string {
	return "api_token"
}

type AuditLog struct {
	ID				int64							 		`db:"id"`
	Action    int16     				 		`db:"action"`
	UserID    int32     				 		`db:"user_id"`
	AppID     int32     				 		`db:"app_id"`
	SecretID  int64     				 		`db:"secret_id"`
	TokenHash string    				 		`db:"token_hash"`
	Date      time.Time 				 		`db:"date"`
}
func (AuditLog) TableName() string {
	return "audit_log"
}

type Cert struct {
	ID           int             		`db:"id"`
	Public       string          		`db:"public"`
	Private      string          		`db:"private"`
	CreationDate pgtype.Timestamptz `db:"creation_date"`
	AppID        int32           		`db:"app_id"`
	Metadata     pgtype.JSONB  	 		`db:"metadata"`
}
func (Cert) TableName() string {
	return "cert"
}

type Policy struct {
	ID			 int64					     		`db:"id"`
	AppID    int                 		`db:"app_id"`
	UserID   int32               		`db:"user_id"`
	Rules    pgtype.JSONB           `db:"rules"`
	Metadata pgtype.JSONB           `db:"metadata"`
}
func (Policy) TableName() string {
	return "policy"
}

type ServerLog struct {
	ID			 int64					     		`db:"id"`
	MSG	     string                 `db:"msg"`
	Type     int16               		`db:"type"`
	Stack    string                 `db:"stack"`
	Date     time.Time  				 		`db:"date"`
}
func (ServerLog) TableName() string {
	return "server_log"
}

type SessionToken struct {
	UserID	 int64	    				    `db:"userid"`
	Hash	   string                 `db:"hash"`
	Date     time.Time              `db:"creation_date"`
	Metadata pgtype.JSONB           `db:"metadata"`
}
func (SessionToken) TableName() string {
	return "auth_token"
}
