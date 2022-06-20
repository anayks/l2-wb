package pattern

import "fmt"

// Представим себе такую сложную вещь, как компьютер с совокупностью элементов,
// в котором может быть множество разных элементов
// Сделаем такую вещь, как заказ компьютера

type ComputerOrder struct {
	RAM         []string
	Disks       []string
	CPU         []string
	GPU         []string
	USBHubs     []string
	Motherboard []string
	PowerBlock  []string
	OC          []string
	Monitors    []string
}

// Как мы можем увидеть, здесь есть множество элементов, которые даже не обязательно могут присутствовать
// Если мы попытаемся сделать конструктор этого заказа компьютера (фабрикой), мы увидим, что это в принципе довольно тяжело

func newComputerOrder(RAM, Disks, USBHubs, OC, Monitors, CPU, GPU, Motherboard, PowerBlock, Mouse, Keyboard []string) ComputerOrder {
	return ComputerOrder{
		RAM:         RAM,
		Disks:       Disks,
		CPU:         CPU,
		GPU:         GPU,
		USBHubs:     USBHubs,
		Motherboard: Motherboard,
		PowerBlock:  PowerBlock,
		OC:          OC,
		Monitors:    Monitors,
	}
}

// Чтобы не создавать такой огромный конструктор
// (с учетом того, что у нас реально есть множество различных компьютеров,
// у которых что-то может быть, а может не быть),
// мы можем воспользоваться паттерном Builder

// В нём мы можем установить необходимый порядок установки элементов по необходимости

// Для этого, создаем общий интерфейс строителя

type Builder interface {
	addRam(string)
	addDisk(string)
	addUSBHub(string)
	addOC(string)
	addMonitor(string)
	addCPU(string)
	addGPU(string)
	addMotherboard(string)
	addPowerBlock(string)

	reset()
	getResult() ComputerOrder
}

// А сейчас мы создадим структуру, которая будет создавать только компьютеры MSI

type MSIBuilder struct {
	result ComputerOrder
}

// Просто представим, что у нас MSIBuilder реализует весь этот функционал,
// приклеивая к названию надпись MSI (как в примере внизу)

func (m *MSIBuilder) addRam(r string) {
	m.result.RAM = append(m.result.RAM, "MSI "+r)
}

// Тут просто возвращаем результат работы

func (m *MSIBuilder) getResult() ComputerOrder {
	return m.result
}

// вот тут просто представим реализацию...

func (m *MSIBuilder) addDisk(r string)        {}
func (m *MSIBuilder) addUSBHub(r string)      {}
func (m *MSIBuilder) addOC(r string)          {}
func (m *MSIBuilder) addMonitor(r string)     {}
func (m *MSIBuilder) addCPU(r string)         {}
func (m *MSIBuilder) addGPU(r string)         {}
func (m *MSIBuilder) addMotherboard(r string) {}
func (m *MSIBuilder) addPowerBlock(r string)  {}
func (m *MSIBuilder) reset(r string)          {}

// Как это будет выглядеть?

func BuilderMain() {
	MSIBuilder := &MSIBuilder{}
	MSIBuilder.addCPU("Intel Core i9")
	MSIBuilder.addGPU("Nvidia Geforce RTX 3060")
	result := MSIBuilder.getResult()

	// По сути, здесь у нас структура заказа компьютера, в которой есть только процеcсор и видеокарта
	// Нам не нужно использовать для этого огромный конструктор, в котором есть множество пустых данных
	fmt.Printf("%v", result)
}

// Плюсы:
// 1. Пошаговое создание/изменение структуры
// 2. Можно использовать один код для создания разных структур

// Минусы:
// 1. Введение дополнительных классов
