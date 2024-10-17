package server

import (
	service_user "Vault_copy/services/user"
	"encoding/json"
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

	status := service_user.AuthStandard(authReq.PhoneMail, authReq.Password)

	response := Response{
		Status: status,
	}

	if status == http.StatusOK {
		response.Message = "Authentication successful"
	} else {
		response.Message = "Authentication failed"
	}

	json.NewEncoder(w).Encode(response)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var regReq RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&regReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := service_user.CreateUser(regReq.PhoneMail, regReq.Password, regReq.FullName)

	response := Response{
		Status: status,
	}

	if status == http.StatusOK {
		response.Message = "Registration successful"
	} else {
		response.Message = "Registration failed"
	}

	json.NewEncoder(w).Encode(response)
}

func RunServer() {
	http.HandleFunc("/api/user/auth", AuthHandler)
	http.HandleFunc("/api/user/register", RegisterHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
