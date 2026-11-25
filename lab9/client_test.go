package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

const testToken = "test-token"

var (
	fakeUsers  map[int]User
	fakeNextID int
	lastAuth   string
)

// Отправляет ответ в формате JSON
func reply(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// Пропускает запрос дальше, только если в заголовке верный токен
func protected(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lastAuth = r.Header.Get("Authorization")
		if lastAuth != "Bearer "+testToken {
			reply(w, 401, map[string]string{"error": "Требуется авторизация"})
			return
		}
		h(w, r)
	}
}

// Проверка данных пользователя (как на настоящем сервере)
func validUser(u User) bool {
	return strings.TrimSpace(u.Name) != "" && u.Age >= 1 && u.Age <= 120
}

// Запускает тестовый сервер и направляет на него клиент
func startServer() *httptest.Server {
	fakeUsers = map[int]User{}
	fakeNextID = 1
	lastAuth = ""
	token = ""

	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		var in struct{ Login, Password string }
		json.NewDecoder(r.Body).Decode(&in)
		if in.Login != "user" || in.Password != "12345" {
			reply(w, 401, map[string]string{"error": "Неверный логин или пароль"})
			return
		}
		reply(w, 200, map[string]string{"token": testToken})
	})

	mux.HandleFunc("GET /users", protected(func(w http.ResponseWriter, r *http.Request) {
		list := []User{}
		for id := 1; id < fakeNextID; id++ {
			if u, ok := fakeUsers[id]; ok {
				list = append(list, u)
			}
		}
		reply(w, 200, list)
	}))

	mux.HandleFunc("GET /users/{id}", protected(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			reply(w, 400, map[string]string{"error": "Неверный id"})
			return
		}
		u, ok := fakeUsers[id]
		if !ok {
			reply(w, 404, map[string]string{"error": "Пользователь не найден"})
			return
		}
		reply(w, 200, u)
	}))

	mux.HandleFunc("POST /users", protected(func(w http.ResponseWriter, r *http.Request) {
		var u User
		if json.NewDecoder(r.Body).Decode(&u) != nil || !validUser(u) {
			reply(w, 400, map[string]string{"error": "Неверные данные"})
			return
		}
		u.ID = fakeNextID
		fakeNextID++
		fakeUsers[u.ID] = u
		reply(w, 201, u)
	}))

	mux.HandleFunc("PUT /users/{id}", protected(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		if _, ok := fakeUsers[id]; !ok {
			reply(w, 404, map[string]string{"error": "Пользователь не найден"})
			return
		}
		var u User
		if json.NewDecoder(r.Body).Decode(&u) != nil || !validUser(u) {
			reply(w, 400, map[string]string{"error": "Неверные данные"})
			return
		}
		u.ID = id
		fakeUsers[id] = u
		reply(w, 200, u)
	}))

	mux.HandleFunc("DELETE /users/{id}", protected(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		if _, ok := fakeUsers[id]; !ok {
			reply(w, 404, map[string]string{"error": "Пользователь не найден"})
			return
		}
		delete(fakeUsers, id)
		w.WriteHeader(204)
	}))

	// служебные маршруты: проверяем, как клиент реагирует на неправильные ответы сервера
	mux.HandleFunc("GET /broken", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "это не json")
	})
	mux.HandleFunc("GET /fail", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "oops", 500)
	})

	ts := httptest.NewServer(mux)
	baseURL = ts.URL
	return ts
}

// Подставляет "ввод пользователя" вместо клавиатуры.
func input(text string) {
	scan = bufio.NewScanner(strings.NewReader(text + strings.Repeat("\n", 5)))
}

// Выполняет функцию и возвращает всё, что она вывела на экран
func capture(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	return string(out)
}

// Добавляет пользователя прямо на тестовый сервер
func addUser(name string, age int) {
	fakeUsers[fakeNextID] = User{ID: fakeNextID, Name: name, Age: age}
	fakeNextID++
}

// Функция request: заголовок с токеном, разные виды ошибок
func TestRequest(t *testing.T) {
	ts := startServer()
	defer ts.Close()

	token = testToken
	var users []User
	if !request("GET", "/users", nil, &users) {
		t.Error("запрос с токеном должен быть успешным")
	}
	if lastAuth != "Bearer "+testToken {
		t.Errorf("заголовок Authorization: %q", lastAuth)
	}

	token = ""
	out := capture(func() {
		if request("GET", "/users", nil, &users) {
			t.Error("запрос без токена должен завершиться ошибкой")
		}
	})
	if lastAuth != "" {
		t.Errorf("без токена заголовок не должен отправляться, получено %q", lastAuth)
	}
	if !strings.Contains(out, "Ошибка 401: Требуется авторизация") {
		t.Errorf("неверное сообщение: %s", out)
	}

	out = capture(func() { request("GET", "/broken", nil, &users) })
	if !strings.Contains(out, "сервер вернул неожиданный ответ") {
		t.Errorf("неверное сообщение: %s", out)
	}

	out = capture(func() { request("GET", "/fail", nil, nil) })
	if !strings.Contains(out, "Ошибка 500: Internal Server Error") {
		t.Errorf("неверное сообщение: %s", out)
	}

	ts.Close()
	out = capture(func() {
		if request("GET", "/users", nil, nil) {
			t.Error("запрос к недоступному серверу должен завершиться ошибкой")
		}
	})
	if !strings.Contains(out, "не удалось связаться с сервером") {
		t.Errorf("неверное сообщение: %s", out)
	}
}

