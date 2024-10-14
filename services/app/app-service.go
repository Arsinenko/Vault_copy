package service_app

import (
	"Vault_copy/db_operations"
	"Vault_copy/db_operations/cryptoOperation"
	"Vault_copy/db_operations/models"
	LogService "Vault_copy/services/log"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/pgtype"
)

// TODO - NOT COMPLETED
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