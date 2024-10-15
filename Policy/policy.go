package Policy

import (
	"Vault_copy/db_operations"
	"Vault_copy/db_operations/models"
	LogService "Vault_copy/services/log"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

const (
	readingPolicy  = "read"
	creatingPolicy = "create"
	updatingPolicy = "update"
	deletingPolicy = "delete"
	listingPolicy  = "list"
	sudoPolicy     = "sudo"
	patchingPolicy = "patch"
)

func GetRules(UserID int32, AppID int32, _log_hash string) []string {
	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[API_AppChangeName]::db_operations.InitDB()", _log_hash)

	}
	// get rules where UserID = UserID and AppID = AppID
	var policy models.Policy
	res := db.First(&policy, "user_id = ? AND app_id = ?", UserID, AppID)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[API_AppChangeName]::db_operations.InitDB()", _log_hash)

	}
	//check if policy is not empty
	if policy.ID != 0 {
		LogService.PushAuditLog(LogService.EventChangeAppName, UserID, AppID, 0, _log_hash)

	}
	var rules []string
	if err := json.Unmarshal(policy.Rules.Bytes, &rules); err != nil {
		//LogService.Push_server_log(LogService.ErrorJSONUnmarshal, LogService.TErrorJSONUnmarshal, "[API_AppChangeName]::json.Unmarshal(policy.Rules)", _log_hash)

	}

	return rules
}

// JSONB Interface for JSONB Field of yourTableName Table
type JSONB []interface{}

// Value Marshal
func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
