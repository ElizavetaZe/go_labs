package main

import (
	"errors"
	"fmt"
	"sync"
)

/**Структура запроса пользователя*/
type Request struct {
	A     float64
	B     float64
	Oper  string
	Reply chan Response
}

/**Структура ответа*/
type Response struct {
	Result float64
	Err    error
}

/**Функция, реализующая работы калькулятора*/
func calculator(requests <-chan Request) {
	for r := range requests {
		go func(r Request) {
			var resp Response
			switch r.Oper {
			case "+":
				resp.Result = r.A + r.B
			case "-":
				resp.Result = r.A - r.B
			case "*":
				resp.Result = r.A * r.B
			case "/":
				if r.B == 0 {
					resp.Err = errors.New("деление на ноль")
				} else {
					resp.Result = r.A / r.B
				}
			default:
				resp.Err = errors.New("неизвестная операция")
			}
			r.Reply <- resp
		}(r)
	}
}

func client(id int, a float64, oper string, b float64, requesets chan<- Request, wg *sync.WaitGroup) {
	defer wg.Done()
	reply := make(chan Response)
	requesets <- Request{A: a, B: b, Oper: oper, Reply: reply}
	resp := <-reply
	if resp.Err != nil {
		fmt.Printf("Клиент %d: %.2f %s %.2f -> ошибка: %v\n", id, a, oper, b, resp.Err)
		return
	}
	fmt.Printf("Клиент %d: %.2f %s %.2f = %.2f\n", id, a, oper, b, resp.Result)
}

func runTask5() {
	fmt.Println("Задание 5: многопоточный калькулятор")
	requests := make(chan Request)
	go calculator(requests)
	var wg sync.WaitGroup
	wg.Add(4)
	go client(1, 10, "+", 5, requests, &wg)
	go client(2, 10, "-", 3, requests, &wg)
	go client(3, 6, "*", 7, requests, &wg)
	go client(4, 8, "/", 0, requests, &wg)
	wg.Wait()

	close(requests)
}
