package secret

import (
	"Vault_copy/db_operations"
	"Vault_copy/db_operations/cryptoOperation"
	"Vault_copy/db_operations/models"
	LogService "Vault_copy/services/log"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

func CreateSecret(Data []byte, AppID int32, Metadata string) int { // SID, ya hz ne pomny zachem eto)))))
	logHash := hex.EncodeToString(cryptoOperation.SHA256(append([]byte(string(AppID)+Metadata), Data...))) // TODO fix FMT
	LogService.PushAuditLog(LogService.EventTryCreateSecret, 0, AppID, 0, logHash)

	db, err := db_operations.InitDB()
	if err != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[CreateSecret]::db_operations.InitDB()", logHash)
		return http.StatusInternalServerError
	}

	encryptedSecret, err := cryptoOperation.EncryptSecret(Data)
	if err != nil {
		// Обработка ошибки
		return http.StatusInternalServerError
	}

	var secret models.Secret
	secret.Data = []byte(encryptedSecret)
	secret.AppID = AppID
	secret.CreationDate = time.Now()
	secret.Metadata = Metadata

	db.Create(&secret) // TODO

	LogService.PushAuditLog(LogService.EventCreateSecret, 0, secret.AppID, secret.ID, logHash)
	return http.StatusOK
}

func DeleteSecret(SecretID int64, AppID int32) int {
	logHash := hex.EncodeToString(cryptoOperation.SHA256([]byte(strconv.FormatInt(SecretID, 10) + string(AppID))))
	LogService.PushAuditLog(LogService.EventTryDeleteSecret, 0, AppID, SecretID, logHash)

	db, err := db_operations.InitDB()
	if err != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[DeleteSecret]::db_operations.InitDB()", logHash)
		return http.StatusInternalServerError
	}

	var secret models.Secret
	res := db.First(&secret, "ID = ?", SecretID)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[DeleteSecret]::db_operations.InitDB()", logHash)
		return http.StatusInternalServerError
	}

	res = db.Delete(&secret)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[DeleteSecret]::db_operations.InitDB()", logHash)
		return http.StatusInternalServerError
	}

	LogService.PushAuditLog(LogService.EventDeleteSecret, 0, AppID, SecretID, logHash)
	return http.StatusOK
}

func GetSecrets(AppID int32) ([]models.Secret, int) {
	logHash := hex.EncodeToString(cryptoOperation.SHA256([]byte(strconv.FormatInt(int64(AppID), 10))))
	LogService.PushAuditLog(LogService.EventTryGetSecret, 0, AppID, 0, logHash)

	db, err := db_operations.InitDB()
	if err != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[GetSecrets]::db_operations.InitDB()", logHash)
		return nil, http.StatusInternalServerError
	}

	var secrets []models.Secret
	res := db.Where("app_id = ?", AppID).Find(&secrets)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[GetSecrets]::db.Find(&secrets)", logHash)
		return nil, http.StatusInternalServerError
	}

	var decryptedSecrets []models.Secret
	for _, secret := range secrets {
		decryptedData, err := cryptoOperation.DecryptSecret(string(secret.Data))
		if err != nil {
			// Логирование ошибки и пропуск этого секрета
			LogService.Push_server_log(LogService.ErrorGeneral, "Failed to decrypt secret", "[GetSecrets]::cryptoOperation.DecryptSecret()", logHash)
			continue
		}
		secret.Data = decryptedData
		decryptedSecrets = append(decryptedSecrets, secret)
	}

	LogService.PushAuditLog(LogService.EventGetSecret, 0, AppID, 0, logHash)
	return decryptedSecrets, http.StatusOK
}
