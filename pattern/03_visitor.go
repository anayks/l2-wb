package pattern

// import "reflect"

// Представим, что у нас есть огромный пул различных элементов
// На примере какой-нибудь игры, представим сущности: автомобиль, игрок, объект на земле, вертолет, чекпоинт и т.д.
// Нам нужно перебрать абсолютно все элементы, связанные с этим пуллом, а они разные!
// Например, раз в час у нас происходит таймер, в котором тикает игровой опыт у игрока, страховка в автомобиле, время объекта на земле и т.д...

var PoolEntity []interface{}

type Player struct{}
type Vehicle struct{}
type Heli struct{}
type Checkpoint struct{}

// Чтобы перебрать элементы, нам нужно ручками перебирать каждый элемент и смотреть на его тип,
// чтобы прям точно быть уверенным, что у него есть функция, которая нам нужна
// есть вариант это делать через switch, как в примере внизу, или через рефлексию, сравнивая типы

// func ActionByType(t interface{}) {
// 	switch t.(type) {
// 	case Player:
// 		{
// 			// Player action
// 			break
// 		}
// 	case Vehicle:
// 		{
// 			// Vehicle action
// 			break
// 		}
// 	case Heli:
// 		{
// 			// Heli action
// 			break
// 		}
// 	case Checkpoint:
// 		{
// 			// Checkpoint action
// 			break
// 		}
// 	}
// }

// func main() {
// 	for _, v := range PoolEntity {
// 		ActionByType(v)
// 	}
// }

// Но есть вариант удобнее...
// В отдельном пакете мы можем реализовать паттерн посетителя, в котором будут сокрыты действия
// Которые мы будем совершать с элементами
// И при этом нам не нужно будет сравнивать по типам, так как по этому паттерну посетитель будет
// Посещать структуру, а структура уже будет вызывать нужный ей метод внутри посетителя

type Visitor interface {
	visitPlayer(*Player)
	visitVehicle(*Vehicle)
	visitHeli(*Heli)
	visitCheckpoint(*Checkpoint)
}

type TimeVisitor struct{}

func (t *TimeVisitor) visitPlayer(p *Player) {
	// Что-то делаем с игроком
}

func (t *TimeVisitor) visitVehicle(p *Vehicle) {
	// Что-то делаем с автомобилем
}

func (t *TimeVisitor) visitHeli(p *Heli) {
	// Что-то делаем с вертолетом
}

func (t *TimeVisitor) visitCheckpoint(p *Checkpoint) {
	// Что-то делаем с чекпоинтом
}

// делаем так, чтобы пул был не с пустым интерфейсом, а с тем, что нужен для посетителя

type Entity interface {
	accept(*TimeVisitor) // Принять посетителя
}

func (p *Player) accept(t *TimeVisitor) {
	t.visitPlayer(p)
}

func (p *Vehicle) accept(t *TimeVisitor) {
	t.visitVehicle(p)
}

func (p *Heli) accept(t *TimeVisitor) {
	t.visitHeli(p)
}

func (p *Checkpoint) accept(t *TimeVisitor) {
	t.visitCheckpoint(p)
}

var NewPoolEntity []Entity

func main() {
	vis := TimeVisitor{} // Создаем посетителя, который перебирает список различных элементов

	for _, v := range NewPoolEntity { // перебираем пул
		v.accept(&vis) // Каждый из типов принимает к себе посетителя и передает ему внутри уже то, что нужно для их типа
	}
}

// Плюсы:
// 1. Возможность сокрытия реализации при обработке сложной структуры или разносортных данных
// 2. Не нужно проверять типы через Switch, что будет работать, хоть немного, но быстрее и более явно
// 3. В структуре можно аккумулировать какие-либо значения при переборе
// 4. Все функции, связанный с конкретной работой, находятся в одном месте в пакете с посетителем

// Минусы
// 1. Если классы в сложной структуре постоянно добавляются, то в каждом нужно реализовывать метод принятия посетителя
// 2. Можно нарушить сокрытие данных и придется менять элементы, если в пакете с Visitor'ом придется пользоваться приватными данными
// 3. Явно больше кода, чем со switch
