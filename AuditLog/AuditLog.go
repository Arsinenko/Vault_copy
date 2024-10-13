package AuditLog

import (
	"Vault_copy/db_operations"
	"Vault_copy/db_operations/models"
	"time"
)

const (
	EventAuth int16 = iota + 1
	EventRegister
)

// audit_log
func CreateAuditLog(action int16, idUser int32) {
	//actions register: 1, auth: 2
	db, err := db_operations.InitDB()
	if err != nil {
		panic(err)
	}
	var auditLog models.AuditLog
	auditLog.Action = action
	auditLog.UserID = idUser
	auditLog.Date = time.Now()
	db.Create(&auditLog)
}
