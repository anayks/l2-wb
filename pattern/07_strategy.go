package pattern

import (
	"fmt"
	"time"
)

// Паттерн стратегия заключается в том, чтобы через общий интерфейс получать результат
// С помощью разных алгоритмов или методов, но через один общий интерфейс

// Представим, что у нас есть возможность считать цены на товары, и у нас в сервисе есть различные акции
// Мы создаем общий интерфейс для рассчета акции

type PriceStrategy interface {
	getPrice(int) int
}

type NewYearPrice struct{}

func (t NewYearPrice) getPrice(g int) int {
	return int(g / 2)
}

type BasicPrice struct{}

func (t BasicPrice) getPrice(g int) int {
	return g
}

type WeekendPrice struct{}

func (t WeekendPrice) getPrice(g int) int {
	return int(g / 10 * 9)
}

type StrategyApp struct {
	NowStrategy PriceStrategy // Приложение будет хранить структуру текущей стратегии
}

const (
	DEFAULT_SIM_CARD_PRICE = 1000 // Стандартная цена товара
)

func StrategyMain() {
	app := StrategyApp{}
	date := time.Now()

	month := date.Month()
	day := date.Day()

	weekDay := date.Weekday()

	// выбираем стратегию в зависимости от текущего дня
	if month == 1 && day == 1 {
		app.NowStrategy = NewYearPrice{}
	} else if weekDay == 6 {
		app.NowStrategy = WeekendPrice{}
	} else {
		app.NowStrategy = BasicPrice{}
	}

	fmt.Printf("SIM card price by strategy: %v", app.NowStrategy.getPrice(DEFAULT_SIM_CARD_PRICE)) // Просчитываем стоимость сим-карты в зависимости от текущей
}

// Плюсы
// 1. Изоляция реализации конкретной стратегии
// 2. Возможность менять стратегию "на лету"

// Минусы
// 1. Клиенту (разработчику, не связанному с пакетом стратегии) нужно знать, в чём заключаются стратегии
// 2. Если стратегий будет много и они будут меняться, то код засорится классами
