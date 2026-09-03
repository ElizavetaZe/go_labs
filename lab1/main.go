package main

/**Импорты*/
import (
	"fmt"
	"time"
)

/**Вывод текущей даты и времени*/
func task1() {
	now := time.Now()
	fmt.Println("Текущее время и дата:", now)
}

/**Создание переменных и выведение их на экран*/
func task2() {
	var i int = 12
	var f float64 = 17.542
	var s string = "hello"
	var b bool = false
	fmt.Println("Переменные различных типов:", i, f, s, b)
}

/**Создание и вывод переменных, используя краткую форму*/
func task3() {
	i, f, s, b := 10, 14.5, "hi", true
	fmt.Println("Переменные различных типов:", i, f, s, b)
}

/**Арифметические операции с двумя целыми числами*/
func task4(a, b int) (sum, diff, mult, div, rem int) {
	if b == 0 {
		fmt.Println("Деление на ноль невозможно")
		return
	}
	sum = a + b
	diff = a - b
	mult = a * b
	div = a / b
	rem = a % b
	return
}

/**Вычисление суммы и разности двух чисел с плавающей запятой*/
func task5(a, b float64) (sum, diff float64) {
	sum = a + b
	diff = a - b
	return
}

/**Вычисление среднего значения трех чисел*/
func task6(a, b, c int) float64 {
	avg := float64(a+b+c) / 3.0
	return avg
}

/**Главная функция, в которой демонстрируем решение всех задач*/
func main() {
	fmt.Println("Задание 1:")
	task1()

	fmt.Println("Задание 2:")
	task2()

	fmt.Println("Задание 3:")
	task3()

	fmt.Println("Задание 4:")
	s, d, m, di, r := task4(15, 4)
	fmt.Printf("Сумма: %d Разность: %d Произведение: %d Частное: %d Остаток: %d\n", s, d, m, di, r)

	fmt.Println("Задание 5:")
	var x, y float64
	fmt.Println("Введите два числа с плавающей точкой:")
	fmt.Scan(&x, &y)
	s1, d1 := task5(x, y)
	fmt.Printf("Сумма: %.2f Разность: %.2f \n", s1, d1)

	fmt.Println("Задание 6:")
	var a, b, c int
	fmt.Println("Введите три целых числа:")
	fmt.Scan(&a, &b, &c)
	fmt.Println("Среднее арифметическое трех чисел:", task6(a, b, c))
}
