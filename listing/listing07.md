Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v:= <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:
```
У нас есть два канала, которые генерируются с помощью функции, и туда записываются значения
Значения записываются в случайное время раз от 0 до 1000 мсек
После того как данные были записаны, каналы a и b закрываются
В функции merge мы через select читаем с этих двух каналов
Данные читаются с этих двух каналов, после чего посыпятся нули, так как нет проверки на закрытость канала,
И чтение будет получать стандартное значение из первого варианта (то есть, ноль)

Решение: 
1. Закрывать канал с, чтобы функция main завершила программу
2. Чтобы сделать 1, нам нужно выйти из цикла и закрыть канал, для этого мы будем проверять, обнулены ли оба канала
Если они обнулены -- выходим и закрываем канал c
```

Само решение
```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					break
				}
				c <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					break
				}
				c <- v
			}
			if a == nil && b == nil {
				close(c)
				return
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Println(v)
	}
}
```
