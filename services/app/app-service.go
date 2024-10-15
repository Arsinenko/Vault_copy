package service_app

import (
	"Vault_copy/db_operations"
	"Vault_copy/db_operations/cryptoOperation"
	"Vault_copy/db_operations/models"
	LogService "Vault_copy/services/log"
	"encoding/hex"
	"encoding/json"
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

// FINAL - STATIC API
func I_get_app(ID int32) *models.App {
	_log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(ID))))

	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[I_get_app]::db_operations.InitDB()", _log_hash)
		return nil
	}

	var app models.App
	res := db.First(&app, "ID = ?", ID)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[I_get_app]::db_operations.InitDB()", _log_hash)
		return nil
	}

	return &app
}

// TODO: check user policy, security checks, audit log
func API_AppChangeName(UserID int32, AppID int32, name string) {

	_log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(AppID) + name)))
	LogService.PushAuditLog(LogService.EventTryChangeAppName, UserID, AppID, 0, _log_hash)

	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[API_AppChangeName]::db_operations.InitDB()", _log_hash)
		return
	}
	// get rules where UserID = UserID and AppID = AppID
	var policy models.Policy
	res := db.First(&policy, "user_id = ? AND app_id = ?", UserID, AppID)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[API_AppChangeName]::db_operations.InitDB()", _log_hash)
		return
	}
	//check if policy is not empty
	if policy.ID != 0 {
		LogService.PushAuditLog(LogService.EventChangeAppName, UserID, AppID, 0, _log_hash)
		return
	}
	var rules []string
	if err := json.Unmarshal(policy.Rules.Bytes, &rules); err != nil {
		//LogService.Push_server_log(LogService.ErrorJSONUnmarshal, LogService.TErrorJSONUnmarshal, "[API_AppChangeName]::json.Unmarshal(policy.Rules)", _log_hash)
		return
	}
	AppChangeName(AppID, name)
}

func AppChangeDescription(UserID int32, AppID int32, description string) int {
	_log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(AppID) + description)))

	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[AppChangeDescription]::db_operations.InitDB()", _log_hash)
		return http.StatusInternalServerError
	}

	var app models.App
	res := db.First(&app, "ID = ?", AppID)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[API_AppChangeDescription]::db_operations.InitDB()", _log_hash)
		return http.StatusInternalServerError
	}
	app.Description = description
	err := db.Save(&app).Error
	if err != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[API_AppChangeDescription]::db.Save(&app)", _log_hash)
		return http.StatusInternalServerError
	}
	LogService.PushAuditLog(LogService.EventChangeAppDescription, UserID, AppID, 0, _log_hash)
	return http.StatusOK
}

// TODO
func AppChangeName(AppID int32, name string) int {
	_log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(AppID) + name)))

	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[AppChangeName]::db_operations.InitDB()", _log_hash)
		return http.StatusInternalServerError
	}

	var app models.App
	res := db.First(&app, "ID = ?", AppID)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[API_AppChangeName]::db_operations.InitDB()", _log_hash)
		return http.StatusInternalServerError
	}
	app.Name = name
	err := db.Save(&app).Error
	if err != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[API_AppChangeName]::db.Save(&app)", _log_hash)
		return http.StatusInternalServerError
	}
	LogService.PushAuditLog(LogService.EventChangeAppName, 0, AppID, 0, _log_hash)
	return http.StatusOK
}

//db.Exec("")
//panic("NotImplemented")
//}
