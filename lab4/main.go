package main

/**Импорты*/
import (
	"fmt"
	"strings"
)

/**Функция, возвращающая средний возраст людей в карте*/
func avg(m map[string]int) {
	sum := 0
	count := 0
	for _, value := range m {
		sum += value
		count += 1
	}
	result := float64(sum) / float64(count)
	fmt.Printf("Средний возраст людей из карты: %.2f\n", result)
}

/**Функция, возвращающая перевернутый массив*/
func reverseArr(a []int) []int {
	left := 0
	right := len(a) - 1
	for right > left {
		a[left], a[right] = a[right], a[left]
		left++
		right--
	}
	return a
}

/**Главная функция, в которой демонстрируем решение всех задач*/
func main() {
	fmt.Println("Задание 1:")
	people := make(map[string]int)
	people["Jane"] = 18
	people["Alice"] = 31
	people["Mike"] = 52
	fmt.Println(people)

	fmt.Println("Задание 2:")
	avg(people)

	fmt.Println("Задание 3:")
	delete(people, "Alice")
	fmt.Println("Карта после удаления: ", people)

	fmt.Println("Задание 4:")
	var str string
	fmt.Println("Введите строку: ")
	fmt.Scan(&str)
	result := strings.ToUpper(str)
	fmt.Println("Строка в верхнем регистре: ", result)

	fmt.Println("Задание 5:")
	var n, x int
	fmt.Println("Введите количество чисел, которое хотите ввести: ")
	fmt.Scan(&n)
	fmt.Println("Введите числа: ")
	sum := 0
	for n > 0 {
		fmt.Scan(&x)
		sum += x
		n--
	}
	fmt.Println("Сумма введенных чисел: ", sum)

	fmt.Println("Задание 6:")
	fmt.Println("Введите размер массива: ")
	fmt.Scan(&n)
	arr := make([]int, n)
	fmt.Println("Введите элементы массива: ")
	for i := range arr {
		fmt.Scan(&arr[i])
	}
	fmt.Println("Исходный массив: ", arr)
	fmt.Println("Перевернутый массив: ", reverseArr(arr))
}
