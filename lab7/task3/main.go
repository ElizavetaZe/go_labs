package main

/**Импорты*/
import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
)

/**Функция, обрабатывающая клиента*/
func process(c net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
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

/**Главная функция сервера. Сервер корректно завершает соединения Ctrl+C*/
func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка ", err)
		return
	}
	fmt.Println("Сервер запущен на порту 8080")
	//канал для сигнала остановки
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	var wg sync.WaitGroup
	go func() {
		<-stop
		fmt.Println("\nПолучен сигнал остановки")
		listen.Close()
	}()

	for {
		connection, err := listen.Accept()
		if err != nil {
			fmt.Println("Ошибка подключения ", err)
			break
		}
		wg.Add(1)
		go process(connection, &wg)
	}
	wg.Wait()
	fmt.Println("Сервер остановлен")
}
