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

	"Vault_copy/AuditLog"
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

func AuthStandart(phone_mail string, password string) int {
	user, err := get_usr(phone_mail)
	if err != nil {
		return http.StatusNotFound
	}

	usr_salt1, err_s1 := hex.DecodeString(user.Password[:32])
	if err_s1 != nil {
		panic(err_s1)
	}
	usr_hash, err_h := hex.DecodeString(user.Password[32:96])
	if err_h != nil {
		panic(err_h)
	}
	usr_salt2, err_s2 := hex.DecodeString(user.Password[96:])
	if err_s2 != nil {
		panic(err_s2)
	}

	rn_hash := passHash(password, usr_salt1, usr_salt2)
	auth_ok := pass_cmpP(usr_hash, rn_hash) == 0

	if auth_ok {
		AuditLog.CreateAuditLog(AuditLog.EventAuth, user.ID)
		return http.StatusOK

	} else {
		return http.StatusUnauthorized
	}
}

func CreateUser(phone_mail string, password string, full_name string) int {
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
	AuditLog.CreateAuditLog(AuditLog.EventRegister, user.ID)

	return http.StatusOK
}

func SecUser() {

}

func MakeToken(userID int32) {

}
