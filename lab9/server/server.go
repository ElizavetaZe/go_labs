package main

/**Импорты*/
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

/**Подключение базы данных*/
const baza = "postgres://postgres:ВАШ_ПАРОЛЬ@localhost:5432/postgres?sslmode=disable"

/**Пользователь*/
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

/**Подключение к базе*/
var db *sql.DB

/**Ошибка для вывода клиенту - задание 3*/
type AppErr struct {
	Status  int
	Message string
}

func (e *AppErr) Error() string { return e.Message }

var errorNotFound = &AppErr{http.StatusNotFound, "Пользователь не найден"}

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

/**Функция преобразующая ошибки в ответ клиенту - задание 3*/
func handle(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}
		if appErr, ok := err.(*AppErr); ok {
			writeJSON(w, appErr.Status, map[string]string{"error": appErr.Message})
			return
		}
		fmt.Println("Внутренняя ошибка:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Внутренняя ошибка сервера"})
	}
}

/**Отправление ответа в формате JSON*/
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

/**Получение ID из адреса*/
func getID(r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return 0, &AppErr{http.StatusBadRequest, "Неверный id"}
	}
	return id, nil
}

/**Проверка данных пользователя - задание 3*/
func validate(u User) error {
	if strings.TrimSpace(u.Name) == "" {
		return &AppErr{http.StatusBadRequest, "Имя не может быть пустым"}
	}
	if u.Age < 1 || u.Age > 120 {
		return &AppErr{http.StatusBadRequest, "Возраст должен быть от 1 до 120"}
	}
	return nil
}

/**Чтение пользователя из тела запроса - задание 3*/
func readUser(r *http.Request) (User, error) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return u, &AppErr{http.StatusBadRequest, "Неверный JSON"}
	}
	return u, validate(u)
}

/**Функция читает числовой параметр запроса - задание 4*/
func readPar(r *http.Request, key string, def int) (int, error) {
	s := r.URL.Query().Get(key)
	if s == "" {
		return def, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return 0, &AppErr{http.StatusBadRequest, "Неверный параметр " + key}
	}
	return n, nil
}

/**GET /users - список всех пользователей*/
func getUsers(w http.ResponseWriter, r *http.Request) error {
	//фильтрация - 4 задание
	name := r.URL.Query().Get("name")
	age, err := readPar(r, "age", 0)
	if err != nil {
		return err
	}

	//пагинация - 4 задание
	page, err := readPar(r, "page", 1)
	if err != nil {
		return err
	}
	limit, err := readPar(r, "limit", 10)
	if err != nil {
		return err
	}
	if page < 1 || limit < 1 || limit > 100 {
		return &AppErr{http.StatusBadRequest, "page должен быть от 1, limit от 1 до 100"}
	}
	offset := (page - 1) * limit

	rows, err := db.Query(`SELECT id, name, age FROM users 
		WHERE ($1::text = '' OR name ILIKE '%' || $1::text || '%')
		AND ($2::int = 0 OR age = $2::int)
		ORDER BY id
		LIMIT $3 OFFSET $4`, name, age, limit, offset)
	if err != nil {
		return err
	}
	defer rows.Close()
	list := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			return err
		}
		list = append(list, u)
	}
	writeJSON(w, http.StatusOK, list)
	return nil
}

/**GET /users/{id} - один пользователь по id*/
func getUser(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	var u User
	err = db.QueryRow("SELECT id, name, age FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Age)
	if err == sql.ErrNoRows {
		return errorNotFound
	}
	if err != nil {
		return err
	}
	writeJSON(w, http.StatusOK, u)
	return nil
}

/**POST /users - добавить пользователя*/
func addUser(w http.ResponseWriter, r *http.Request) error {
	u, err := readUser(r)
	if err != nil {
		return err
	}
	err = db.QueryRow("INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id", u.Name, u.Age).Scan(&u.ID)
	if err != nil {
		return err
	}
	writeJSON(w, http.StatusCreated, u)
	return nil
}

/**PUT /users/{id} - изменить пользователя по id*/
func updateUser(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	u, err := readUser(r)
	if err != nil {
		return err
	}
	res, err := db.Exec("UPDATE users SET name = $1, age = $2 WHERE id = $3", u.Name, u.Age, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errorNotFound
	}
	u.ID = id
	writeJSON(w, http.StatusOK, u)
	return nil
}

/**DELETE /users/{id} - удалить пользователя по id*/
func deleteUser(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	res, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errorNotFound
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

/**Маршрутизатор со всеми маршрутами - задание 5*/
func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", handle(getUsers))
	mux.HandleFunc("GET /users/{id}", handle(getUser))
	mux.HandleFunc("POST /users", handle(addUser))
	mux.HandleFunc("PUT /users/{id}", handle(updateUser))
	mux.HandleFunc("DELETE /users/{id}", handle(deleteUser))
	return mux
}

/**Главная функция - сервер*/
func main() {
	var err error
	db, err = sql.Open("postgres", baza)
	if err != nil {
		fmt.Println("Ошибка ", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT NOT NULL, age INT NOT NULL)`)
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных ", err)
		return
	}
	fmt.Println("Сервер запущен: http://localhost:8080")
	http.ListenAndServe(":8080", newAuthMux())
}
