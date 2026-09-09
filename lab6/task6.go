package main

/**Импорты*/
import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

/**Функция, возвращающая перевернутутю строку*/
func reserveStr(s string) (result string) {
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		result += string(runes[i])
	}
	return result
}

/**Функция, читающая файл построчно*/
func read(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		lines = append(lines, scan.Text())
	}
	return lines, scan.Err()
}

/**Функция, записывающая результаты в новый файл*/
func write(fileName string, lines []string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

/**Функция worker*/
func worker(i int, tasks <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		time.Sleep(time.Millisecond * 50)
		revers := reserveStr(task)
		result := fmt.Sprintf("Воркер %d: %s -> %s", i, task, revers)
		results <- result
	}
}
