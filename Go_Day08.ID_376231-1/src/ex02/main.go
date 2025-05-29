package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Инициализация приложения (аналог InitApplication)
	myApp := app.New()

	// Создание окна (аналог Window_Create)
	myWindow := myApp.NewWindow("School 21")
	myWindow.SetContent(widget.NewLabel("Hello, School 21!"))

	windowSize := fyne.Size{Width:300, Height: 200,}
	myWindow.Resize(windowSize)

	// Отображение окна и запуск основного цикла (аналог RunApplication)
	myWindow.ShowAndRun()
}