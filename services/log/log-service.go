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
	ErrorDBExec
	ErrorDBSave
	ErrorDBDelete

	ErrorCreateApp
	ErrorCreateSecret
	ErrorGetUsr
	ErrorIGetApp
	ErrorRuleCheck

	ErrorISetAppName
	ErrorISetAppDesc

	ErrorIEncPolicy
	ErrorIDecPolicy

	ErrorJSONUnmarshal
	ErrorJSONMarshal
)

const TErrorHexDecode = "Error decode from hex"

const TErrorDBInit = "Failed to initialize database connection"
const TErrorDBExec = "Failed to execute querry"
const TErrorDBSave = "Failed to save module in DB"
const TErrorDBDelete = "Failed delete entry from database"

const TErrorCreateApp = "Failed to create app"
const TErrorCreateSecret = "Failed to push secret into database"
const TErrorGetUsr = "Failed to get user from database"
const TErrorIGetApp = "Failed to get app from database"
const TErrorRuleCheck = "Failed to check rules"

const TErrorISetAppName = "Failed to set app name"
const TErrorISetAppDesc = "Failed to set app description"

const TErrorIEncPolicy = "Failed to encode policy"
const TErrorIDecPolicy = "Failed to decode policy"

const TErrorJSONUnmarshal = "Failed to unmarshal json"
const TErrorJSONMarshal = "Failed to marshal json"

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
