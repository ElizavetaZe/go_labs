package main

/**Импорты*/
import (
	"fmt"
	"math/rand"
)

/**Функция для рассчета факториала*/
func factorial(n int) {
	var mult int = 1
	for n > 1 {
		mult *= n
		n--
	}
	fmt.Println("Факториал числа равен ", mult)
}

/**Генерация случайных чисел*/
func generate(n int) {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(100)
	}
	fmt.Println(arr)
}

/**Функция для вычисления суммы числового ряда*/
func sumRow(n int) {
	sum := 0
	for n > 0 {
		sum += n
		n--
	}
	fmt.Println("Сумма числового ряда", sum)
}
