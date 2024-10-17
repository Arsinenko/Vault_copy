package server

import (
	serviceApp "Vault_copy/services/app"
	serviceUser "Vault_copy/services/user"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/pgtype"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
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

	status := serviceUser.AuthStandard(authReq.PhoneMail, authReq.Password)

	response := Response{
		Status: status,
	}

	if status == http.StatusOK {
		response.Message = "Authentication successful"
	} else {
		response.Message = "Authentication failed"
	}

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

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func CreateAppHandler(w http.ResponseWriter, r *http.Request) {
	var app struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		OwnerID     int32  `json:"owner_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := serviceApp.CreateApp(app.Name, app.Description, app.OwnerID, pgtype.JSONB{})
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "App creation attempt", Status: status}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func ChangeAppNameHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int32  `json:"user_id"`
		AppID  int32  `json:"app_id"`
		Name   string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := serviceApp.API_AppChangeName(req.UserID, req.AppID, req.Name)
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "App name change attempt", Status: status}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func ChangeAppDescriptionHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID      int32  `json:"user_id"`
		AppID       int32  `json:"app_id"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := serviceApp.API_AppChangeDescription(req.UserID, req.AppID, req.Description)
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "App description change attempt", Status: status}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func CreateSecretHandler(w http.ResponseWriter, r *http.Request) {
	var secret struct {
		SID      string       `json:"sid"`
		Data     []byte       `json:"data"`
		AppID    int32        `json:"app_id"`
		Metadata pgtype.JSONB `json:"metadata"`
	}
	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := serviceUser.CreateSecret(secret.SID, secret.Data, secret.AppID, secret.Metadata)
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "Secret creation attempt", Status: status}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func RunServer() {
	r := mux.NewRouter()
	r.HandleFunc("/api/user/auth", AuthHandler).Methods("POST")
	r.HandleFunc("/api/user/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/app/create", CreateAppHandler).Methods("POST")
	r.HandleFunc("/api/app/{app_id}/name", ChangeAppNameHandler).Methods("PUT")
	r.HandleFunc("/api/app/{app_id}/description", ChangeAppDescriptionHandler).Methods("PUT")
	r.HandleFunc("/api/app/{app_id}/secret", CreateSecretHandler).Methods("POST")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
