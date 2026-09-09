package main

/**Функция, генерирующая последовательность Фибоначчи*/
func fibonachi(ch chan int) {
	a, b := 0, 1
	for i := 0; i < 10; i++ {
		ch <- a
		a, b = b, a+b
	}
	close(ch)
}
