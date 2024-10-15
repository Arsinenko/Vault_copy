package LogService

import (
	"Vault_copy/db_operations"
	"Vault_copy/db_operations/models"
	"time"
)

// TODO: Complete events
const (
	EventAuth int16 = iota + 1
	EventTryAuth

	EventRegister
	EventTryRegister

	EventTryRegisterAlreadyExists

	EventTryRegisterBadPassword
	EventTryRegisterBadLogin

	EventSaltError
	EventDecodePasswdError
	EventUnauthorized

	EventCreateApp
	EventTryCreateApp

	EventCreateSecret
	EventTryCreateSecret

	EventChangeAppName
	EventChangeAppNameTry
	EventChangeAppNameForbidden
	EventChangeAppNameNotFound

	EventChangeAppDesc
	EventChangeAppDescTry
	EventChangeAppDescForbidden
	EventChangeAppDescNotFound

	EventTryChangeAppDescription
	EventChangeAppDescription
)

// audit_log
func PushAuditLog(action int16, idUser int32, AppId int32, SecretID int64, hash string) {
	db, err := db_operations.InitDB()
	if err != nil {
		panic(err)
	}

	var auditLog models.AuditLog
	auditLog.Action = action
	auditLog.UserID = idUser
	auditLog.Date = time.Now()
	auditLog.SecretID = SecretID
	auditLog.AppID = AppId
	auditLog.TokenHash = hash

	db.Create(&auditLog)
}
