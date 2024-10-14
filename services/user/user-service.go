package service_user

import (
	"Vault_copy/db_operations"
	"Vault_copy/db_operations/cryptoOperation"
	"Vault_copy/db_operations/models"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
	"unsafe"

	LogService "Vault_copy/services/log"

	"github.com/gotranspile/runtimec/libc"
	"github.com/jackc/pgx/pgtype"
)

// import "C"

func passHash(pass string, salt1 []byte, salt2 []byte) []byte {
	passb := []byte(pass)
	s := salt1[:]
	s = append(s, passb[:]...)
	s = append(s, salt2[:]...)
	return cryptoOperation.SHA256(s)
}
func pass_cmpP(hash_1 []byte, hash_2 []byte) int {
	return libc.MemCmpP(unsafe.Pointer(&hash_1[0]), unsafe.Pointer(&hash_2[0]), 32)
}

func AuthToken() {
 // TODO
}

func get_usr(phone_mail string) (models.User, error) {
	db, e := db_operations.InitDB()
	if e != nil {
		panic(e)
	}

	is_mail := strings.IndexByte(phone_mail, '@') != -1
	opt_mail_phone := map[bool]string{true: "email", false: "phone_number"}

	var user models.User
	res := db.First(&user, opt_mail_phone[is_mail]+"= ?", phone_mail)
	if res.Error != nil {
		return user, res.Error
	}
	return user, nil
}

func AuthStandard(phone_mail string, password string) int {
	_log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(phone_mail + password)));
	LogService.PushAuditLog(LogService.EventTryAuth, 0, 0, 0, _log_hash);

	user, err := get_usr(phone_mail)
	if err != nil {
		return http.StatusNotFound
	}

	usr_salt1, err_s1 := hex.DecodeString(user.Password[:32])
	if err_s1 != nil {
		LogService.Push_server_log(LogService.ErrorHexDecode, LogService.TErrorHexDecode, "[AuthStandard]::hex:decode(usr_salt1)", _log_hash)
		panic(err_s1)
	}

	usr_hash, err_h := hex.DecodeString(user.Password[32:96])
	if err_h != nil {
		LogService.Push_server_log(LogService.ErrorHexDecode, LogService.TErrorHexDecode, "[AuthStandard]::hex:decode(usr_hash)", _log_hash)
		panic(err_h)
	}

	usr_salt2, err_s2 := hex.DecodeString(user.Password[96:])
	if err_s2 != nil {
		LogService.Push_server_log(LogService.ErrorHexDecode, LogService.TErrorHexDecode, "[AuthStandard]::hex:decode(usr_salt2)", _log_hash)
		panic(err_s2)
	}

	rn_hash := passHash(password, usr_salt1, usr_salt2)
	auth_ok := pass_cmpP(usr_hash, rn_hash) == 0

	if auth_ok {
		LogService.PushAuditLog(LogService.EventAuth, user.ID, 0, 0, _log_hash)
		return http.StatusOK
	} else {
		LogService.PushAuditLog(LogService.EventUnauthorized, user.ID, 0, 0, _log_hash)
		return http.StatusUnauthorized
	}
}

func CreateUser(phone_mail string, password string, full_name string) int {
	_log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(phone_mail + password + full_name)));
	LogService.PushAuditLog(LogService.EventTryRegister, 0, 0, 0, _log_hash);

	act_usr, err := get_usr(phone_mail)
	if err == nil {
		return http.StatusInternalServerError
	} // [http]::Conflict - User with that phone or mail already exists
	if act_usr.Email == phone_mail || act_usr.PhoneNumber == phone_mail {
		return http.StatusConflict
	}

	// TODO: 1. Check password and phone/mail legit. 2. Check pass length

	if len(password) < 8 {
		return http.StatusBadRequest
	}

	if len(phone_mail) < 3 {
		return http.StatusBadRequest
	}

	db, e := db_operations.InitDB()
	if e != nil {
		panic(e)
	}

	is_mail := strings.IndexByte(phone_mail, '@') != -1

	var salt_1 = cryptoOperation.SALT(16)
	var salt_2 = cryptoOperation.SALT(16)
	var hash = passHash(password, salt_1, salt_2)

	var s = salt_1[:]
	s = append(s, hash[:]...)
	s = append(s, salt_2[:]...)

	var user models.User
	user.CreationDate = time.Now()
	user.Password = hex.EncodeToString(s)
	user.Metadata = pgtype.JSONB{}
	user.FullName = full_name

	if is_mail {
		user.Email = phone_mail
		user.PhoneNumber = "nil"
	} else {
		user.PhoneNumber = phone_mail
		user.Email = "nil"
	}

	db.Create(user)
	LogService.PushAuditLog(LogService.EventRegister, user.ID, 0, 0, _log_hash)

	return http.StatusOK
}

func CreateApp(Name string, Description string, OwnerID int32, metadata pgtype.JSONB) int {
	_log_hash := hex.EncodeToString(cryptoOperation.SHA256(append([]byte(Name + Description + string(OwnerID)), metadata.Bytes...)));
	LogService.PushAuditLog(LogService.EventTryCreateApp, OwnerID, 0, 0, _log_hash);

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

func CreateSecret(SID string, Data []byte, AppID int32, Metadata pgtype.JSONB) int { // SID, ya hz ne pomny zachem eto)))))
	_log_hash := hex.EncodeToString(append(cryptoOperation.SHA256(append([]byte(SID+string(AppID)), Data...)), Metadata.Bytes...));
	LogService.PushAuditLog(LogService.EventTryCreateSecret, 0, AppID, 0, _log_hash);

	db, err := db_operations.InitDB()
	if err != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[CreateSecret]::db_operations.InitDB()", _log_hash)
		return http.StatusInternalServerError
	}

	var secret models.Secret
	secret.SID = SID
	secret.Data = Data
	secret.AppID = AppID
	secret.CreationDate = time.Now()
	secret.Metadata = Metadata
	err = db.Create(&secret).Error
	if err != nil {
		LogService.Push_server_log(LogService.ErrorCreateSecret, LogService.TErrorCreateSecret, "[CreateSecret]::db.Create(&secret)", _log_hash)
		return http.StatusInternalServerError
	}

	LogService.PushAuditLog(LogService.EventCreateSecret, 0, secret.AppID, secret.ID, _log_hash)
	return http.StatusOK
}

func SecUser() {

}

func MakeToken(userID int32) {

}
