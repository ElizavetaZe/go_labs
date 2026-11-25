package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func main() {
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		fmt.Println("Не удалось загрузить сертификат клиента", err)
		return
	}

	caPEM, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		fmt.Println("Не удалось прочитать certs/ca.crt:", err)
		return
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caPEM)

	congig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		ServerName:   "localhost",
		MinVersion:   tls.VersionTLS12,
	}

	if len(os.Args) > 1 && os.Args[1] == "nocert" {
		congig.Certificates = nil
	}

	conn, err := tls.Dial("tcp", "localhost:8080", congig)
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		return
	}
	defer conn.Close()

	state := conn.ConnectionState()
	server := state.PeerCertificates[0].Subject.CommonName
	fmt.Printf("Соединение защищено (%s). Сертификат сервер проверен: %s\n", tls.VersionName(state.Version), server)

	fmt.Println("Введите сообщение:")
	scan := bufio.NewScanner(os.Stdin)
	scan.Scan()
	message := scan.Text()

	conn.Write([]byte(message))
	fmt.Println("Отправлено:", message)

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return
	}
	fmt.Println("Ответ сервера:", string(buffer[:n]))
}
