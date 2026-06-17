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

	chatContent := container.NewVBox(widget.NewLabel("Chat"))
	chatVScroll := container.NewVScroll(chatContent)

	chatInputEntry := widgets.NewSingleLineRichEntry()
	chatInputEntry.SetPlaceHolder("Enter message...")

	s := binding.NewString()
	_ = s.Set("Size will appear here")

	// Send message when clicking the button
	chatSendButton := widget.NewButton("Send", func() {
		messageLabel := widget.NewLabel(chatInputEntry.Text)
		chatContent.Add(messageLabel)
		chatInputEntry.SetText("")
	})

	// Send message when pressing Shift+Enter
	chatInputEntry.OnSubmitted = func(text string) {
		messageLabel := widget.NewLabel(chatInputEntry.Text)
		chatContent.Add(messageLabel)
		chatInputEntry.SetText("")
	}

	chatInput := container.NewBorder(nil, nil, nil, chatSendButton, chatInputEntry)
	chatBox := container.NewBorder(nil, chatInput, nil, nil, chatVScroll)

	windowMain.SetContent(chatBox)
	windowMain.ShowAndRun()
}
