Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
Ситуация такая же, как и в listing 4
Нам нужно сравнивать конкретно с nil типом указателя на customError
Иначе оно будет не равно nil и выведется error
```

Как можно сделать: (поменять сравнение типов)
```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != (*customError)(nil) {
		println("error")
		return
	}
	println("ok")
}
```

Как можно сделать 2: (поменять тип вывода функции)
```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() (err error) {
	err = &customError{}
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```
