package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var scanner = bufio.NewScanner(os.Stdin)

/**Чтение строки с клавиатуры*/
func readLine(line string) string {
	fmt.Print(line)
	if !scanner.Scan() {
		os.Exit(0)
	}
	return scanner.Text()
}

/**Превращение ключа любой длины в ключ длиной ровна 32 байта*/
func makeKey(password string) []byte {
	sum := sha256.Sum256([]byte(password))
	return sum[:]
}

/**Создание шифра AES в режиме GCM*/
func newGCM(password string) (cipher.AEAD, error) {
	block, err := aes.NewCipher(makeKey(password))
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

/**Шифрование строки; возвращает в виде 16ричной строки*/
func encrypt(text, password string) (string, error) {
	gcm, err := newGCM(password)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	encrypted := gcm.Seal(nil, nonce, []byte(text), nil)
	result := append(nonce, encrypted...)
	return hex.EncodeToString(result), nil
}

/**Дешифрование*/
func decrypt(encryptedHex, password string) (string, error) {
	data, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", errors.New("зашифрованная строка должна состоять из символов 0-9 и a-f")
	}
	gcm, err := newGCM(password)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("Зашифрованная строка слишком короткая")
	}
	nonce := data[:gcm.NonceSize()]
	encrypted := data[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", errors.New("не удалось расшифровать: неверный ключ или данные изменены")
	}
	return string(plain), nil
}

/**Зашифровать строку*/
func encryptMenu() {
	text := readLine("Введите стркоку: ")
	password := readLine("Введите ключ: ")
	if password == "" {
		fmt.Println("Ключ не может быть пустым")
		return
	}
	result, err := encrypt(text, password)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Зашифрованная строка:", result)
}

/**Расшифровать строку*/
func decryptMenu() {
	encrypted := strings.TrimSpace(readLine("Введите зашифрованную строку:"))
	password := readLine("Введите ключ: ")
	plain, err := decrypt(encrypted, password)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Расшифрованная строка:", plain)
}

func main() {
	for {
		fmt.Println("\nШифрование AES")
		fmt.Println("1. Зашифровать строку")
		fmt.Println("2. Расшифровать строку")
		fmt.Println("0. Выход")
		switch strings.TrimSpace(readLine("Выберите пункт: ")) {
		case "1":
			encryptMenu()
		case "2":
			decryptMenu()
		case "0":
			return
		default:
			fmt.Println("Неверный пункт")
		}
	}
}
