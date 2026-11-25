package main

/**Импорты*/
import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/**Функция для рассчета факториала*/
func factorial(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	res := 1
	for i := 2; i <= n; i++ {
		res *= i
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("Факториал числа равен ", res)
}

/**Генерация случайных чисел*/
func generate(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		fmt.Println("Случайное число:", rand.Intn(100))
		time.Sleep(150 * time.Millisecond)
	}
}

/**Функция для вычисления суммы числового ряда*/
func sumRow(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0
	for i := 1; i < n; i++ {
		sum += i
		time.Sleep(120 * time.Millisecond)
	}
	fmt.Println("Сумма числового ряда", sum)
}

func runTask1() {
	fmt.Println("Задание 1: Создание и запуск горутин")
	var n, count, x int
	fmt.Println("Введите число для факториала, количество случайных чисел и число для суммы ряда:")
	fmt.Scan(&n, &count, &x)
	var wg sync.WaitGroup
	wg.Add(3)
	go factorial(n, &wg)
	go generate(count, &wg)
	go sumRow(x, &wg)
	wg.Wait()
}
