package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const ApplicationName = "Fyline"

func main() {
	app := app.New()
	windowMain := app.NewWindow(ApplicationName)

	message := widget.NewLabel("Skraaaaaaaa")
	button := widget.NewButton("Update", func() {
		formatted := time.Now().Format("Time: 03:04:05")
		message.SetText(formatted)
	})

	newWindowButton := widget.NewButton("Open new window", func() {
		spawnedWindow := app.NewWindow("New Window")
		spawnedWindow.SetContent(widget.NewLabel("I was spawned using a button"))
		spawnedWindow.Resize(fyne.NewSize(200, 200))
		spawnedWindow.Show()
	})

	canvasText := canvas.NewText("First message", color.NRGBA{G: 0xff, A: 0xff})
	createAutoUpdatingText(canvasText)

	windowMain.SetContent(container.NewVBox(message, button, newWindowButton, canvasText))
	windowMain.Show()

	app.Run()
}

func createAutoUpdatingText(canvasText *canvas.Text) {
	// Run this in the background / another thread
	go func() {
		// Every second the ticker ticks and it updates the text with the current time
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			fyne.Do(func() {
				canvasText.Text = time.Now().Format(time.TimeOnly)
				canvasText.Refresh()
			})
		}
	}()
}
