package main

/**Импорты*/
import (
	"fmt"
	"unicode/utf8"
)

/**Функиция для определение четности/нечетности введеного числа*/
func task1(n int) {
	if n%2 == 0 {
		fmt.Printf("Число %d четно\n", n)
	} else {
		fmt.Printf("Число %d нечетно\n", n)
	}
}

/**Функция возвращающая "Positive", "Negative" или "Zero"*/
func task2(n int) string {
	if n > 0 {
		return "Positive"
	}
	if n < 0 {
		return "Negative"
	}
	return "Zero"
}

/**Функция, возвращающая длину строки*/
func task4(s string) {
	fmt.Println(utf8.RuneCountInString(s))
}

/**Структура Rectangle*/
type Rectangle struct {
	Height int
	Width  int
}

/**Метод для вычисления площади прямоугольника*/
func (r Rectangle) Area() int {
	return r.Height * r.Width
}

/**Функция, возвращающая среднее значение двух чисел*/
func task6(a, b int) float64 {
	avg := float64(a+b) / 2.0
	return avg
}

/**Главная функция, в которой демонстрируем решение всех задач*/
func main() {
	fmt.Println("Задание 1: ")
	var number int
	fmt.Println("Введите число: ")
	fmt.Scan(&number)
	task1(number)

	fmt.Printf("Задание 2: \nПроверим на введенном ранее числе\nЧисло %d: ", number)
	result := task2(number)
	fmt.Println(result)

	fmt.Println("Задание 3: ")
	x := 1
	for x <= 10 {
		fmt.Printf("%d ", x)
		x++
	}

	fmt.Println("\nЗадание 4: ")
	var str string
	fmt.Println("Введите строку: ")
	fmt.Scan(&str)
	fmt.Println("Количество символов в строке: ")
	task4(str)

	fmt.Println("Задание 5: ")
	rect := Rectangle{12, 3}
	fmt.Println("Площадь прямоугольника: ", rect.Area())

	fmt.Println("Задание 6:")
	var a, b int
	fmt.Println("Введите два целых числа:")
	fmt.Scan(&a, &b)
	fmt.Println("Среднее арифметическое двух чисел:", task6(a, b))
}
