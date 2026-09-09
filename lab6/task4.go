package main

/**Импорты*/
import (
	"fmt"
	"sync"
	"time"
)

/**Переменные для счетчика и синхронизации для 4 задания*/
var count4 int
var mutex sync.Mutex

/**Функция увеличавающая счетчик*/
func increment(n int) {
	for i := 0; i < 10; i++ {
		mutex.Lock()
		count4++
		fmt.Printf("Горутина %d: счетчик %d\n", n, count4)
		mutex.Unlock()
		time.Sleep(time.Millisecond * 10)
	}
}
