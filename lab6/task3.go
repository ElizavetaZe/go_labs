package main

/**Импорты*/
import (
	"fmt"
	"math/rand"
	"time"
)

/**Функция, генерирующая случайные значения без сохранения в массив*/
func genRandom(out chan<- int) {
	for i := 0; i < 10; i++ {
		out <- rand.Intn(100)
		time.Sleep(100 * time.Millisecond)
	}
	close(out)
}

/**Функция, проверяющая четность/нечетность*/
func evenOdd(in <-chan int, nums chan<- int, msgs chan<- string) {
	for x := range in {
		nums <- x
		if x%2 == 0 {
			msgs <- fmt.Sprintf("%d четно", x)
		} else {
			msgs <- fmt.Sprintf("%d нечетно", x)
		}
	}
	close(nums)
	close(msgs)
}

func runTask3() {
	fmt.Println("Задание 3: select")
	raw := make(chan int)
	nums := make(chan int)
	msgs := make(chan string)
	go genRandom(raw)
	go evenOdd(raw, nums, msgs)

	for nums != nil || msgs != nil {
		select {
		case v, ok := <-nums:
			if !ok {
				nums = nil
				continue
			}
			fmt.Println("Сгенерированное число:", v)
		case m, ok := <-msgs:
			if !ok {
				msgs = nil
				continue
			}
			fmt.Println(m)
		}
	}
}
