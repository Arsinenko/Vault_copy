package AuditLog

import (
	"Vault_copy/db_operations"
	"Vault_copy/db_operations/models"
	"time"
)

// TODO: Complete events
const (
	EventAuth int16 = iota + 1
	EventRegister
	EventSaltError
	EventDecodePasswdError
	EventUnauthorized
	EventCreateApp
	EventCreateSecret
)

// audit_log
func CreateAuditLog(action int16, idUser int32, AppId int32, TokenHash string) {

	db, err := db_operations.InitDB()
	if err != nil {
		panic(err)
	}

	var auditLog models.AuditLog
	auditLog.Action = action
	auditLog.UserID = idUser
	auditLog.Date = time.Now()
	auditLog.AppID = AppId
	auditLog.TokenHash = TokenHash

	db.Create(&auditLog)
}
