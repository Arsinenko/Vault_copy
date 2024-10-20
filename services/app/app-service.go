package service_app

import (
	"Vault_copy/db_operations"
	"Vault_copy/db_operations/cryptoOperation"
	"Vault_copy/db_operations/models"
	IAPI "Vault_copy/internal"
	LogService "Vault_copy/services/log"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/pgtype"
)

// TODO - NOT COMPLETED --- security checks,
func CreateApp(Name string, Description string, OwnerID int32, metadata pgtype.JSONB) int {
	logHash := hex.EncodeToString(cryptoOperation.SHA256(append([]byte(Name+Description+fmt.Sprint(OwnerID)), metadata.Bytes...)))
	LogService.PushAuditLog(LogService.EventTryCreateApp, OwnerID, 0, 0, logHash)

	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[CreateApp]::db_operations.InitDB()", logHash)
		return http.StatusInternalServerError
	}
	var app models.App
	app.Name = Name
	app.Description = Description
	app.OwnerID = OwnerID
	app.Metadata = "{}"
	app.APIPath = strings.ToLower(strings.ReplaceAll(Name, " ", "_"))
	app.CreationDate = time.Now()
	db.Create(&app)

	LogService.PushAuditLog(LogService.EventCreateApp, app.OwnerID, app.ID, 0, logHash)

	return http.StatusOK
}

// TODO: security checks, check name text legit. last error
// FINAL - STATIC API - 2 LAYER API
func API_AppChangeName(UserID int32, AppID int32, name string) int {
	logHash := hex.EncodeToString(cryptoOperation.SHA256([]byte(fmt.Sprint(UserID) + fmt.Sprint(AppID) + name)))
	LogService.PushAuditLog(LogService.EventChangeAppNameTry, UserID, AppID, 0, logHash)

	rule, e := IAPI.I_get_policy_rule(AppID, UserID, IAPI.I_rule_change_app_name)
	if e != nil {
		LogService.Push_server_log(LogService.ErrorRuleCheck, LogService.TErrorRuleCheck, "[API_AppChangeName]::Policy.CheckRule()", logHash)
		return http.StatusInternalServerError
	}
	if !rule {
		LogService.PushAuditLog(LogService.EventChangeAppNameForbidden, UserID, AppID, 0, logHash)
		return http.StatusForbidden
	}

	app, e := IAPI.I_set_app_name(AppID, name)
	if e != nil {
		LogService.Push_server_log(LogService.ErrorISetAppName, LogService.TErrorISetAppName, "[API_AppChangeName]::iappsrv.I_set_app_name(AppID)", logHash)
		return http.StatusInternalServerError
	}
	if app == nil {
		LogService.PushAuditLog(LogService.EventChangeAppNameNotFound, UserID, AppID, 0, logHash)
		return http.StatusNotFound
	}

	LogService.PushAuditLog(LogService.EventChangeAppName, UserID, AppID, 0, logHash)
	return http.StatusOK
}

func API_AppGetName(UserID int32, AppID int32) (string, int) {
	logHash := hex.EncodeToString(cryptoOperation.SHA256([]byte(fmt.Sprint(UserID) + fmt.Sprint(AppID))))
	LogService.PushAuditLog(LogService.EventChangeAppNameTry, UserID, AppID, 0, logHash)

	rule, e := IAPI.I_get_policy_rule(AppID, UserID, IAPI.I_rule_view_app_info)
	if e != nil {
		LogService.Push_server_log(LogService.ErrorRuleCheck, LogService.TErrorRuleCheck, "[API_AppChangeName]::Policy.CheckRule()", logHash)
		return "", http.StatusInternalServerError
	}
	if !rule {
		LogService.PushAuditLog(LogService.EventChangeAppNameForbidden, UserID, AppID, 0, logHash)
		return "", http.StatusForbidden
	}

	res_name, e := IAPI.I_get_app_name(AppID)
	if e != nil {
		LogService.Push_server_log(LogService.ErrorISetAppName, LogService.TErrorISetAppName, "[API_AppChangeName]::iappsrv.I_set_app_name(AppID)", logHash)
		return "", http.StatusInternalServerError
	}
	if res_name == "" {
		LogService.PushAuditLog(LogService.EventChangeAppNameNotFound, UserID, AppID, 0, logHash)
		return "",  http.StatusNotFound
	}

	LogService.PushAuditLog(LogService.EventChangeAppName, UserID, AppID, 0, logHash)
	return res_name, http.StatusOK
}

// TODO: security checks, check description text legit
// FINAL - STATIC API - 2 LAYER API
func API_AppChangeDescription(UserID int32, AppID int32, description string) int {
	logHash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(UserID) + string(AppID) + description)))
	LogService.PushAuditLog(LogService.EventChangeAppDescTry, UserID, AppID, 0, logHash)

	rule, e := IAPI.I_get_policy_rule(AppID, UserID, IAPI.I_rule_change_app_desc)
	if e != nil {
		LogService.Push_server_log(LogService.ErrorRuleCheck, LogService.TErrorRuleCheck, "[API_AppChangeDescription]::Policy.CheckRule()", logHash)
		return http.StatusInternalServerError
	}
	if !rule {
		LogService.PushAuditLog(LogService.EventChangeAppDescForbidden, UserID, AppID, 0, logHash)
		return http.StatusForbidden
	}

	app, e := IAPI.I_set_app_desc(AppID, description)
	if e != nil {
		LogService.Push_server_log(LogService.ErrorISetAppDesc, LogService.TErrorISetAppDesc, "[API_AppChangeDescription]::iappsrv.I_set_app_desc(AppID)", logHash)
		return http.StatusInternalServerError
	}
	if app == nil {
		LogService.PushAuditLog(LogService.EventChangeAppDescNotFound, UserID, AppID, 0, logHash)
		return http.StatusNotFound
	}

	LogService.PushAuditLog(LogService.EventChangeAppDesc, UserID, AppID, 0, logHash)
	return http.StatusOK
}

