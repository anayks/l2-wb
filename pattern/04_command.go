package pattern

// Паттерн команда, в основном, нужен для того, чтобы структурировать логику
// И отделить её от другого слоя логики
// Например, представим, что у нас есть множество разных мест, где нужно вызвать сохранение состояния всего приложения

type App struct {
	SaveApp Command
	sql     MySQL
	saved   int
}

func newApp() App {
	app := App{
		sql: MySQL{},
	}
	saveCommand := SaveCommand{
		app: &app,
	}
	app.SaveApp = &saveCommand
	return App{}
}

func (a *App) moveFile() {
	// что-то делаем
	a.SaveApp.Execute()
}

func (a *App) changeAppName() {
	// что-то делаем
	a.SaveApp.Execute()
}

func (a *App) deleteFile() {
	// что-то делаем
	a.SaveApp.Execute()
}

func (a *App) createFile() {
	// что-то делаем
	a.SaveApp.Execute()
}

type Command interface {
	Execute()
}

type SaveCommand struct {
	app *App
}

type MySQL struct{} // Какой-нибудь абстрактный MySQL

func (m *MySQL) query(s string) {
	// делаем s запрос внутри mysql
}

func (sc *SaveCommand) Execute() {
	sc.app.sql.query("UPDATE App SET a = 1, b = 2, c = 3 WHERE id = 1")
	sc.app.saved++
}

func CommandMain() {
	app := newApp()

	app.createFile()
	app.moveFile()
	app.changeAppName()
}

// Таким образом, нам не нужно писать MySQL запрос или что-то ещё каждый раз,
// Когда мы делаем какое-либо действие
// Мы просто можем спокойно вызвать команду и оно это сделает

// Плюсы
// 1. Благодаря общему интерфейсу, который можно дополнить, можно реализовать стек вызовов команд и манипулировать им - откатывать, вызывать много раз или вызвать последний
// 2. Меньше кода, который будет повторяться, благодаря одной созданной команде, которая может использоваться повсеместно
// 3. Можно собирать сложные команды из более простых

// Минусы:
// 1. Загрязнение программы классами
