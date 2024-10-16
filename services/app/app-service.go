package service_app

import (
	"Vault_copy/db_operations"
	"Vault_copy/db_operations/cryptoOperation"
	"Vault_copy/db_operations/models"
	IAPI "Vault_copy/internal"
	LogService "Vault_copy/services/log"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/pgtype"
)

// TODO - NOT COMPLETED --- security checks,
func CreateApp(Name string, Description string, OwnerID int32, metadata pgtype.JSONB) int {
	_log_hash := hex.EncodeToString(cryptoOperation.SHA256(append([]byte(Name+Description+string(OwnerID)), metadata.Bytes...)))
	LogService.PushAuditLog(LogService.EventTryCreateApp, OwnerID, 0, 0, _log_hash)

	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[CreateApp]::db_operations.InitDB()", _log_hash)
		return http.StatusInternalServerError
	}
	var app models.App
	app.Name = Name
	app.Description = Description
	app.OwnerID = OwnerID
	app.Metadata = metadata
	app.APIPath = strings.ToLower(strings.ReplaceAll(Name, " ", "_"))
	app.CreationDate = time.Now()
	err := db.Create(&app).Error
	if err != nil {
		LogService.Push_server_log(LogService.ErrorCreateApp, LogService.TErrorCreateApp, "[CreateApp]::db.Create(&app)", _log_hash)
		return http.StatusInternalServerError
	}
	LogService.PushAuditLog(LogService.EventCreateApp, app.OwnerID, app.ID, 0, _log_hash)

	return http.StatusOK
}



// TODO: security checks, check name text legit
// FINAL - STATIC API - 2 LAYER API
func API_AppChangeName(UserID int32, AppID int32, name string) int {
	_log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(UserID) + string(AppID) + name)))
	LogService.PushAuditLog(LogService.EventChangeAppNameTry, UserID, AppID, 0, _log_hash)

	rule, e := IAPI.I_get_policy_rule(AppID, UserID, IAPI.I_rule_change_app_name)
	if e != nil {
		LogService.Push_server_log(LogService.ErrorRuleCheck, LogService.TErrorRuleCheck, "[API_AppChangeName]::Policy.CheckRule()", _log_hash)
		return http.StatusInternalServerError;
	}
	if !rule {
		LogService.PushAuditLog(LogService.EventChangeAppNameForbidden, UserID, AppID, 0, _log_hash)
		return http.StatusForbidden;
	}

	app, e := IAPI.I_set_app_name(AppID, name)
	if e != nil {
		LogService.Push_server_log(LogService.ErrorISetAppName, LogService.TErrorISetAppName, "[API_AppChangeName]::iappsrv.I_set_app_name(AppID)", _log_hash)
		return http.StatusInternalServerError;
	}
	if app == nil {
		LogService.PushAuditLog(LogService.EventChangeAppNameNotFound, UserID, AppID, 0, _log_hash)
		return http.StatusNotFound;
	}

	LogService.PushAuditLog(LogService.EventChangeAppName, UserID, AppID, 0, _log_hash)
	return http.StatusOK;
}

// TODO: security checks, check description text legit
// FINAL - STATIC API - 2 LAYER API
func API_AppChangeDescription(UserID int32, AppID int32, description string) int {
  _log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(UserID) + string(AppID) + description)))
	LogService.PushAuditLog(LogService.EventChangeAppDescTry, UserID, AppID, 0, _log_hash)

	rule, e := IAPI.I_get_policy_rule(AppID, UserID, IAPI.I_rule_change_app_desc)
	if e != nil {
		LogService.Push_server_log(LogService.ErrorRuleCheck, LogService.TErrorRuleCheck, "[API_AppChangeDescription]::Policy.CheckRule()", _log_hash)
		return http.StatusInternalServerError;
	}
	if !rule {
		LogService.PushAuditLog(LogService.EventChangeAppDescForbidden, UserID, AppID, 0, _log_hash)
		return http.StatusForbidden;
	}

	app, e := IAPI.I_set_app_desc(AppID, description)
	if e != nil {
		LogService.Push_server_log(LogService.ErrorISetAppDesc, LogService.TErrorISetAppDesc, "[API_AppChangeDescription]::iappsrv.I_set_app_desc(AppID)", _log_hash)
		return http.StatusInternalServerError;
	}
	if app == nil {
		LogService.PushAuditLog(LogService.EventChangeAppDescNotFound, UserID, AppID, 0, _log_hash)
		return http.StatusNotFound;
	}

	LogService.PushAuditLog(LogService.EventChangeAppDesc, UserID, AppID, 0, _log_hash)
	return http.StatusOK;
}