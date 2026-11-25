package main

/**Импорты*/
import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

/**Функция, реализующая middleware - оборачивает обработчик и логирует метод, URL и время выполнения*/
func middleware(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		n.ServeHTTP(w, r)
		fmt.Printf("Метод: %s Путь: %s Время выполнения: %s\n", r.Method, r.URL.Path, time.Since(start))
	})
}

/**Главная функция простого HTTP-сервера*/
func main() {
	//создание маршрутизатора
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/data", dataHandler)

	fmt.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", middleware(mux))
}
