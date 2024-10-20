package server

import (
	serviceApp "Vault_copy/services/app"
	serviceSecret "Vault_copy/services/secret"
	serviceUser "Vault_copy/services/user"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/pgtype"
)

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data"`
}

type AuthRequest struct {
	PhoneMail string `json:"phone_mail"`
	Password  string `json:"password"`
}

type RegisterRequest struct {
	PhoneMail string `json:"phone_mail"`
	Password  string `json:"password"`
	FullName  string `json:"full_name"`
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authReq AuthRequest
	err := json.NewDecoder(r.Body).Decode(&authReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, token := serviceUser.AuthStandard(authReq.PhoneMail, authReq.Password)

	response := Response{
		Status: status,
		Data: token,
	}

	if status == http.StatusOK {
		response.Message = "Authentication successful"

		// Устанавливаем cookie с токеном
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	} else {
		response.Message = "Authentication failed"
	}

	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var regReq RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&regReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := serviceUser.Register(regReq.PhoneMail, regReq.Password, regReq.FullName)

	response := Response{
		Status: status,
	}

	if status == http.StatusOK {
		response.Message = "Registration successful"
	} else {
		response.Message = "Registration failed"
	}

	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, err := r.Cookie("token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uid, err := serviceUser.AuthWithToken(token.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	status := serviceUser.DeleteUser(uid, 0)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "User deletion attempt", Status: status}
	_err := json.NewEncoder(w).Encode(response)
	if _err != nil {
		return
	}
}

func CreateAppHandler(w http.ResponseWriter, r *http.Request) {
	var app struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := r.Cookie("token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uid, err := serviceUser.AuthWithToken(token.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	status := serviceApp.CreateApp(app.Name, app.Description, uid, pgtype.JSONB{})

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "App creation attempt", Status: status}
	err_ := json.NewEncoder(w).Encode(response)
	if err_ != nil {
		return
	}
}

func ChangeAppNameHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name   string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	AppIDstr := vars["app_id"]
	AppID, er := strconv.Atoi(AppIDstr)
	if er != nil {
		http.Error(w, er.Error(), http.StatusBadRequest)
		return
	}

	token, err := r.Cookie("token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uid, err := serviceUser.AuthWithToken(token.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	status := serviceApp.API_AppChangeName(uid, int32(AppID), req.Name)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "App name change attempt", Status: status}
	err_ := json.NewEncoder(w).Encode(response)
	if err_ != nil {
		return
	}
}

func ChangeAppDescriptionHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	AppIDstr := vars["app_id"]
	AppID, er := strconv.Atoi(AppIDstr)
	if er != nil {
		http.Error(w, er.Error(), http.StatusBadRequest)
		return
	}

	token, err := r.Cookie("token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uid, err := serviceUser.AuthWithToken(token.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	status := serviceApp.API_AppChangeDescription(uid, int32(AppID), req.Description)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "App description change attempt", Status: status}
	err_ := json.NewEncoder(w).Encode(response)
	if err_ != nil {
		return
	}
}

func CreateSecretHandler(w http.ResponseWriter, r *http.Request) {
	var secret struct {
		SID   string `json:"sid"`
		Data  []byte `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	AppIDstr := vars["app_id"]
	AppID, er := strconv.Atoi(AppIDstr)
	if er != nil {
		http.Error(w, er.Error(), http.StatusBadRequest)
		return
	}

	token, err := r.Cookie("token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uid, auth_err := serviceUser.AuthWithToken(token.Value)
	if auth_err != nil {
		http.Error(w, auth_err.Error(), http.StatusUnauthorized)
		return
	}
	if uid < 0 {
		http.Error(w, auth_err.Error(), http.StatusUnauthorized)
		return
	}

	status := serviceSecret.CreateSecret(secret.Data, int32(AppID), "{}")

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "Secret creation attempt", Status: status}
	err_ := json.NewEncoder(w).Encode(response)
	if err_ != nil {
		return
	}
}

func HTTP_app_get_name(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uid, err := serviceUser.AuthWithToken(token.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	AppIDstr := vars["app_id"]
	AppID, er := strconv.Atoi(AppIDstr)
	if er != nil {
		// Handle the error (e.g., invalid app_id)
		http.Error(w, "Invalid app_id", http.StatusBadRequest)
		return
	}

	res_name, res_status := serviceApp.API_AppGetName(uid, int32(AppID))

	w.WriteHeader(res_status)
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "App get name attempt", Status: res_status, Data: res_name}
	err_ := json.NewEncoder(w).Encode(response)
	if err_ != nil {
		return
	}
}

func GetSecretsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, token_err := r.Cookie("token")
	if token_err != nil {
		http.Error(w, token_err.Error(), http.StatusBadRequest)
		return
	}

	uid, err := serviceUser.AuthWithToken(token.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if uid < 0 {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	AppIDstr := vars["app_id"]
	AppID, er := strconv.Atoi(AppIDstr)
	if er != nil {
		http.Error(w, er.Error(), http.StatusBadRequest)
		return
	}

	secrets, status := serviceSecret.GetSecrets(int32(AppID))

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "Secrets get attempt", Status: status, Data: secrets}
	err_ := json.NewEncoder(w).Encode(response)
	if err_ != nil {
		return
	}
}

// TODO test
func DeleteSecretHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		ID    int64 `json:"id"`
		AppID int32 `json:"app_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := serviceSecret.DeleteSecret(req.ID, req.AppID)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "Secret deletion attempt", Status: status}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}

}

func RunServer() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/user/auth", AuthHandler).Methods("POST")
	r.HandleFunc("/api/v1/user/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/v1/user/delete", DeleteUserHandler).Methods("POST") //TODO test it
	r.HandleFunc("/api/v1/app/create", CreateAppHandler).Methods("POST")

	// [GET] /api/v1/app/{app_id}/[date_update, date_create]]
	r.HandleFunc("/api/v1/app/{app_id}/name", ChangeAppNameHandler).Methods("PUT")
	r.HandleFunc("/api/v1/app/{app_id}/name", HTTP_app_get_name).Methods("GET")
	r.HandleFunc("/api/v1/app/{app_id}/description", ChangeAppDescriptionHandler).Methods("PUT") // +++ [GET] /api/app/{app_id}/description -- return name of app
	r.HandleFunc("/api/v1/app/{app_id}/secret", CreateSecretHandler).Methods("POST")
	r.HandleFunc("/api/v1/app/{app_id}/secrets", GetSecretsHandler).Methods("GET")

	err := http.ListenAndServe(":" + os.Getenv("APP_PORT"), r)
	if err != nil {
		return
	}
}
