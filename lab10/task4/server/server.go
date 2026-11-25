package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
)

func process(c net.Conn) {
	defer c.Close()

	tlsConn := c.(*tls.Conn)
	if err := tlsConn.Handshake(); err != nil {
		fmt.Println("Ошибка TLS:", err)
		return
	}

	client := tlsConn.ConnectionState().PeerCertificates[0].Subject.CommonName
	fmt.Print("Подключился клиент:", client)

	data := make([]byte, 1024)
	n, err := c.Read(data)
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return
	}
	message := string(data[:n])
	fmt.Println("\nПолучено:", message)

	c.Write([]byte("Сообщение получено: " + message))
}

func main() {
	cert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		fmt.Println("Не удалось загрузить сертификат сервера", err)
		return
	}
	caPEM, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		fmt.Println("Не удалось прочитать certs/ca.crt:", err)
		return
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caPEM)

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	listener, err := tls.Listen("tcp", ":8080", config)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Защищенный сервер запущен на порту 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка подключения:", err)
			continue
		}
		go process(conn)
	}
}
