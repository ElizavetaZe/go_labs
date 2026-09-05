package mathutils

/**Функция для вычисление факториала числа n*/
func Factorial(n int) int {
	var mult int = 1
	for n > 1 {
		mult *= n
		n--
	}
	return mult
}
