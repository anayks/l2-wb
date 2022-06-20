package pattern

import "fmt"

// Цепочка вызовов позволяет вызывать функции друг за другом буквально в пределах одной строчки
// Функция возвращает сразу объект, к которому мы можем обращаться и мы можем сразу к результату вызова
// функции сразу вызвать следующую функцию
// Можно сделать это на примере числа

type ChainValue struct {
	value int
}

func (c ChainValue) addValue(v int) ChainValue {
	c.value += v
	return c
}

func (c ChainValue) removeValue(v int) ChainValue {
	c.value -= v
	return c
}

func (c ChainValue) multiplyValue(v int) ChainValue {
	c.value += v
	return c
}

func (c ChainValue) divideValue(v int) ChainValue {
	c.value += v
	return c
}

func (c ChainValue) print() ChainValue {
	fmt.Printf("\nResult value: %v", c.value)
	return c
}

// Как это будет выглядеть?

func ChainMain() {
	c := ChainValue{}
	c.addValue(3).divideValue(3).removeValue(1).addValue(2).multiplyValue(2)
}

// Функция возвращает результат в виде объекта, к которому мы можем вызывать следующую функцию - это и есть цепочка

// Плюсы:
// 1. Можно сделать иммутабельность данных, таким образом, чтобы каждый вызов создавал новый объект (может быть полезно, чтобы сохранять состояние начального объекта)
// 2. Можем сохранять последовательность действий в структуре и откатывать её по необходимости
// 3. Решает проблему вложенности функций a(b(c())) путем последовательного вызова

// Минусы:
// 1. Иногда код может быть ОЧЕНЬ длинным, нужно знать меру
