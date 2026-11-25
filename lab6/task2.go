package main

import (
	"fmt"
	"sync"
)

/**Функция, генерирующая последовательность Фибоначчи*/
func fibonachi(ch chan<- int) {
	a, b := 0, 1
	for i := 0; i < 10; i++ {
		ch <- a
		a, b = b, a+b
	}
	close(ch)
}

func printer(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for x := range ch {
		fmt.Print(x, " ")
	}
	fmt.Println()
}

func runTask2() {
	fmt.Println("Задание 2: каналы")

	ch := make(chan int)
	var wg sync.WaitGroup

	wg.Add(1)
	go fibonachi(ch)
	go printer(ch, &wg)
	wg.Wait()
}
