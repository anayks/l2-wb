package pattern

// Паттерн фабрики полезен в 2 случаях:
// 1. Переиспользование данных
// 2. Инициализация структуры в зависимости от ситуации

// 1.
// Представим, что у нас есть обращение к базе данных. Соединение с базой данных может оборваться

type FMySQL struct{}

func (m *FMySQL) Ping() bool {
	return true
}

var appConnect *FMySQL

func GetMySQLConnect() *FMySQL {
	if appConnect == nil { // Если соединение не существует, возвращаем новое соединение
		appConnect = &FMySQL{}
		return appConnect
	}

	if appConnect.Ping() == false { // Если соединение не стабильно и не пингуется, возвращаем новое соединение
		appConnect = &FMySQL{}
		return appConnect
	}

	return appConnect // Всё работает корректно, возвращаем старое соединение, чтобы не создавать новое
}

// 2.
// Представим, что у нас то же соединение MySQL зависит от драйвера операционной системы

type Driver interface {
	Query()
}

type MySQLOptimized struct {
	MySQL
	Driver Driver
}

type WinDriver struct{}

func (d WinDriver) Query() {}

type MacOSDriver struct{}

func (d MacOSDriver) Query() {}

type LinuxDriver struct{}

func (d LinuxDriver) Query() {}

var appConnect2 *MySQLOptimized

func getNowOS() string {
	return "win" // "macos", "linux"
}

func getMySQLOptimized() *MySQLOptimized {
	os := getNowOS()

	appConnect2 = &MySQLOptimized{}

	if os == "win" { // Устанавливаем драйвер в зависимости от ОС
		appConnect2.Driver = WinDriver{}
	} else if os == "macos" {
		appConnect2.Driver = MacOSDriver{}
	} else if os == "linux" {
		appConnect2.Driver = LinuxDriver{}
	}

	return appConnect2
}

// Плюсы
// 1. Избавляемся от конкретных типов структур и от инициализации с условиями, перенося это все в отдельный пакет, например
// 2. Код инициализации находится в одном месте
// 3. Упрощается добавление новых типов, например, драйверов

// Минусы
// 1. Может получиться большая иерархия классов с фабриками или большое количество функций, реализующих их
