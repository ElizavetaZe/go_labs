package main

/**Импорты*/
import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

/**Функция, возвращающая перевернутутю строку*/
func reverseStr(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

/**Функция, читающая файл построчно*/
func readLines(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

/**Функция worker*/
func worker(id int, tasks <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range tasks {
		time.Sleep(50 * time.Millisecond)
		results <- fmt.Sprintf("Воркер %d: %s -> %s", id, line, reverseStr(line))
	}
}

func runTask6() {
	fmt.Println("Задание 6: пул воркеров")
	var n int
	fmt.Println("Введите количество воркеров:")
	fmt.Scan(&n)
	lines, err := readLines("input.txt")
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}
	tasks := make(chan string)
	results := make(chan string, len(lines))

	var wg sync.WaitGroup
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	for _, line := range lines {
		tasks <- line
	}
	close(tasks)

	wg.Wait()
	close(results)

	var out []string
	for res := range results {
		fmt.Println(res)
		out = append(out, res)
	}

	data := strings.Join(out, "\n") + "\n"
	if err := os.WriteFile("output.txt", []byte(data), 0644); err != nil {
		fmt.Println("Ошибка записи файла:", err)
	}
}
