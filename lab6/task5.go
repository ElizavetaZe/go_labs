package main

/**Структура запроса пользователя*/
type Request struct {
	A           float64
	B           float64
	Oper        string
	Result_chan chan float64
}

/**Функция, реализующая работы калькулятора*/
func calculator(requests chan Request) {
	for r := range requests {
		var result float64
		switch r.Oper {
		case "+":
			result = r.A + r.B
		case "-":
			result = r.A - r.B
		case "*":
			result = r.A * r.B
		case "/":
			result = r.A / r.B
		}
		r.Result_chan <- result
	}
}
