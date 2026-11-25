package mathutils

import "fmt"

/**Функция для вычисление факториала числа n*/
func Factorial(n int) int {
	if n < 0 {
		fmt.Println("Число не может быть отрицательным")
		return 0
	}
	var mult int = 1
	for n > 1 {
		mult *= n
		n--
	}
	return mult
}
