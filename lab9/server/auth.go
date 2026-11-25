package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
)

/**Логин и пароль*/
const (
	loginUser    = "user"
	passwordUser = "12345"
)

/**Токены*/
var (
	tokens   = map[string]bool{}
	tokensMu sync.Mutex
)

/**Создание случайного токена*/
func newToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

/**POST /login - проверка логина и пароля и выдача токена*/
func login(w http.ResponseWriter, r *http.Request) error {
	var in struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return &AppErr{http.StatusBadRequest, "Неверный JSON"}
	}
	if in.Login != loginUser || in.Password != passwordUser {
		return &AppErr{http.StatusUnauthorized, "Неверный логин или пароль"}
	}
	token := newToken()
	tokensMu.Lock()
	tokens[token] = true
	tokensMu.Unlock()
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
	return nil
}

/**middleware - пропускает запрос дальше, только если в заголовке Authorization верный токен*/
func requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		tokensMu.Lock()
		ok := tokens[token]
		tokensMu.Unlock()
		if !ok {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Требуется авторизация"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

/**Маршруты сервера: /login открыт всем, все остальное только с токеном*/
func newAuthMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", handle(login))
	mux.Handle("/", requireAuth(newMux()))
	return mux
}
