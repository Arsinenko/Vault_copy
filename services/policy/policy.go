package Policy

import (
	"Vault_copy/db_operations"
	"Vault_copy/db_operations/cryptoOperation"
	"Vault_copy/db_operations/models"
	iappsrv "Vault_copy/internal"
	LogService "Vault_copy/services/log"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/pgtype"
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

const (
	CanChangeAppName = "can_change_app_name"
	CanChangeAppDesc = "can_change_app_desc"
)

// golang map for store default values
var DefaultRuleSetUser = map[string] bool {
	"CanChangeAppName": false,
	"CanChangeAppDescription" : false,
}
var DefaultRuleSetOwner = map[string] bool {
	"CanChangeAppName": true,
	"CanChangeAppDescription" : true,
}


// TODO ! Handle errors
func GetRules(UserID int32, AppID int32) map[string] bool {
	_log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(UserID) + string(AppID))));

	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[GetRules]::db_operations.InitDB()", _log_hash)
		return nil;
	}
	// get rules where UserID = UserID and AppID = AppID
	var policy models.Policy
	res := db.First(&policy, "user_id = ? AND app_id = ?", UserID, AppID)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[GetRules]::db.First()", _log_hash)
    return nil;
	}

	if policy.ID != 0 {
		return nil;
	}
	
	var rules map[string] bool
	if err := json.Unmarshal(policy.Rules.Bytes, &rules); err != nil {
		//LogService.Push_server_log(LogService.ErrorJSONUnmarshal, LogService.TErrorJSONUnmarshal, "[API_AppChangeName]::json.Unmarshal(policy.Rules)", _log_hash)
    return nil
	}

	return rules
}

// FINAL - STATIC API - 3 LAYER API
func CheckRule(UserID int32, AppID int32, Rule string) (bool, error) {
	_log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(UserID) + string(AppID))));

	app, e := iappsrv.I_get_app(AppID);
	if e != nil {
		LogService.Push_server_log(LogService.ErrorIGetApp, LogService.TErrorIGetApp, "[CheckRule]::AppSrv.I_get_app(AppID)", _log_hash)
		return false, e
	}
	if app == nil {
		return false, errors.New("404");
	}

	ruls := GetRules(UserID, AppID);
	if ruls == nil {
		if app.OwnerID == UserID {
			return DefaultRuleSetOwner[Rule] || false, nil;
		} else {
			return DefaultRuleSetUser[Rule] || false, nil;
		}
	} else {
		return ruls[Rule] || false, nil;
	}
}


// JSONB Interface for JSONB Field of yourTableName Table
type JSONB []interface{}

// Value Marshal
func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a);
}

// Scan Unmarshal
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}


func CreatePolicy(UserID int32, AppID int, Rules pgtype.JSONB) (int64, error) {
	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[CreatePolicy]::db_operations.InitDB()", "")
		return 0, e
	}
	// get rules where UserID = UserID and AppID = AppID
	var policy models.Policy
	res := db.First(&policy, "user_id = ? AND app_id = ?", UserID, AppID)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[CreatePolicy]::db.First()", "")
		return 0, res.Error
	}

	if policy.ID != 0 {
		return 0, nil
	}

	policy.UserID = UserID
	policy.AppID = AppID
	policy.Rules = Rules

	res = db.Create(&policy)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[CreatePolicy]::db.Create(&policy)", "")
		return 0, res.Error
	}
	return policy.ID, nil
}
func DeletePolicy(UserID int32, AppID int) error {
	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[DeletePolicy]::db_operations.InitDB()", "")
		return e
	}
	res := db.Delete(&models.Policy{}, "user_id = ? AND app_id = ?", UserID, AppID)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[DeletePolicy]::db.Delete(&models.Policy{})", "")
		return res.Error
	}
	return nil
}