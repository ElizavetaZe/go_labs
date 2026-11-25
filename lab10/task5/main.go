package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Секретный ключ
var secret = []byte("lab10-secret-key")

const tokenLifetime = time.Hour

// Учётная запись
type Account struct {
	Password string
	Role     string
}

// Список учётных записей
var accounts = map[string]Account{
	"admin": {"admin123", "admin"},
	"user":  {"user123", "user"},
}

// Данные токена
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Ресурс
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Пользователи хранятся в памяти
var (
	users  = map[int]User{1: {1, "Anna", 20}, 2: {2, "Anton", 25}}
	nextID = 3
	mu     sync.Mutex
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func validUser(u User) bool {
	return strings.TrimSpace(u.Name) != "" && u.Age >= 1 && u.Age <= 120
}

// Проверяет токен из заголовка
func parseToken(header string) (*Claims, error) {
	tokenString := strings.TrimPrefix(header, "Bearer ")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !token.Valid {
		return nil, errors.New("неверный токен")
	}
	return claims, nil
}

// Middleware: пропускает запрос дальше, только если токен верный, а роль пользователя входит в список roles
func protect(roles []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := parseToken(r.Header.Get("Authorization"))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Нужен действительный токен")
			return
		}
		if !slices.Contains(roles, claims.Role) {
			writeError(w, http.StatusForbidden, "Недостаточно прав")
			return
		}
		next(w, r)
	}
}

// POST /login - проверяет логин и пароль и выдаёт токен
func login(w http.ResponseWriter, r *http.Request) {
	var in struct{ Username, Password string }
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		writeError(w, http.StatusBadRequest, "Неверный JSON")
		return
	}
	account, ok := accounts[in.Username]
	if !ok || account.Password != in.Password {
		writeError(w, http.StatusUnauthorized, "Неверный логин или пароль")
		return
	}

	claims := Claims{
		Username: in.Username,
		Role:     account.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось создать токен")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

// GET /users
func getUsers(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	list := []User{}
	for id := 1; id < nextID; id++ {
		if u, ok := users[id]; ok {
			list = append(list, u)
		}
	}
	mu.Unlock()
	writeJSON(w, http.StatusOK, list)
}

// GET /users/{id}
func getUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	mu.Lock()
	u, ok := users[id]
	mu.Unlock()
	if !ok {
		writeError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}
	writeJSON(w, http.StatusOK, u)
}

// POST /users
func createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if json.NewDecoder(r.Body).Decode(&u) != nil || !validUser(u) {
		writeError(w, http.StatusBadRequest, "Неверные данные")
		return
	}
	mu.Lock()
	u.ID = nextID
	nextID++
	users[u.ID] = u
	mu.Unlock()
	writeJSON(w, http.StatusCreated, u)
}

// PUT /users/{id}
func updateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	var u User
	if json.NewDecoder(r.Body).Decode(&u) != nil || !validUser(u) {
		writeError(w, http.StatusBadRequest, "Неверные данные")
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := users[id]; !ok {
		writeError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}
	u.ID = id
	users[id] = u
	writeJSON(w, http.StatusOK, u)
}

// DELETE /users/{id}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	mu.Lock()
	_, ok := users[id]
	delete(users, id)
	mu.Unlock()
	if !ok {
		writeError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Маршрутизатор
func newMux() *http.ServeMux {
	anyRole := []string{"user", "admin"}
	adminOnly := []string{"admin"}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", login)
	mux.HandleFunc("GET /users", protect(anyRole, getUsers))
	mux.HandleFunc("GET /users/{id}", protect(anyRole, getUser))
	mux.HandleFunc("POST /users", protect(adminOnly, createUser))
	mux.HandleFunc("PUT /users/{id}", protect(adminOnly, updateUser))
	mux.HandleFunc("DELETE /users/{id}", protect(adminOnly, deleteUser))
	return mux
}

func main() {
	fmt.Println("Сервер запущен: http://localhost:8080")
	http.ListenAndServe(":8080", newMux())
}
