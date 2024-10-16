package internal_api

import (
	"Vault_copy/db_operations"
	"Vault_copy/db_operations/cryptoOperation"
	"Vault_copy/db_operations/models"
	"fmt"

	LogService "Vault_copy/services/log"
	"encoding/hex"
	"encoding/json"
)

const (
	I_rule_change_app_name = "change_app_name"
	I_rule_change_app_desc = "change_app_desc"
)

// * FINAL - STATIC API - L3 INTERNAL API
// Decode user rules as map[string] bool
func I_dec_policy(UserID int32, AppID int32) (*map[string] bool, error) {
	_log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(fmt.Sprint(UserID) + fmt.Sprint(AppID))));

	// init database
	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[I_dec_policy]::db_operations.InitDB()", _log_hash)
		return nil, e
	}

	// try to find entry
	var policy *models.Policy
	res := db.First(&policy, "user_id = ? AND app_id = ?", UserID, AppID)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[I_dec_policy]::db.First()", _log_hash)
		return nil, res.Error
	}

	// user isn't in app
	if (policy == nil) { return nil, nil }

	// decode json
	var json_o map[string] bool;
	if err := json.Unmarshal(policy.Rules.Bytes, &json_o); err != nil {
		LogService.Push_server_log(LogService.ErrorJSONUnmarshal, LogService.TErrorJSONUnmarshal, "[I_dec_policy]::json.Unmarshal()", _log_hash)
    return nil, err;
	}

	return &json_o, nil;
}

// * FINAL - STATIC API - L3 INTERNAL API
// Encode input user rules map[string] bool and store it in database
func I_enc_policy(UserID int32, AppID int32, json_o map[string] bool) (*models.Policy, error) {
	_log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(fmt.Sprint(UserID) + fmt.Sprint(AppID) + fmt.Sprint(json_o))));

	// init database
	db, e := db_operations.InitDB()
	if e != nil {
		LogService.Push_server_log(LogService.ErrorDBInit, LogService.TErrorDBInit, "[I_enc_policy]::db_operations.InitDB()", _log_hash)
		return nil, e
	}

	// try to find entry
	var policy *models.Policy
	res := db.First(policy, "user_id = ? AND app_id = ?", UserID, AppID)
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBExec, LogService.TErrorDBExec, "[I_enc_policy]::db.First()", _log_hash)
		return nil, res.Error
	}
	// user isn't in app
	if (policy == nil) { return nil, nil }

	// encode json to bytes
	raw_str, err := json.Marshal(json_o);
	if err != nil {
		LogService.Push_server_log(LogService.ErrorJSONUnmarshal, LogService.TErrorJSONUnmarshal, "[I_enc_policy]::json.Unmarshal()", _log_hash)
    return nil, err;
	}

	// update and save policy entry
	policy.Rules.Bytes = raw_str;
	res = db.Save(&policy);
	if res.Error != nil {
		LogService.Push_server_log(LogService.ErrorDBSave, LogService.TErrorDBSave, "[I_enc_policy]::db.Save()", _log_hash)
		return nil, res.Error
	}

	return policy, nil;
}

// * FINAL - STATIC API - L3 INTERNAL API
// Set user policy rule if user exists in app
func I_set_policy_rule(UserID int32, AppID int32, Rule string, value bool) (*models.Policy, error) {
	var lv string;
	if (value) { lv = "1" } else { lv = "0"}
	_log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(fmt.Sprint(UserID) + fmt.Sprint(AppID) + Rule + string(lv))));

	// decode existent policy entry from database
	json_o, err := I_dec_policy(UserID, AppID);
	if (err != nil) {
		LogService.Push_server_log(LogService.ErrorIDecPolicy, LogService.TErrorIDecPolicy, "[I_set_policy_rule]::I_dec_policy()", _log_hash)
		return nil, err
	}
	if (json_o == nil) { return nil, nil }

	// Set value at rule
	(*json_o)[Rule] = value;

	// save changed policy entry into database
	policy, err := I_enc_policy(UserID, AppID, *json_o);
	if (err != nil) { 
		LogService.Push_server_log(LogService.ErrorIEncPolicy, LogService.TErrorIEncPolicy, "[I_set_policy_rule]::I_dec_policy()", _log_hash)
		return nil, err
	}
	if (policy == nil) { return nil, nil }

	return policy, nil;
}

// * FINAL - STATIC API - L3 INTERNAL API
// Get user policy rule value
func I_get_policy_rule(UserID int32, AppID int32, Rule string) (bool, error) {
	_log_hash := hex.EncodeToString(cryptoOperation.SHA256([]byte(string(UserID) + string(AppID))));

	// decode existent policy entry from database
	json_o, err := I_dec_policy(UserID, AppID);
	if (err != nil) {
		LogService.Push_server_log(LogService.ErrorIDecPolicy, LogService.TErrorIDecPolicy, "[I_get_policy_rule]::I_dec_policy()", _log_hash)
		return false, err
	}
	if (json_o == nil) { return false, nil } // false - user not in app
	
	// false - policy critical error | impossible state
	return (*json_o)[Rule] || false, nil;
}