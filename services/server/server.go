package server

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := Response{Message: "Hello, World"}
	json.NewEncoder(w).Encode(response)
}

func RunServer() {
	http.HandleFunc("/", HelloHandler)
	http.ListenAndServe(":8080", nil)
}
