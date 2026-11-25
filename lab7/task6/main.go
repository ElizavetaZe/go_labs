package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const page = `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Чат</title></head>
<body>
<pre id="chat"></pre>
<input id="text">
<button onclick="send()">Отправить</button>
<script>
var ws = new WebSocket("ws://" + location.host + "/ws");
var chat = document.getElementById("chat");
var input = document.getElementById("text");
 
// пришло сообщение от сервера: добавляем его на страницу
ws.onmessage = function (e) {
  chat.textContent += e.data + "\n";
};
 
// отправляем сообщение на сервер
function send() {
  ws.send(input.value);
  input.value = "";
}
</script>
</body>
</html>`

var upgrader = websocket.Upgrader{}
var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

/**Функция, обрабатывающая клиента: каждого в своей горутине*/
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer conn.Close()

	//добавляем клиента
	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		//рассылаем сообщение
		mu.Lock()
		for c := range clients {
			c.WriteMessage(websocket.TextMessage, msg)
		}
		mu.Unlock()
	}

	//убираем клиента, если ушел со страницы
	mu.Lock()
	delete(clients, conn)
	mu.Unlock()
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(page))
	})
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("Чат запущен: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
