package internal_api

import (
	"Vault_copy/db_operations"
	"Vault_copy/db_operations/cryptoOperation"
	"Vault_copy/db_operations/models"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	LogService "Vault_copy/services/log"

	"github.com/jackc/pgx/pgtype"
	"github.com/jinzhu/gorm"
)

const (
	I_rule_view_app_info   = "view_app_info"
	I_rule_change_app_name = "change_app_name"
	I_rule_change_app_desc = "change_app_desc"
)

var defaultRules = map[string]bool{
	I_rule_view_app_info:   true,
	I_rule_change_app_name: false,
	I_rule_change_app_desc: false,
}

func I_policy_agreement(UserID int32, AppID int32) (bool, error) {
	//add a record app with default rules
	//policy, err := I_add_user_policy(UserID, AppID)
	//if err != nil {
	//	return false, err
	//}
	//return policy != nil, nil
	// Check if a user policy exists
	db, _e := db_operations.InitDB()
	if _e != nil {
		return false, _e
	}

	if !I_policy_exists(db, UserID) {
		return false, nil
	}

	// Check if the user has accepted the app rules
	jsonRules, err := I_dec_policy(UserID, AppID)
	if err != nil {
		return false, err
	}
	if jsonRules == nil {
		return false, nil
	}

	// Check if all rules are accepted
	for rule, accepted := range *jsonRules {
		if !accepted {
			return false, nil
		}

		// If the rule doesn't exist in the default rules, it's an error
		if _, ok := defaultRules[rule]; !ok {
			return false, fmt.Errorf("unknown rule %s", rule)
		}
	}

	return true, nil
}

// I_policy_exists checks if a user policy exists in the database.
func I_policy_exists(db *gorm.DB, UserID int32) bool {
	var user models.User
	return db.First(&user, UserID).Error == nil
}

// I_dec_policy decodes user rules from the database into a map.
func I_dec_policy(UserID int32, AppID int32) (*map[string]bool, error) {
	logHash := generateLogHash(UserID, AppID)

	db, err := db_operations.InitDB()
	if err != nil {
		logError(LogService.ErrorDBInit, "[I_dec_policy]::db_operations.InitDB()", logHash)
		return nil, err
	}

	var policy models.Policy
	if err := db.First(policy, "user_id = ? AND app_id = ?", UserID, AppID).Error; err != nil {
		logError(LogService.ErrorDBExec, "[I_dec_policy]::db.First()", logHash)
		return nil, err
	}
	if policy.AppID < 0 {
		return nil, nil
	}

	// Decode JSON rules
	var jsonRules map[string]bool
	if err := json.Unmarshal(policy.Rules.Bytes, &jsonRules); err != nil {
		logError(LogService.ErrorJSONUnmarshal, "[I_dec_policy]::json.Unmarshal()", logHash)
		return nil, err
	}

	return &jsonRules, nil
}

// I_enc_policy encodes and stores user rules in the database.
func I_enc_policy(UserID int32, AppID int32, jsonRules map[string]bool) (*models.Policy, error) {
	logHash := generateLogHash(UserID, AppID, jsonRules)

	db, err := db_operations.InitDB()
	if err != nil {
		logError(LogService.ErrorDBInit, "[I_enc_policy]::db_operations.InitDB()", logHash)
		return nil, err
	}

	if !I_policy_exists(db, UserID) {
		return nil, fmt.Errorf("user does not exist")
	}

	var policy *models.Policy
	if err := db.First(&policy, "user_id = ? AND app_id = ?", UserID, AppID).Error; err != nil {
		logError(LogService.ErrorDBExec, "[I_enc_policy]::db.First()", logHash)
		return nil, err
	}
	if policy == nil {
		return nil, nil
	}

	// Encode JSON to bytes
	rawStr, err := json.Marshal(jsonRules)
	if err != nil {
		logError(LogService.ErrorJSONMarshal, "[I_enc_policy]::json.Marshal()", logHash)
		return nil, err
	}

	// Update and save policy entry
	policy.Rules.Bytes = rawStr
	policy.DateChanged.Time = time.Now()

	if err := db.Save(&policy).Error; err != nil {
		logError(LogService.ErrorDBSave, "[I_enc_policy]::db.Save()", logHash)
		return nil, err
	}

	return policy, nil
}

// I_set_policy_rule sets a specific policy rule for a user.
func I_set_policy_rule(UserID int32, AppID int32, Rule string, value bool) (*models.Policy, error) {
	logHash := generateLogHash(UserID, AppID, Rule, value)

	jsonRules, err := I_dec_policy(UserID, AppID)
	if err != nil {
		logError(LogService.ErrorIDecPolicy, "[I_set_policy_rule]::I_dec_policy()", logHash)
		return nil, err
	}
	if jsonRules == nil {
		return nil, nil
	}

	// Set the rule value
	(*jsonRules)[Rule] = value

	// Save updated policy
	return I_enc_policy(UserID, AppID, *jsonRules)
}

// I_get_policy_rule retrieves the value of a specific policy rule.
func I_get_policy_rule(UserID int32, AppID int32, Rule string) (bool, error) {
	logHash := generateLogHash(UserID, AppID)

	jsonRules, err := I_dec_policy(UserID, AppID)
	if err != nil {
		logError(LogService.ErrorIDecPolicy, "[I_get_policy_rule]::I_dec_policy()", logHash)
		return false, err
	}
	if jsonRules == nil {
		return false, nil
	}

	return (*jsonRules)[Rule], nil
}

