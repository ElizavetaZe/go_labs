package main

/**Импорты*/
import (
	"fmt"
	"net"
)

/**Главная функция клиента*/
func main() {
	connection, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка подключения ", err)
		return
	}
	defer connection.Close()

	message := "Здравствуй, сервер"
	connection.Write([]byte(message))
	fmt.Println("Отправлено: ", message)

	buffer := make([]byte, 1024)
	n, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка чтения ", err)
		return
	}
	fmt.Println("Ответ сервера: ", string(buffer[:n]))
}
