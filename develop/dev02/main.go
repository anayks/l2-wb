package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Будем использовать паттерн состояния
// Всего есть 3 состояния чтения, когда:
// 1. ReadSymb - первое состояние, машина может читать только бекслеш или символ
// 2. ReadAnything - состояние, когда машина может прочитать, что угодно.
// Если это число - это количество, если это символ - то обрабатываем символ, если бекслеш - переходим к 3 состоянию
// 3. ReadingEscapeValue - чтение любого символа после бекслеша.

type unpack struct { // Текущая распаковка как структура
	lastSymbol rune
	newString  []rune
	state      unpackingState
}

func newUnpack() *unpack {
	unpacking := unpack{
		newString: make([]rune, 0),
	} // инициализируем структуру распаковки
	unpacking.setState(&readSymb{}) // задаем стандарное состояние
	return &unpacking
}

func (u *unpack) setLastSymbol(id int, r rune, array []rune) { // Установка последнего символа
	if u.lastSymbol != rune(0) {
		u.newString = append(u.newString, u.lastSymbol)
	}
	u.lastSymbol = r
	if id == len(array)-1 {
		u.newString = append(u.newString, u.lastSymbol)
	}
}

func (u *unpack) insertLastSymbol() {
	u.newString = append(u.newString, u.lastSymbol)
	u.lastSymbol = rune(0)
}

func (u *unpack) setState(v unpackingState) { // Установка
	u.state = v
	v.setApp(u)
}

type unpackingState interface { // интерфейс обработчика состояния
	execute(int, rune, []rune) error // Обработать id с указанным символом и строкой обработки
	setApp(*unpack)                  // Указать для текущего состояния приложение, как метод управления состоянием
}

type readSymb struct {
	app *unpack // Изнутри состояния будем управлять, ссылаясь на приложение
}

func (r *readSymb) execute(id int, value rune, array []rune) error {
	if value == '\\' {
		r.app.setState(&readingEscapeValue{})
		return nil
	} else if unicode.IsNumber(value) {
		return fmt.Errorf("error on reading symbol #%d (%v): state is reading only symbols", id, value)
	}

	r.app.setLastSymbol(id, value, array)
	r.app.setState(&ReadAnything{})
	return nil
}

func (r *readSymb) setApp(a *unpack) {
	r.app = a
}

type ReadAnything struct {
	app *unpack
}

func (r *ReadAnything) execute(id int, value rune, array []rune) error {
	if value == '\\' {
		r.app.insertLastSymbol()
		r.app.setState(&readingEscapeValue{})
		return nil
	} else if unicode.IsNumber(value) {
		last := r.app.lastSymbol
		r.app.newString = append(r.app.newString, []rune(strings.Repeat(string(last), int(value-'0')))...)
		r.app.lastSymbol = rune(0)
		r.app.setState(&readSymb{})
		return nil
	}

	r.app.setLastSymbol(id, value, array)
	r.app.setState(&ReadAnything{})
	return nil
}

func (r *ReadAnything) setApp(a *unpack) {
	r.app = a
}

type readingEscapeValue struct {
	app *unpack
}

func (r *readingEscapeValue) execute(id int, value rune, array []rune) error {
	r.app.setLastSymbol(id, value, array)
	r.app.setState(&ReadAnything{})
	return nil
}

func (r *readingEscapeValue) setApp(a *unpack) {
	r.app = a
}

func UnpackString(v string) (string, error) {
	unpacking := newUnpack() // инициализируем стандартное состояне
	runes := []rune(v)       // переводим строку в слайс рун

	for id, v := range runes {
		if err := unpacking.state.execute(id, v, runes); err != nil { // если при обработке возникла ошибка, возвращаем пустую строку и ошибку
			return "", err
		}
	}

	return string(unpacking.newString), nil // возвращаем строку без ошибку
}
