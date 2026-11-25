// curl http://localhost:8080/hello
// curl -X POST http://localhost:8080/data -d '{"Name":"Анна","Age":20}'
package main

/**Импорты*/
import (
	"encoding/json"
	"fmt"
	"net/http"
)

/**Структура для POST-запроса*/
type Data struct {
	Name string
	Age  int
}

/**Функция для GET-запроса*/
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Только GET", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Здравствуйте, клиент")
}

/**Функция для POST-запроса*/
func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST", http.StatusMethodNotAllowed)
		return
	}
	var data Data
	//декодируем json в структуру Data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	fmt.Printf("Получено: %+v\n", data)
	w.Header().Set("Content-Type", "application/json")
	//кодируем в json
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

/**Главная функция простого HTTP-сервера, регестрируем обработчики для двух путей*/
func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/data", dataHandler)
	fmt.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", nil)
}
