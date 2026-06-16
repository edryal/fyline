package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/edryal/fyline/internal/client/widgets"
)

const ApplicationName = "Fyline"

func main() {
	application := app.New()
	windowMain := application.NewWindow(ApplicationName)
	windowMain.Resize(fyne.NewSize(1280, 720))

	chatVScroll := container.NewVScroll(widget.NewLabel("Chat"))

	chatInput := widgets.NewSingleLineRichEntry()
	chatInput.SetPlaceHolder("Enter message...")

	s := binding.NewString()
	_ = s.Set("Size will appear here")

	chatBox := container.NewBorder(nil, chatInput, nil, nil, chatVScroll)

	windowMain.SetContent(chatBox)
	windowMain.ShowAndRun()
}
