package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	sessionCookie   = "session_id"
	sessionLifeTime = 30 * time.Minute
)

/**Логины и пароли*/
var accounts = map[string]string{
	"admin": "admin123",
	"user":  "user123",
}

/**Сессия*/
type Session struct {
	Username  string
	CSRFToken string
	Expires   time.Time
}

/**Все активные сессии*/
var (
	sessions   = map[string]Session{}
	sessionsMu sync.Mutex
)

/** Ресурс*/
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

/**Создание случайной строки*/
func randomString() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

/**Находит сессию текущего пользователя по cookie*/
func getSession(r *http.Request) (Session, bool) {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		return Session{}, false
	}
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	session, ok := sessions[cookie.Value]
	if !ok {
		return Session{}, false
	}
	if time.Now().After(session.Expires) {
		delete(sessions, cookie.Value)
		return Session{}, false
	}
	return session, true
}

/**Пропускает запрос, только если у пользователя есть действующая сессия*/
func protect(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := getSession(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "Войдите в систему")
			return
		}
		if r.Method != http.MethodGet && r.Header.Get("X-CSRF-Token") != session.CSRFToken {
			writeError(w, http.StatusForbidden, "Неверный CSRF-токен")
			return
		}
		next(w, r)
	}
}

/**POST /login - проверка логина и пароля; создание сессии*/
func login(w http.ResponseWriter, r *http.Request) {
	var in struct{ Username, Password string }
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		writeError(w, http.StatusBadRequest, "Неверный JSON")
		return
	}
	password, ok := accounts[in.Username]
	if !ok || password != in.Password {
		writeError(w, http.StatusUnauthorized, "Неверный логин или пароль")
		return
	}
	sessionID := randomString()
	csrfToken := randomString()
	sessionsMu.Lock()
	sessions[sessionID] = Session{
		Username:  in.Username,
		CSRFToken: csrfToken,
		Expires:   time.Now().Add(sessionLifeTime),
	}
	sessionsMu.Unlock()
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(sessionLifeTime.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	writeJSON(w, http.StatusOK, map[string]string{
		"username":   in.Username,
		"csrf_token": csrfToken,
	})
}

/**POST /logout - завершение сессии*/
func logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookie); err == nil {
		sessionsMu.Lock()
		delete(sessions, cookie.Value)
		sessionsMu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{Name: sessionCookie, Value: "", Path: "/", MaxAge: -1})
	w.WriteHeader(http.StatusNoContent)
}

/**GET /session - информация о текущей сессии*/
func sessionInfo(w http.ResponseWriter, r *http.Request) {
	session, _ := getSession(r)
	writeJSON(w, http.StatusOK, map[string]string{
		"username":   session.Username,
		"csrf_token": session.CSRFToken,
		"expires":    session.Expires.Format(time.RFC3339),
	})
}

/** GET /users*/
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

/** GET /users/{id}*/
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

/** POST /users*/
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

/** PUT /users/{id}*/
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

/** DELETE /users/{id}*/
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

/**Маршрутизатор*/
func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", login)
	mux.HandleFunc("POST /logout", protect(logout))
	mux.HandleFunc("GET /sessions", protect(sessionInfo))
	mux.HandleFunc("GET /users", protect(getUsers))
	mux.HandleFunc("GET /users/{id}", protect(getUser))
	mux.HandleFunc("POST /users", protect(createUser))
	mux.HandleFunc("PUT /users/{id}", protect(updateUser))
	mux.HandleFunc("DELETE /users/{id}", protect(deleteUser))
	return mux
}

func main() {
	fmt.Println("Сервер запущен: http://localhost:8080")
	if err := http.ListenAndServe(":8080", newMux()); err != nil {
		fmt.Println("Ошибка:", err)
	}
}
