// сервер включить в task1 для задания 2
// для задания 3 включить сервер в task3
package main

/**Импорты*/
import (
	"bufio"
	"fmt"
	"net"
	"os"
)

/**Главная функция клиента с вводом сообщения из консоли*/
func main() {
	connection, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка подключения ", err)
		return
	}
	defer connection.Close()

	scan := bufio.NewScanner(os.Stdin)
	fmt.Println("Введите сообщение: ")
	scan.Scan()
	message := scan.Text()

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
