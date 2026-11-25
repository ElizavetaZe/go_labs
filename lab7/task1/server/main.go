package main

/**Импорты*/
import (
	"fmt"
	"net"
)

/**Функция, обрабатывающая клиента*/
func process(c net.Conn) {
	defer c.Close()

	data := make([]byte, 1024)
	n, err := c.Read(data)
	if err != nil {
		fmt.Println("Ошибка чтения ", err)
		return
	}
	message := string(data[:n])
	fmt.Println("Получено: ", message)

	c.Write([]byte("Клиент, вы сказали: " + message))
}

/**Главная функция*/
func main() {

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка ", err)
		return
	}
	defer listen.Close()
	fmt.Println("Сервер запущен на порту 8080")

	for {
		connection, err := listen.Accept()
		if err != nil {
			fmt.Println("Ошибка подключения ", err)
			continue
		}
		process(connection)
	}
}
