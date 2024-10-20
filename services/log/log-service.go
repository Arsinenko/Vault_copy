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

	ErrorCreateToken

	ErrorEnc
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

const TErrorCreateToken = "Failed to create token"

const TErrorEnc = "Failed to encrypt data"

// audit_log
func Push_server_log(type_l int16, msg string, stack string, hash string) {
	db, err := db_operations.InitDB()
	if err != nil {
		panic(err)
	}

	var log_e models.ServerLog
	log_e.MSG = msg
	log_e.Type = type_l
	log_e.Stack = stack
	log_e.Date = time.Now()
	// TODO add hash
	db.Create(&log_e) // TODO handle errors
}
