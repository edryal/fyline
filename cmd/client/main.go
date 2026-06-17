package main

import (
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/edryal/fyline/internal/client/assets"
	"github.com/edryal/fyline/internal/client/widgets"
)

const ApplicationName = "Fyline"
const Username = "Catalin"

func main() {
	application := app.New()
	windowMain := application.NewWindow(ApplicationName)
	windowMain.Resize(fyne.NewSize(1280, 720))

	// TODO: move this stuff either in another file for configurable variables
	// or make them ENV variables
	currentUsername := binding.NewString()
	currentUsername.Set(Username)

	chatContent := container.NewVBox(widget.NewLabel("Chat"))
	compactChat := container.NewThemeOverride(chatContent, assets.CompactTheme{Theme: theme.DefaultTheme()})
	chatVScroll := container.NewVScroll(compactChat)

	chatInputEntry := widgets.NewSingleLineRichEntry()
	chatInputEntry.SetPlaceHolder("Enter message...")

	// Send message when clicking the button
	chatSendButton := widget.NewButton("Send", func() {
		sendMessageAndClearEntry(currentUsername, chatInputEntry, chatContent)
	})

	// Send message when pressing Shift+Enter
	chatInputEntry.OnSubmitted = func(text string) {
		sendMessageAndClearEntry(currentUsername, chatInputEntry, chatContent)
	}

	chatInput := container.NewBorder(nil, nil, nil, chatSendButton, chatInputEntry)
	chatBox := container.NewBorder(nil, chatInput, nil, nil, chatVScroll)

	windowMain.SetContent(chatBox)
	windowMain.ShowAndRun()
}

// TODO: add validator that will highlight wth red the entry when it cannot send the message
func sendMessageAndClearEntry(currentUsername binding.String, chatInputEntry *widgets.SingleLineRichEntry, chatContent *fyne.Container) {
	input := strings.TrimSpace(chatInputEntry.Text)
	if input == "" {
		return
	}

	// TODO: when we'll implement history and we'll have to display time for older messages
	// Do a check if currentTime < Today's date, show the exact date + time
	// currentTime will have to be like currentUsername, a binding, not raw string
	currentTime := time.Now().Format(time.Kitchen)

	header := container.NewHBox(
		widgets.NewChatUsernameLabel(currentUsername),
		widgets.NewChatTimestampLabel(currentTime),
	)

	// Play with themeing padding to lower padding between header and body
	message := container.NewVBox(
		header,
		widgets.NewChatBodyLabel(input),
	)

	// Message becomes a block that can be hovered
	block := widgets.NewMessageBlock(message)

	// TODO: for later when we add more elements we'll most probably want extra padding
	// Was used to add bottom padding between messages
	spacedMessage := container.New(
		layout.NewCustomPaddedLayout(0, 0, 0, 0),
		block,
	)

	chatContent.Add(spacedMessage)
	chatInputEntry.SetText("")
}
