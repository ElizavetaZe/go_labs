package main

/**Импорты*/
import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/**Адрес сервера*/
var baseURL = "http://localhost:8080"

/**Данные пользователя*/
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

/**Ответ сервера с ошибкой - 2 задание*/
type ErrorResponse struct {
	Error string `json:"error"`
}

/**Ответ сервера на вход - 4 задание*/
type LoginResponse struct {
	Token string `json:"token"`
}

/**Токен сессии - 4 задание*/
var token string

var scan = bufio.NewScanner(os.Stdin)

/**Чтение строки, введенной пользователем*/
func ask(line string) string {
	fmt.Print(line)
	if !scan.Scan() {
		os.Exit(0)
	}
	return strings.TrimSpace(scan.Text())
}

/**Функция, спрашивающая имя и возраст (для PUT и POST)*/
func askUser() ([]byte, bool) {
	name := ask("Имя: ")
	age, err := strconv.Atoi(ask("Возраст: "))
	if err != nil {
		fmt.Println("Возраст должен быть числом")
		return nil, false
	}
	body, _ := json.Marshal(User{Name: name, Age: age})
	return body, true
}

/**Текст ошибки из ответа сервера*/
func errorMessage(status int, data []byte) string {
	var e ErrorResponse
	if json.Unmarshal(data, &e) == nil && e.Error != "" {
		return e.Error
	}
	return http.StatusText(status)
}

/**Отправляет запрос и обрабатывает ответ; если result не  nil, JSON записывается в него 2 задание*/
func request(method, path string, body []byte, result any) bool {
	req, err := http.NewRequest(method, baseURL+path, bytes.NewReader(body))
	if err != nil {
		fmt.Println("Ошибка: неверный запрос:", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	//сервер недоступен
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка: не удалось связаться с сервером. Проверьте, что он запущен")
		return false
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка: не удалось прочитать ответ сервера")
		return false
	}
	//сервер ответил ошибкой
	if resp.StatusCode >= 400 {
		fmt.Printf("Ошибка %d: %s\n", resp.StatusCode, errorMessage(resp.StatusCode, data))
		return false
	}
	//успех
	if result != nil {
		if err := json.Unmarshal(data, result); err != nil {
			fmt.Println("Ошибка: сервер вернул неожиданный ответ")
			return false
		}
	}
	return true
}

/**Один пользователь - 2 задание*/
func printUser(u User) {
	fmt.Printf("ID:	%d\nИмя: %s\nВозраст: %d\n", u.ID, u.Name, u.Age)
}

/**Список пользователей - 2 задание*/
func printUsers(users []User) {
	if len(users) == 0 {
		fmt.Println("Пользователей нет")
		return
	}
	fmt.Printf("%-5s %-20s %s\n", "ID", "Имя", "Возраст")
	for _, u := range users {
		fmt.Printf("%-5d %-20s %d\n", u.ID, u.Name, u.Age)
	}
	fmt.Println("Всего:", len(users))
}

/**Вывести всех пользователей*/
func listUsers() {
	var users []User
	if request("GET", "/users", nil, &users) {
		printUsers(users)
	}
}

/**Вывести пользователя по id*/
func showUser() {
	if !loggedIn() {
		return
	}
	id := ask("Введите id: ")
	var u User
	if request("GET", "/users/"+id, nil, &u) {
		printUser(u)
	}
}

/**Добавить пользователя*/
func createUser() {
	if !loggedIn() {
		return
	}
	body, ok := askUser()
	if !ok {
		return
	}
	var u User
	if request("POST", "/users", body, &u) {
		fmt.Println("Пользователь добавлен:")
		printUser(u)
	}
}

/**Обновить пользователя*/
func updateUser() {
	if !loggedIn() {
		return
	}
	id := ask("Введите id: ")
	var u User
	if !request("GET", "/users/"+id, nil, &u) {
		return
	}
	fmt.Println("\nТекущие данные:")
	printUser(u)
	fmt.Println("\nНовые данные (Enter - оставить без изменений):")
	if name := ask("Имя: "); name != "" {
		u.Name = name
	}
	if text := ask("Возраст: "); text != "" {
		age, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println("Возраст должен быть числом")
			return
		}
		u.Age = age
	}
	body, _ := json.Marshal(u)
	if request("PUT", "/users/"+id, body, &u) {
		fmt.Println("\nПользователь обновлен:")
		printUser(u)
	}
}

/**Удалить пользователя*/
func deleteUser() {
	if !loggedIn() {
		return
	}
	id := ask("Введите id: ")
	if request("DELETE", "/users/"+id, nil, nil) {
		fmt.Println("Пользователь удален")
	}
}

/**Вход - 4 задание*/
func signIn() {
	login := ask("Логин: ")
	password := ask("Пароль: ")
	body, _ := json.Marshal(map[string]string{"login": login, "password": password})
	var result LoginResponse
	if request("POST", "/login", body, &result) {
		token = result.Token //сохранение токена
		fmt.Println("Вход выполнен")
		fmt.Println("токен:", token)
	}
}

/**Проверка, что пользователь вошев в аккаунт - 4 задание*/
func loggedIn() bool {
	if token == "" {
		fmt.Println("Сначала войдите в аккаунт (пункт 6)")
		return false
	}
	return true
}

/**Выход - 4 задание*/
func signOut() {
	token = ""
	fmt.Println("Вы вышли из аккаунта")
}

func main() {
	for {
		fmt.Println("\nМеню")
		fmt.Println("1. Вывести всех пользователей")
		fmt.Println("2. Вывести пользователя по id")
		fmt.Println("3. Добавить пользователя")
		fmt.Println("4. Обновить пользователя")
		fmt.Println("5. Удалить пользователя")
		fmt.Println("6. Войти")
		fmt.Println("7. Выйти из аккаунта")
		fmt.Println("0. Выход из приложения")
		switch ask("Выберите пункт: ") {
		case "1":
			listUsers()
		case "2":
			showUser()
		case "3":
			createUser()
		case "4":
			updateUser()
		case "5":
			deleteUser()
		case "6":
			signIn()
		case "7":
			signOut()
		case "0":
			return
		default:
			fmt.Println("Неверный пункт меню")
		}
	}
}
