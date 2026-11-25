package main

/**Импорты*/
import (
	"fmt"
	"sync"
)

/**Функция увеличавающая счетчик*/
func counter(useMutex bool) int {
	var (
		count int
		mu    sync.Mutex
		wg    sync.WaitGroup
	)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000000; j++ {
				if useMutex {
					mu.Lock()
				}
				count++
				if useMutex {
					mu.Unlock()
				}
			}
		}()
	}
	wg.Wait()
	return count
}

func runTask4() {
	fmt.Println("Задание 4: мьютексы")
	fmt.Println("Без мьютекса:", counter(false))
	fmt.Println("С мьютексом:", counter(true))
}
