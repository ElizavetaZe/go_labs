package main

/**Импорты*/
import (
	"fmt"
	"math"
)

/**Структура с именем и возрастом*/
type Person struct {
	name string
	age  int
}

/**Метода для вывода информации о человеке*/
func (p Person) String() string {
	return fmt.Sprintf("Имя: %s, Возраст: %d", p.name, p.age)
}

/**Метод, увеличивающий возраст*/
func (p *Person) birthday() {
	p.age++
}

/**Структура Circle*/
type Circle struct {
	radius int
}

/**Метод для вычисления площади круга*/
func (c Circle) areaCircle() {
	result := math.Pi * math.Pow(float64(c.radius), 2)
	fmt.Printf("Площадь круга: %.3f\n", result)
}

/**Структура Rectangle*/
type Rectangle struct {
	height int
	width  int
}

/**Интерфейс Shape*/
type Shape interface {
	Area() float64
}

/**Реализация интерфейса Shape для Circle*/
func (c Circle) Area() float64 {
	return math.Pi * math.Pow(float64(c.radius), 2)
}

/**Реализация интерфейса Shape для Rectangle*/
func (r Rectangle) Area() float64 {
	return float64(r.height) * float64(r.width)
}

/**Функция, выводящая площадь каждого объекта*/
func areaAll(s []Shape) {
	for i, shape := range s {
		area := shape.Area()
		fmt.Printf("Площадь %d объекта: %.2f\n", i+1, area)
	}
}

/**Структура Book*/
type Book struct {
	title  string
	author string
	year   int
	page   int
}

/**Интерфейс Stringer*/
type Stringer interface {
	String() string
}

/**Реализация интерфейса Stringer для Book*/
func (b Book) String() string {
	return fmt.Sprintf("Название: %s, Автор: %s, Год: %d, Количество страниц: %d", b.title, b.author, b.year, b.page)
}

/**Главная функция, в которой демонстрируем решение всех задач*/
func main() {
	fmt.Println("Задание 1:")
	p1 := Person{name: "Alice", age: 12}
	p2 := Person{name: "Rob", age: 34}
	fmt.Println("Информация о людях:\n", p1, "\n", p2)

	fmt.Println("Задание 2:")
	p2.birthday()
	fmt.Println("Информация о людях:\n", p1, "\n", p2)

	fmt.Println("Задание 3:")
	c1 := Circle{radius: 2}
	c1.areaCircle()

	fmt.Println("Задание 4:")
	c := Circle{radius: 5}
	r := Rectangle{height: 10, width: 5}
	fmt.Printf("Площадь круга: %.2f\nПлощадь прямоугольника: %.2f\n", c.Area(), r.Area())

	fmt.Println("Задание 5:")
	shapes := []Shape{c, r}
	fmt.Println("Все площади:")
	areaAll(shapes)

	fmt.Println("Задание 6:")
	book := Book{title: "Война и мир", author: "Лев Толстой", year: 1869, page: 1984}
	fmt.Println("Информация о книге:\n", book)
}
