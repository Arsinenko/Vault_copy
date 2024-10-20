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
)

// import "C"

func passHash(pass string, salt1 []byte, salt2 []byte) []byte {
	passb := []byte(pass)
	s := salt1[:]
	s = append(s, passb[:]...)
	s = append(s, salt2[:]...)
	return cryptoOperation.SHA256(s)
}
func pass_cmpP(hash1 []byte, hash2 []byte) int {
	return libc.MemCmpP(unsafe.Pointer(&hash1[0]), unsafe.Pointer(&hash2[0]), 32)
}

func AuthToken() {
	// TODO
}

// FINAL - STATIC API
func get_usr(phoneMail string) (*models.User, error) {
	logHash := hex.EncodeToString(cryptoOperation.SHA256([]byte(phoneMail)))

	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[get_usr]::db_operations.InitDB()", logHash)
		return nil, e
	}

	isMail := strings.IndexByte(phoneMail, '@') != -1
	optMailPhone := map[bool]string{true: "email", false: "phone_number"}

	var user *models.User
	res := db.First(&user, optMailPhone[isMail]+"= ?", phoneMail)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[get_usr]::db_operations.InitDB()", logHash)
		return nil, res.Error
	}
	if user == nil {
		return nil, nil;
	}
	
	return user, nil
}

// TODO: security checks
// MVP READY, STATIC API, FINAL (REST)
func AuthStandard(phoneMail string, password string) int {
	logHash := hex.EncodeToString(cryptoOperation.SHA256([]byte(phoneMail + password)))
	LogService.PushAuditLog(LogService.EventTryAuth, 0, 0, 0, logHash)

	user, err := get_usr(phoneMail)
	if err != nil {
		LogService.Push_server_log(LogService.ErrorGetUsr, LogService.TErrorGetUsr, "[CreateSecret]::get_usr(phone_mail)", logHash)
		return http.StatusNotFound
	}

	//usrSalt1, errS1 := hex.DecodeString(user.Password[:32])
	//if errS1 != nil {
	//	LogService.Push_server_log(LogService.ErrorHexDecode, LogService.TErrorHexDecode, "[AuthStandard]::hex:decode(usr_salt1)", logHash)
	//	return http.StatusInternalServerError
	//}
	//
	//usrHash, errH := hex.DecodeString(user.Password[32:96])
	//if errH != nil {
	//	LogService.Push_server_log(LogService.ErrorHexDecode, LogService.TErrorHexDecode, "[AuthStandard]::hex:decode(usr_hash)", logHash)
	//	return http.StatusInternalServerError
	//}
	//
	//usrSalt2, errS2 := hex.DecodeString(user.Password[96:])
	//if errS2 != nil {
	//	LogService.Push_server_log(LogService.ErrorHexDecode, LogService.TErrorHexDecode, "[AuthStandard]::hex:decode(usr_salt2)", logHash)
	//	return http.StatusInternalServerError
	//}
	//
	//rnHash := passHash(password, usrSalt1, usrSalt2)
	//authOk := pass_cmpP(usrHash, rnHash) == 0
	//TODO - check security
	authOk := cryptoOperation.CheckPasswordHash(password, user.Password)

	if authOk {
		LogService.PushAuditLog(LogService.EventAuth, user.ID, 0, 0, logHash)
		return http.StatusOK
	} else {
		LogService.PushAuditLog(LogService.EventUnauthorized, user.ID, 0, 0, logHash)
		return http.StatusUnauthorized
	}
}

// TODO: security checks, check password and phone/mail legit by content.
// MVP READY, STATIC API, FINAL (REST)
func Register(phoneMail string, password string, fullName string) int {
	logHash := hex.EncodeToString(cryptoOperation.SHA256([]byte(phoneMail + password + fullName)))
	LogService.PushAuditLog(LogService.EventTryRegister, 0, 0, 0, logHash)

	_, err := get_usr(phoneMail)
	if err == nil {
		// TODO LOG
		//LogService.Push_server_log(LogService.ErrorGetUsr, LogService.TErrorGetUsr, "[Register]::get_usr(phone_mail)", _log_hash)
		return http.StatusConflict
	}

	// ! ERROR <- broken
	// if (strings.Compare(act_usr.Email, phone_mail) == 0 || strings.Compare(act_usr.PhoneNumber, phone_mail) == 0) {
	// 	LogService.PushAuditLog(LogService.EventTryRegisterAlreadyExists, 0, 0, 0, _log_hash)
	// 	return http.StatusConflict
	// }
	// <--- STOPED HERE

	if len(password) < 8 {
		LogService.PushAuditLog(LogService.EventTryRegisterBadPassword, 0, 0, 0, logHash)
		return http.StatusBadRequest
	}

	if len(phoneMail) < 5 {
		LogService.PushAuditLog(LogService.EventTryRegisterBadLogin, 0, 0, 0, logHash)
		return http.StatusBadRequest
	}

	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[Register]::db_operations.InitDB()", logHash)
		return http.StatusInternalServerError
	}

	isMail := strings.IndexByte(phoneMail, '@') != -1

	//var salt1 = cryptoOperation.SALT(16)
	//var salt2 = cryptoOperation.SALT(16)
	//var hash = passHash(password, salt1, salt2)
	//
	//var s = salt1[:]
	//s = append(s, hash[:]...)
	//s = append(s, salt2[:]...)

	hashedPass, err := cryptoOperation.HashPassword(password)

	if err != nil {
		//LogService.Push_server_log(LogService.ErrorHashPassword, LogService.TErrorHashPassword, "[Register]::cryptoOperation.HashPassword()", logHash)
		return http.StatusInternalServerError
	}

	var user models.User
	user.CreationDate = time.Now()
	user.Password = hashedPass
	user.FullName = fullName
	user.Metadata = "{}"
	user.TwoFactorKey = "nil"

	if isMail {
		user.Email = phoneMail
		user.PhoneNumber = "nil"
	} else {
		user.PhoneNumber = phoneMail
		user.Email = "nil"
	}

	db.Create(&user)

	LogService.PushAuditLog(LogService.EventRegister, user.ID, 0, 0, logHash)
	return http.StatusOK
}

func DeleteUser(UserID int32, AppID int32) int {
	logHash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(UserID) + string(AppID))))

	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[DeleteUser]::db_operations.InitDB()", logHash)
		return http.StatusInternalServerError
	}

	var user models.User

	res := db.First(&user, "ID = ?", UserID)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[DeleteUser]::db_operations.InitDB()", logHash)
		return http.StatusInternalServerError
	}

	res = db.Delete(&user)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[DeleteUser]::db_operations.InitDB()", logHash)
		return http.StatusInternalServerError
	}

	LogService.PushAuditLog(LogService.EventDeleteUser, UserID, AppID, 0, logHash)
	return http.StatusOK
}

func MakeToken(userID int32) {

}
