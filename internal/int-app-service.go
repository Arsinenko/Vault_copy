package internal_api

import (
	"Vault_copy/db_operations"
	"Vault_copy/db_operations/cryptoOperation"
	"Vault_copy/db_operations/models"
	LogService "Vault_copy/services/log"
	"encoding/hex"
)

// * FINAL - STATIC API - 3 LAYER API
func I_get_app(ID int32) (*models.App, error) {
	logHash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(ID))))

	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[I_get_app]::db_operations.InitDB()", logHash)
		return nil, e
	}

	var app models.App
	res := db.First(&app, "ID = ?", ID)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[I_get_app]::db_operations.InitDB()", logHash)
		return nil, e
	}

	return &app, nil
}

// * FINAL - STATIC API - 3 LAYER API
func I_set_app_name(ID int32, Name string) (*models.App, error) {
	logHash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(ID) + Name)))

	app, e := I_get_app(ID)
	if e != nil {
		LogService.Push_server_log(LogService.ErrorIGetApp, LogService.TErrorIGetApp, "[I_app_name]::I_get_app(ID)", logHash)
		return nil, e
	}
	if app == nil {
		return nil, nil
	}

	app.Name = Name

	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[I_app_name]::db_operations.InitDB()", logHash)
		return nil, e
	}
	db.Save(app) // TODO Handle error

	return app, nil
}

// * FINAL - STATIC API - 3 LAYER API
func I_set_app_desc(ID int32, Desription string) (*models.App, error) {
	logHash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(ID) + Desription)))

	app, e := I_get_app(ID)
	if e != nil {
		LogService.Push_server_log(LogService.ErrorIGetApp, LogService.TErrorIGetApp, "[I_set_app_desc]::I_get_app(ID)", logHash)
		return nil, e
	}
	if app == nil {
		return nil, nil
	}

	app.Description = Desription

	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[I_set_app_desc]::db_operations.InitDB()", logHash)
		return nil, e
	}
	db.Save(app) // TODO Handle error

	return app, nil
}

// ! BEFORE <-- move policy into internal api segment
// func I_app_add_user(UserID int32, AppID int32, ) *models.Policy{
// 	_log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(UserID) + string(AppID))))
// 	res, e := policy.CreatePolicy(UserID, AppID)
// 	if e != nil {
// 		LogService.Push_server_log(LogService.ErrorIGetApp, LogService.TErrorIGetApp, "[I_app_name]::I_get_app(ID)", _log_hash)
// 		return nil
// 	}
// 	return res
// }
// func I_app_remove_user(UserID int32, AppID int32, ) *models.Policy{
// 	_log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(UserID) + string(AppID))))
// 	res, e := policy.RemovePolicy(UserID, AppID)
// 	if e != nil {
// 		LogService.Push_server_log(LogService.ErrorIGetApp, LogService.TErrorIGetApp, "[I_app_name]::I_get_app(ID)", _log_hash)
// 		return nil
// 	}
// 	return res
// }

// I_app_description
// I_app_add_user (!policy)
// I_app_remove_user (!policy)
