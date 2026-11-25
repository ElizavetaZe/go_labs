package main

import (
	"database/sql"
	"net/http/httptest"
	"strings"
	"testing"
)

/**Тестовая база*/
var testDB = strings.Replace(baza, "/postgres?", "/lab8_test?", 1)

/**Подключение к тестовой базе, очищение при каждом тесте*/
func prepare(t *testing.T) {
	var err error
	db, err = sql.Open("postgres", testDB)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id   SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		age  INT NOT NULL
	)`)
	if err != nil {
		t.Fatal("Не удалось подключиться к тестовой базе: ", err)
	}
	_, err = db.Exec("TRUNCATE users RESTART IDENTITY")
	if err != nil {
		t.Fatal(err)
	}
}

/**Отправка запроса в сервер (без запуска настоящего)*/
func call(method, url, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, strings.NewReader(body))
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, req)
	return rec
}

/**POST /users*/
func TestCreateUser(t *testing.T) {
	prepare(t)
	defer db.Close()

	rec := call("POST", "/users", `{"name":"Анна","age":20}`)
	if rec.Code != 201 {
		t.Errorf("статус %d, ожидался 201", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Анна") {
		t.Errorf("В ответе нет имени: %s", rec.Body.String())
	}

	for _, body := range []string{`{"name":"Борис","age":1}`, `{"name":"Вера","age":120}`} {
		rec := call("POST", "/users", body)
		if rec.Code != 201 {
			t.Errorf("тело %q: статус %d, ожидался 201", body, rec.Code)
		}
	}

	badBodies := []string{
		`это не json`,
		``,
		`{"name":"","age":20}`,
		`{"name":"   ","age":20}`,
		`{"age":20}`,
		`{"name":"Анна"}`,
		`{"name":"Анна","age":0}`,
		`{"name":"Анна","age":121}`,
		`{"name":"Анна","age":-5}`,
		`{"name":"Анна","age":"двадцать"}`,
	}
	for _, body := range badBodies {
		rec := call("POST", "/users", body)
		if rec.Code != 400 {
			t.Errorf("тело %q: статус %d, ожидался 400", body, rec.Code)
		}
	}
}

/**GET /users*/
func TestGetUsers(t *testing.T) {
	prepare(t)
	defer db.Close()

	rec := call("GET", "/users", "")
	if rec.Code != 200 {
		t.Errorf("статус %d, ожидался 200", rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Errorf("ожидался [], получено %s", rec.Body.String())
	}
	call("POST", "/users", `{"name":"Anna","age":20}`)
	call("POST", "/users", `{"name":"Anton","age":25}`)
	call("POST", "/users", `{"name":"Boris","age":25}`)

	cases := map[string]int{
		"/users":                3,
		"/users?name=An":        2,
		"/users?age=25":         2,
		"/users?name=An&age=25": 1,
		"/users?limit=2":        2,
		"/users?limit=2&page=2": 1,
		"/users?page=5":         0,
	}
	for url, want := range cases {
		rec := call("GET", url, "")
		if rec.Code != 200 {
			t.Errorf("%s: статус %d, ожидался 200", url, rec.Code)
			continue
		}
		got := strings.Count(rec.Body.String(), `"id"`)
		if got != want {
			t.Errorf("%s: вернулось %d пользователей, ожидалось %d", url, got, want)
		}
	}

	badParameters := []string{"page=0", "page=abc", "limit=0", "limit=101", "limit=abc", "age=abc", "age=-1"}
	for _, p := range badParameters {
		rec := call("GET", "/users?"+p, "")
		if rec.Code != 400 {
			t.Errorf("?%s: статус %d, ожидался 400", p, rec.Code)
		}
	}
}

/**GET /usres/{id}*/
func TestGetUser(t *testing.T) {
	prepare(t)
	defer db.Close()
	call("POST", "/users", `{"name":"Анна","age":20}`)

	rec := call("GET", "/users/1", "")
	if rec.Code != 200 {
		t.Errorf("статус %d, ожидался 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Анна") {
		t.Errorf("в ответе нет имени: %s", rec.Body.String())
	}

	rec = call("GET", "/users/999", "")
	if rec.Code != 404 {
		t.Errorf("/users/999: статус %d, ожидался 404", rec.Code)
	}

	rec = call("GET", "/users/abc", "")
	if rec.Code != 400 {
		t.Errorf("/users/abc: статус %d, ожидался 400", rec.Code)
	}
}

/**PUT /users/{id}*/
func TestUpdateUser(t *testing.T) {
	prepare(t)
	defer db.Close()
	call("POST", "/users", `{"name":"Анна","age":20}`)

	rec := call("PUT", "/users/1", `{"name":"Анастасия","age":21}`)
	if rec.Code != 200 {
		t.Errorf("статус %d, ожидался 200", rec.Code)
	}
	rec = call("GET", "/users/1", "")
	if !strings.Contains(rec.Body.String(), "Анастасия") {
		t.Errorf("данные не обновились: %s", rec.Body.String())
	}

	rec = call("PUT", "/users/999", `{"name":"Анна","age":21}`)
	if rec.Code != 404 {
		t.Errorf("/users/999: статус %d, ожидался 404", rec.Code)
	}

	rec = call("PUT", "/users/abc", `{"name":"Анна","age":21}`)
	if rec.Code != 400 {
		t.Errorf("/users/abc: статус %d, ожидался 400", rec.Code)
	}

	badBodies := []string{
		`это не json`,
		`{"name":"","age":21}`,
		`{"name":"Анна","age":0}`,
		`{"name":"Анна","age":121}`,
	}
	for _, body := range badBodies {
		rec := call("PUT", "/users/1", body)
		if rec.Code != 400 {
			t.Errorf("тело %q: статус %d, ожидался 400", body, rec.Code)
		}
	}
}

/**DELETE /users/{id}*/
func TestDeleteUser(t *testing.T) {
	prepare(t)
	defer db.Close()
	call("POST", "/users", `{"name":"Анна","age":20}`)

	rec := call("DELETE", "/users/1", "")
	if rec.Code != 204 {
		t.Errorf("статус %d, ожидался 204", rec.Code)
	}

	rec = call("GET", "/users/1", "")
	if rec.Code != 404 {
		t.Errorf("после удаления статус %d, ожидался 404", rec.Code)
	}

	rec = call("DELETE", "/users/999", "")
	if rec.Code != 404 {
		t.Errorf("/users/999: статус %d, ожидался 404", rec.Code)
	}

	rec = call("DELETE", "/users/abc", "")
	if rec.Code != 400 {
		t.Errorf("/users/abc: статус %d, ожидался 400", rec.Code)
	}
}