// I_add_user_policy adds a new user policy entry with default values.
func I_add_user_policy(UserID int32, AppID int32) (*models.Policy, error) {
	logHash := generateLogHash(UserID, AppID)

	db, err := db_operations.InitDB()
	if err != nil {
		logError(LogService.ErrorDBInit, "[I_add_user_policy]::db_operations.InitDB()", logHash)
		return nil, err
	}

	if !I_policy_exists(db, UserID) {
		return nil, fmt.Errorf("user does not exist")
	}

	// Create a new policy with default values
	rawStr, err := json.Marshal(defaultRules)
	if err != nil {
		logError(LogService.ErrorJSONMarshal, "[I_add_user_policy]::json.Marshal()", logHash)
		return nil, err
	}

	policy := &models.Policy{
		UserID:      UserID,
		AppID:       AppID,
		Rules:       pgtype.JSONB{Bytes: rawStr},
		DateCreated: pgtype.Timestamptz{Time: time.Now()},
		DateChanged: pgtype.Timestamptz{Time: time.Now()},
	}

	if err := db.Create(policy).Error; err != nil {
		logError(LogService.ErrorDBSave, "[I_add_user_policy]::db.Create()", logHash)
		return nil, err
	}

	return policy, nil
}

// I_del_user_policy deletes a user policy entry.
func I_del_user_policy(UserID int32, AppID int32) (*models.Policy, error) {
	logHash := generateLogHash(UserID, AppID)

	db, err := db_operations.InitDB()
	if err != nil {
		logError(LogService.ErrorDBInit, "[I_del_user_policy]::db_operations.InitDB()", logHash)
		return nil, err
	}

	var policy models.Policy
	if err := db.First(&policy, "user_id = ? AND app_id = ?", UserID, AppID).Error; err != nil {
		logError(LogService.ErrorDBExec, "[I_del_user_policy]::db.First()", logHash)
		return nil, err
	}

	if err := db.Delete(&policy).Error; err != nil {
		logError(LogService.ErrorDBDelete, "[I_del_user_policy]::db.Delete()", logHash)
		return nil, err
	}

	return &policy, nil
}

// I_get_policy_all retrieves all user policy rules as JSON.
func I_get_policy_all(UserID int32, AppID int32) (*map[string]bool, error) {
	logHash := generateLogHash(UserID, AppID)

	jsonRules, err := I_dec_policy(UserID, AppID)
	if err != nil {
		logError(LogService.ErrorIDecPolicy, "[I_get_policy_all]::I_dec_policy()", logHash)
		return nil, err
	}
	return jsonRules, nil
}

// I_set_policy_all overrides the given policy rules.
func I_set_policy_all(UserID int32, AppID int32, jsonStr string) (*models.Policy, error) {
	logHash := generateLogHash(UserID, AppID, jsonStr)

	var jsonRules map[string]bool
	if err := json.Unmarshal([]byte(jsonStr), &jsonRules); err != nil {
		logError(LogService.ErrorJSONUnmarshal, "[I_set_policy_all]::json.Unmarshal()", logHash)
		return nil, err
	}

	return I_enc_policy(UserID, AppID, jsonRules)
}

// I_get_policy_time_changed retrieves the last changed time of a user policy.
func I_get_policy_time_changed(UserID int32, AppID int32) (time.Time, error) {
	logHash := generateLogHash(UserID, AppID)

	db, err := db_operations.InitDB()
	if err != nil {
		logError(LogService.ErrorDBInit, "[I_get_policy_time_changed]::db_operations.InitDB()", logHash)
		return time.Time{}, err
	}

	var policy models.Policy
	if err := db.First(&policy, "user_id = ? AND app_id = ?", UserID, AppID).Error; err != nil {
		logError(LogService.ErrorDBExec, "[I_get_policy_time_changed]::db.First()", logHash)
		return time.Time{}, err
	}

	return policy.DateChanged.Time, nil
}

// I_get_policy_time_created retrieves the create timestamp of a user policy.
func I_get_policy_time_created(UserID int32, AppID int32) (time.Time, error) {
	logHash := generateLogHash(UserID, AppID)

	db, err := db_operations.InitDB()
	if err != nil {
		logError(LogService.ErrorDBInit, "[I_get_policy_time_created]::db_operations.InitDB()", logHash)
		return time.Time{}, err
	}

	var policy models.Policy
	if err := db.First(&policy, "user_id = ? AND app_id = ?", UserID, AppID).Error; err != nil {
		logError(LogService.ErrorDBExec, "[I_get_policy_time_created]::db.First()", logHash)
		return time.Time{}, err
	}

	return policy.DateCreated.Time, nil
}

// Helper functions

// generateLogHash creates a log hash from the provided inputs.
func generateLogHash(inputs ...interface{}) string {
	hashInput := ""
	for _, input := range inputs {
		hashInput += fmt.Sprint(input)
	}
	return hex.EncodeToString(cryptoOperation.SHA256([]byte(hashInput)))
}

// logError logs an error message with the provided parameters.
func logError(errorType int16, message string, logHash string) {
	LogService.Push_server_log(errorType, LogService.TErrorDBInit, message, logHash)
}