// POST /login, вход и выход
func TestSignIn(t *testing.T) {
	ts := startServer()
	defer ts.Close()

	input("user\nwrong\n")
	out := capture(signIn)
	if token != "" {
		t.Errorf("токен не должен сохраняться, получено %q", token)
	}
	if !strings.Contains(out, "Ошибка 401: Неверный логин или пароль") {
		t.Errorf("неверное сообщение: %s", out)
	}

	input("user\n12345\n")
	out = capture(signIn)
	if token != testToken {
		t.Errorf("токен не сохранён, получено %q", token)
	}
	if !strings.Contains(out, "Вход выполнен") {
		t.Errorf("неверное сообщение: %s", out)
	}

	capture(signOut)
	if token != "" {
		t.Errorf("после выхода токен должен быть пустым, получено %q", token)
	}
}

// GET /users
func TestListUsers(t *testing.T) {
	ts := startServer()
	defer ts.Close()
	token = testToken

	out := capture(listUsers)
	if !strings.Contains(out, "Пользователей нет") {
		t.Errorf("неверный вывод: %s", out)
	}

	addUser("Anna", 20)
	addUser("Anton", 25)
	out = capture(listUsers)
	for _, want := range []string{"Anna", "Anton", "Всего: 2"} {
		if !strings.Contains(out, want) {
			t.Errorf("в выводе нет %q: %s", want, out)
		}
	}
}

// GET /users/{id}
func TestShowUser(t *testing.T) {
	ts := startServer()
	defer ts.Close()
	token = testToken
	addUser("Anna", 20)

	input("1\n")
	out := capture(showUser)
	if !strings.Contains(out, "Anna") || !strings.Contains(out, "20") {
		t.Errorf("неверный вывод: %s", out)
	}

	input("999\n")
	out = capture(showUser)
	if !strings.Contains(out, "Ошибка 404: Пользователь не найден") {
		t.Errorf("неверный вывод: %s", out)
	}

	input("abc\n")
	out = capture(showUser)
	if !strings.Contains(out, "Ошибка 400: Неверный id") {
		t.Errorf("неверный вывод: %s", out)
	}
}

// POST /users
func TestCreateUser(t *testing.T) {
	ts := startServer()
	defer ts.Close()
	token = testToken

	input("Anna\n20\n")
	out := capture(createUser)
	if len(fakeUsers) != 1 || fakeUsers[1].Name != "Anna" || fakeUsers[1].Age != 20 {
		t.Errorf("пользователь не добавлен: %v", fakeUsers)
	}
	if !strings.Contains(out, "Пользователь добавлен") {
		t.Errorf("неверный вывод: %s", out)
	}

	input("Boris\nabc\n")
	out = capture(createUser)
	if len(fakeUsers) != 1 {
		t.Error("пользователь не должен добавляться")
	}
	if !strings.Contains(out, "Возраст должен быть числом") {
		t.Errorf("неверный вывод: %s", out)
	}

	for _, text := range []string{"Boris\n200\n", "\n20\n"} {
		input(text)
		out = capture(createUser)
		if len(fakeUsers) != 1 {
			t.Errorf("ввод %q: пользователь не должен добавляться", text)
		}
		if !strings.Contains(out, "Ошибка 400") {
			t.Errorf("ввод %q: неверный вывод: %s", text, out)
		}
	}
}

// PUT /users/{id}
func TestUpdateUser(t *testing.T) {
	ts := startServer()
	defer ts.Close()
	token = testToken
	addUser("Anna", 20)

	input("1\n\n21\n")
	out := capture(updateUser)
	if u := fakeUsers[1]; u.Name != "Anna" || u.Age != 21 {
		t.Errorf("данные не обновились: %v", u)
	}
	if !strings.Contains(out, "Пользователь обновлен") {
		t.Errorf("неверный вывод: %s", out)
	}

	input("1\nAnastasia\n\n")
	capture(updateUser)
	if u := fakeUsers[1]; u.Name != "Anastasia" || u.Age != 21 {
		t.Errorf("данные не обновились: %v", u)
	}

	input("999\n")
	out = capture(updateUser)
	if !strings.Contains(out, "Ошибка 404") {
		t.Errorf("неверный вывод: %s", out)
	}

	input("1\n\nabc\n")
	out = capture(updateUser)
	if fakeUsers[1].Age != 21 {
		t.Error("возраст не должен меняться")
	}
	if !strings.Contains(out, "Возраст должен быть числом") {
		t.Errorf("неверный вывод: %s", out)
	}
}

// DELETE /users/{id}
func TestDeleteUser(t *testing.T) {
	ts := startServer()
	defer ts.Close()
	addUser("Anna", 20)

	input("1\n")
	capture(deleteUser)
	if len(fakeUsers) != 1 {
		t.Error("без входа пользователь не должен удаляться")
	}
	token = testToken

	input("999\n")
	out := capture(deleteUser)
	if !strings.Contains(out, "Ошибка 404") {
		t.Errorf("неверный вывод: %s", out)
	}
	if len(fakeUsers) != 1 {
		t.Error("существующий пользователь не должен удаляться")
	}

	input("abc\n")
	out = capture(deleteUser)
	if !strings.Contains(out, "Ошибка") {
		t.Errorf("ожидалась ошибка, получено: %s", out)
	}

	input("1\n")
	out = capture(deleteUser)
	if len(fakeUsers) != 0 {
		t.Error("пользователь не удалён")
	}
	if !strings.Contains(out, "Пользователь удален") {
		t.Errorf("неверный вывод: %s", out)
	}
}
