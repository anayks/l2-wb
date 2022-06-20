package pattern

import "math/rand"

// Паттерн состояния хорошо применять, когда есть множество одинаковых действий, которые
// В зависимости от множества разных состояний делают разные вещи
// Паттерн состояния позволяет разделить состояния на классы и уже через них выполнять какие-либо действия
// Как раз в зависимости от состояния
// Это очень похоже предыдущий паттерн, но там стратегию выбирает сам клиент, а здесь "машина" с состоянием
// Сама решает, что делает в зависимости от состояния, которое указано

// Представим себе состояние управления на стрелочки на клавиатуре в качестве интерфейса

type ControlsContext interface {
	ArrowLeft()
	ArrowRight()
	Enter()
}

type PlayerOnFoot struct {
	app *StateApp
} // Состояние игрока пешком

func (p PlayerOnFoot) ArrowLeft() {
	// Заставляем игрока двигаться пешком налево
}

func (p PlayerOnFoot) ArrowRight() {
	// Заставляем игрока двигаться пешком пешком
}

func (p PlayerOnFoot) Enter() {
	// Если игрок нажмет на enter пешком, его с 50% шансов телепортирует в воду или в ближаюшую машину
	r := rand.Intn(100)
	if r > 50 {
		p.app.changeState(PlayerOnVehicle{p.app})
	} else {
		p.app.changeState(PlayerOnWater{p.app})
	}
}

type PlayerOnVehicle struct {
	app *StateApp
} // Состояние игрока пешком

func (p PlayerOnVehicle) ArrowLeft() {
	// Заставляем игрока двигать руль налево
}

func (p PlayerOnVehicle) ArrowRight() {
	// Заставляем игрока двигать руль направо
}

func (p PlayerOnVehicle) Enter() {
	p.app.changeState(PlayerOnFoot{p.app})
	// выходим из машины
}

type PlayerOnWater struct {
	app *StateApp
} // Состояние игрока пешком

func (p PlayerOnWater) ArrowLeft() {
	// Заставляем игрока грести руками налево
}

func (p PlayerOnWater) ArrowRight() {
	// Заставляем игрока грести руками направо
}

func (p PlayerOnWater) Enter() {
	p.app.changeState(PlayerOnFoot{p.app})
	// Выходим из воды
}

type StateApp struct {
	state ControlsContext
}

func (a *StateApp) changeState(c ControlsContext) {
	a.state = c
}

func StateMain() {
	app := &StateApp{}
	app.state = PlayerOnFoot{}

	app.state.ArrowLeft() // Игрок идет налево
	app.state = PlayerOnVehicle{}
	app.state.ArrowRight() // Игрок поворачивает руль направо
	app.state.Enter()      // Выходим из машины
	app.state.Enter()      // Телепортируемся либо в воду, либо в машину
	app.state.Enter()      // Выходим либо из воды, либо из машины
}

// Плюсы
// 1. Избавляет от множества условных операторов
// 2. Упрощает контекст выполнения работы
// 3. Код с определенным состоянием находится в одном месте

// Минусы
// 1. Если состояний мало и они не будут меняться, то смысл реализации паттерна теряется

// Разница между стратегией и состоянием в том, что состояние - это надстройка на стратегий
// и реализация паттерна состояния может сама менять состояния приложения, в то время как
// реализация паттерна стратегии даже не будет знать о других возможных состояних
