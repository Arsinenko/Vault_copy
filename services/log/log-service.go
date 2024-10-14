package LogService

import (
	"Vault_copy/db_operations"
	"Vault_copy/db_operations/models"
	"time"
)

// TODO: Complete events
const (
	ErrorGeneral int16 = iota + 1
	ErrorDatabase
	ErrorHexDecode
	ErrorDBInit
	ErrorCreateApp
)

const TErrorHexDecode = "Error decode from hex";
const TErrorDBInit = "Failed to initialize database connection";
const TErrorCreateApp = "Failed to create app";

// audit_log
func Push_server_log(type_l int16, msg string, stack string, hash string) {
	db, err := db_operations.InitDB()
	if err != nil {
		panic(err)
	}

	var auditLog models.ServerLog
	auditLog.MSG = msg
	auditLog.Type = type_l
	auditLog.Stack = stack
	auditLog.Date = time.Now()
	db.Create(&auditLog) // TODO handle errors 
}
