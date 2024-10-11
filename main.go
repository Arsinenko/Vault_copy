package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Структура пользователя
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// "База данных" пользователей в виде карты
var users = make(map[int]User)
var mu sync.Mutex // для потокобезопасной работы с картой

// Функция для добавления нового пользователя
func addUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Если пользователь уже существует
	if _, exists := users[user.ID]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	users[user.ID] = user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Функция для изменения существующего пользователя
func updateUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Получение ID пользователя из URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Если пользователя нет в базе
	if _, exists := users[id]; !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Обновляем информацию о пользователе
	users[id] = user
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Функция для получения списка всех пользователей
func getUsers(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func main() {
	// Маршруты для нашего API
	http.HandleFunc("/users", getUsers)         // Получение всех пользователей
	http.HandleFunc("/user/add", addUser)       // Добавление пользователя
	http.HandleFunc("/user/update", updateUser) // Обновление пользователя

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
