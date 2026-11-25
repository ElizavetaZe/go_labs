package main

import (
	"bufio"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
)

/**Файлы с ключами*/
const (
	privateFile = "private.pem"
	publicFile  = "public.pem"
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

/**Создание пары ключей RSA и их сохранение в файлы*/
func generateKeys() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})
	if err := os.WriteFile(privateFile, privPEM, 0600); err != nil {
		return err
	}
	return os.WriteFile(publicFile, pubPEM, 0644)
}

/**Чтение закрытого ключа из файла*/
func loadPrivateKey() (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(privateFile)
	if err != nil {
		return nil, errors.New("не удалось прочитать " + privateFile + " (сначала создайте ключи)")
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("файл " + privateFile + " поврежден")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("файл " + privateFile + " поврежден")
	}
	return key, nil
}

/**Чтение окртытого ключа из файла*/
func loadPublicKey() (*rsa.PublicKey, error) {
	data, err := os.ReadFile(publicFile)
	if err != nil {
		return nil, errors.New("не удалось прочитать " + publicFile + " (сначала создайте ключи)")
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("файл " + publicFile + " поврежден")
	}
	key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("файл " + publicFile + " поврежден")
	}
	return key, nil
}

/**Подписывание сообщения: считает хэш и подписывает закрытым ключом*/
func sign(privateKey *rsa.PrivateKey, message string) ([]byte, error) {
	hash := sha256.Sum256([]byte(message))
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
}

/**Проверка подписи открытым ключом*/
func verify(publicKey *rsa.PublicKey, message string, signature []byte) bool {
	hash := sha256.Sum256([]byte(message))
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature) == nil
}

/**Результат проверки подписи*/
func printResult(ok bool) {
	if ok {
		fmt.Println("Подпись верна: сообщение подписано владельцем закрытого ключа и не изменено")
	} else {
		fmt.Println("Подпись НЕВЕРНА: сообщение изменено или подписано другим ключом")
	}
}

/**Пункт 1: создание пары ключей*/
func keysMenu() {
	fmt.Println("Создание ключей")
	if err := generateKeys(); err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Ключи созданы и сохранены в файлы:")
	fmt.Println(" ", privateFile, "- закрытый ключ (секретный, никому не передавать)")
	fmt.Println(" ", publicFile, "- открытый ключ (его можно передавать всем)")
}

/**Пункт 2: подписать сообщение закрытым ключом*/
func signMenu() {
	privateKey, err := loadPrivateKey()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	message := readLine("Введите сообщение: ")
	signature, err := sign(privateKey, message)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Подпись:", hex.EncodeToString(signature))
}

/**Пункт 3: проверить подпись открытым ключом*/
func verifyMenu() {
	publicKey, err := loadPublicKey()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	message := readLine("Введите сообщение: ")
	signatureText := strings.TrimSpace(readLine("Введите подпись: "))
	signature, err := hex.DecodeString(signatureText)
	if err != nil {
		fmt.Println("Подпись должна состоять из символов 0-9 и a-f")
		return
	}
	printResult(verify(publicKey, message, signature))
}

/**Пункт 4: передача подписанных сообещений*/
func demo() {
	privateKey, err := loadPrivateKey()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	publicKey, err := loadPublicKey()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	message := "Женя, отправила тебе посылку. Лиза"
	signature, err := sign(privateKey, message)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Лиза подписала сообщение своим закрытым ключом и оправила Жене:")
	fmt.Println(" Сообщение: ", message)
	fmt.Println(" Подпись: ", hex.EncodeToString(signature)[:40]+"...(начало)")

	//сообщение дошло без изменений
	fmt.Println("\nЖеня проверяет сообщение открытым ключом Лизы:")
	printResult(verify(publicKey, message, signature))
	//подменили текст
	changed := "Женя, отправила тебе 5 посылок. Лиза"
	fmt.Println("\nКто-то подменил сообщение:")
	fmt.Println(" Сообщение:", changed)
	printResult(verify(publicKey, changed, signature))
	//подпись поддельного сообщение своим ключом
	fmt.Println("\nКто-то отправил поддельное сообщение со своей подписью:")
	otherKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fakeSignature, err := sign(otherKey, changed)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	printResult(verify(publicKey, changed, fakeSignature))
}

func main() {
	for {
		fmt.Println("\nЦифровая подпись RSA")
		fmt.Println("1. Создать пару ключей")
		fmt.Println("2. Подписать сообщение")
		fmt.Println("3. Проверить подпись")
		fmt.Println("4. Пример: передача сообщений между двумя сторонами")
		fmt.Println("0. Выход")
		switch strings.TrimSpace(readLine("Выберите пункт: ")) {
		case "1":
			keysMenu()
		case "2":
			signMenu()
		case "3":
			verifyMenu()
		case "4":
			demo()
		case "0":
			return
		default:
			fmt.Println("Неверный пункт")
		}
	}
}
