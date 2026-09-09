package main

/**Импорты*/
import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

/**Главная функция, в которой демонстрируем решение всех задач*/
func main() {
	fmt.Println("Задание 1: Созданиеи запуск горутин")
	var n, count, x int
	fmt.Println("Введите число для факториала: ")
	fmt.Scan(&n)
	fmt.Println("Введите количество случайных значений: ")
	fmt.Scan(&count)
	fmt.Println("Введите число для для суммы ряда: ")
	fmt.Scan(&x)
	go factorial(n)
	go generate(count)
	go sumRow(x)
	time.Sleep(time.Second)

	fmt.Println("Задание 2: Использование каналов для передачи данных")
	ch := make(chan int)
	go fibonachi(ch)
	for x := range ch {
		fmt.Print(x, " ")
	}

	fmt.Println("\nЗадание 3: Применение select для управления данными")
	number_chan := make(chan int)
	result_chan := make(chan string)
	go genRandom(number_chan)
	go even_odd(number_chan, result_chan)
	for i := 0; i < 10; i++ {
		select {
		case a := <-number_chan:
			fmt.Println("Сгенерированное число: ", a)
			time.Sleep(time.Second)
		case b := <-result_chan:
			fmt.Println(b)
		}
	}

	fmt.Println("Задание 4: Синхронизация с помощью мьютексов")
	for i := 1; i < 4; i++ {
		go increment(i)
	}
	time.Sleep(time.Second)
	fmt.Println("Итоговое значение:", count4)

	fmt.Println("Задание 6: Создание пула воркеров")
	var number int
	fmt.Println("Введите количество воркеров:")
	fmt.Scan(&number)

	inputFile := "input.txt"
	outputFile := "output.txt"
	//чтение строк
	lines, err := read(inputFile)
	if err != nil {
		fmt.Println("Ошибка чтения файла ", err)
		return
	}
	//создание каналов
	tasks := make(chan string)
	results := make(chan string, len(lines))
	//запуск пула воркеров
	var wg sync.WaitGroup
	for i := 1; i <= number; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}
	//передаем задачи в канал
	for _, line := range lines {
		tasks <- line
	}
	close(tasks)
	//ожидание завершение всех воркеров
	wg.Wait()
	close(results)

	var resultsList []string
	for res := range results {
		resultsList = append(resultsList, res)
	}

	err = write(outputFile, resultsList)
	if err != nil {
		fmt.Println("Ошибка записи файла ", err)
		return
	}

	fmt.Println("Задание 5: Разработка многопоточного калькулятора")
	req := make(chan Request)
	go calculator(req)
	for {
		var in string
		fmt.Println("Введите первое число или слово 'exit' для выхода")
		fmt.Scan(&in)
		if in == "exit" {
			close(req)
			return
		}
		var a, b float64
		var oper string
		var err error
		a, err = strconv.ParseFloat(in, 64)
		if err != nil {
			fmt.Println("Ошибка. Введите число")
			continue
		}
		fmt.Println("Введите операцию (+, -, *, /):")
		fmt.Scan(&oper)
		if oper != "+" && oper != "-" && oper != "*" && oper != "/" {
			fmt.Println("Ошибка чтения операции")
			continue
		}
		fmt.Println("Введите второе число: ")
		fmt.Scan(&b)
		if oper == "/" && b == 0 {
			fmt.Println("Ошибка.Деление на ноль")
			continue
		}
		reply := make(chan float64)
		req <- Request{A: a, B: b, Oper: oper, Result_chan: reply}
		result := <-reply
		fmt.Printf("%.2f %s %.2f = %.2f\n", a, oper, b, result)
	}
}
