package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

/**Вычисление хэша строки; возвращает его в виде строки из символов*/
func hashString(text, algorithm string) string {
	data := []byte(text)

	switch algorithm {
	case "MD5":
		sum := md5.Sum(data)
		return hex.EncodeToString(sum[:])
	case "SHA-256":
		sum := sha256.Sum256(data)
		return hex.EncodeToString(sum[:])
	case "SHA-512":
		sum := sha512.Sum512(data)
		return hex.EncodeToString(sum[:])
	}
	return ""
}

/**Определение алгоритма по длине хэша*/
func algorithmByLength(hash string) string {
	switch len(hash) {
	case 32:
		return "MD5"
	case 64:
		return "SHA-256"
	case 128:
		return "SHA-512"
	}
	return ""
}

/**Выбор алгоритма*/
func askAlgorithm() string {
	fmt.Println("Алгоритмы: 1 - MD5, 2 - SHA-256, 3 - SHA-512")
	switch strings.TrimSpace(readLine("Выберите: ")) {
	case "1":
		return "MD5"
	case "2":
		return "SHA-256"
	case "3":
		return "SHA-512"
	}
	fmt.Println("Неверный выбор алгоритма")
	return ""
}

/**Вычисление хэша введенной строки*/
func calcHach() {
	algorithm := askAlgorithm()
	if algorithm == "" {
		return
	}
	text := readLine("Введите строку: ")
	fmt.Printf("%s: %s\n", algorithm, hashString(text, algorithm))
}

/**Проверка соответствия*/
func checkHash() {
	text := readLine("Введите строку: ")
	given := strings.TrimSpace(readLine("Введите хэш: "))
	algorithm := algorithmByLength(given)
	if algorithm == "" {
		fmt.Printf("неизвестная длиина хэша (%d символов)", len(given))
		return
	}
	actual := hashString(text, algorithm)
	fmt.Println("Определён алгоритм:", algorithm)
	fmt.Println("Хэш введённой строки:", actual)
	if strings.EqualFold(actual, given) {
		fmt.Println("Хэши совпадают")
	} else {
		fmt.Println("Хэши не совпадают: данные изменены или хэш неверный")
	}
}

var scanner = bufio.NewScanner(os.Stdin)

/**Чтение строки с клавиатуры*/
func readLine(line string) string {
	fmt.Print(line)
	if !scanner.Scan() {
		os.Exit(0)
	}
	return scanner.Text()
}

func main() {
	for {
		fmt.Println("\nХэширование")
		fmt.Println("1. Высилить хэш строки")
		fmt.Println("2. Проверить целостность")
		fmt.Println("0. Выход")
		switch strings.TrimSpace(readLine("Выберите пункт: ")) {
		case "1":
			calcHach()
		case "2":
			checkHash()
		case "0":
			return
		default:
			fmt.Println("Неверный пункт")
		}
	}
}
