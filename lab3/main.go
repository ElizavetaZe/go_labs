package main

/**Импорты*/
import (
	"fmt"
	"math/rand"
	"unicode/utf8"

	"github.com/ElizavetaZe/go_labs/lab3/mathutils"
	"github.com/ElizavetaZe/go_labs/lab3/stringutils"
)

/**Функция для создания массива, заполнения значениями и вывод на экран*/
func createArr() {
	var arr [5]int
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Intn(100)
	}
	fmt.Println(arr)
}

/**Функция для создания среза строк и поиска самой длиной*/
func sliceStr() {
	words := []string{"hello", "Samara", "abracadabra"}
	max := 0
	index := 0
	for i := 0; i < len(words); i++ {
		len_word := utf8.RuneCountInString(words[i])
		if len_word > max {
			max = len_word
			index = i
		}
	}
	fmt.Println("Самая длинная строка: ", words[index], " количество символов: ", max)
}

/**Главная функция, в которой демонстрируем решение всех задач*/
func main() {
	fmt.Println("Задание 1-2: ")
	var number int
	fmt.Println("Введите число: ")
	fmt.Scan(&number)
	fmt.Println("Факториал введеннного числа", mathutils.Factorial(number))

	fmt.Println("Задание 3: ")
	var str string
	fmt.Println("Введите строку: ")
	fmt.Scan(&str)
	fmt.Println("Перевернутая строка: ", stringutils.ReserveStr(str))

	fmt.Println("Задание 4: ")
	fmt.Println("Массив: ")
	createArr()

	fmt.Println("Задание 5: ")
	arr := [7]int{1, 23, 15, 10, 7, 21, 19}
	slice := arr[2:5]
	fmt.Println("Срез массива: ", slice)
	slice = append(slice, 50, 0)
	fmt.Println("Массив после добавления элементов: ", slice)
	i := 2
	slice = append(slice[:i], slice[i+1:]...)
	fmt.Println("Массив после удаления элемента: ", slice)

	fmt.Println("Задание 6: ")
	sliceStr()
}
