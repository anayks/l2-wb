package pattern

import "fmt"

// У нас есть очень много классов, которые как-то относятся к обработке видео

type VideoCodecChanger struct{}

func (vcc *VideoCodecChanger) changeCodec(vf *VideoFile, c string) {
	// Меняем как-нибудь кодек
}

type VideoSoundChanger struct{}

func (vsc *VideoSoundChanger) changeSound(vf *VideoFile, language string) {
	// Меняем как-нибудь звуковой файл
}

type VideoSubtitlesChanger struct{}

func (vsc *VideoSubtitlesChanger) changeSubs(vf *VideoFile, language string) {
	// меняем язык
}

type VideoFile struct{} // сама структура файла

// как-нибудь инициаилизируем файл...
func newVideoFile(path string) *VideoFile {
	return &VideoFile{}
}

// Мы создаем общий класс, который является фасадом для взаимодействия, который будет скрывать множество сложных реализаций других структур, методов и прочего
// Цель фасада - упрощение взаимодействия с множеством структур и "скрытие" реализации от клиента

type VideoConverter struct{}

func (vc *VideoConverter) Convert(path, language, resultCodec string) VideoFile {
	videoFile := newVideoFile(path)
	vsc := &VideoSoundChanger{}
	vsc.changeSound(videoFile, language)
	vsubc := &VideoSubtitlesChanger{}
	vsubc.changeSubs(videoFile, language)
	vcc := &VideoCodecChanger{}
	vcc.changeCodec(videoFile, resultCodec)
	return *videoFile
}

// Вот как у нас будет выглядеть это в коде, максимально простой интерфейс взаимодействия с пакетом, реализация которого от нас может быть скрыта в другом пакете

func FacadeMain() {
	converter := VideoConverter{}
	resultFile := converter.Convert("test-file.mp4", "en", "ogg")
	fmt.Printf("Наш файл: %v", resultFile)
}

// Плюсы паттерна:
// 1. Можно скрыть реализацию от клиента
// 2. Упрощение взаимодействия

// Минусы:
// 1. Может превратиться в "божественный" объект, в котором кашей может быть намешано множество реализаций
