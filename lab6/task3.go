package main

/**Импорты*/
import (
	"fmt"
	"math/rand"
	"time"
)

/**Функция, генерирующая случайные значения без сохранения в массив*/
func genRandom(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- rand.Intn(100)
		time.Sleep(time.Millisecond)
	}
	close(ch)
}

/**Функция, проверяющая четность/нечетность*/
func even_odd(num_ch chan int, result chan string) {
	for x := range num_ch {
		if x%2 == 0 {
			result <- fmt.Sprintf("%d четно", x)
		} else {
			result <- fmt.Sprintf("%d нечетно", x)
		}
	}
	close(result)
}
