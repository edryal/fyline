// Package helpers builds UI panels for the Fyline client. It is deliberately
// "dumb": it constructs and returns widgets and containers, and knows nothing
// about networking, rendering incoming messages, or send behavior. That keeps
// it free of any import cycle (it never imports the net client) and means you
// can reshape layouts here without touching wiring code — main.go attaches all
// behavior to the widgets these functions return.
package helpers

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/edryal/fyline/internal/client/assets"
	"github.com/edryal/fyline/internal/client/ui/widgets"
)

// ChatPanel holds the pieces of the chat UI that the caller needs to wire up.
// Returned by NewChatPanel. Fields are exported so main can attach handlers
// (OnTapped, OnSubmitted) and render messages into Content. The struct lives
// here next to its constructor, but main is free to hold it, embed it, or pull
// out individual fields — nothing in this package depends on how it's used.
type ChatPanel struct {
	// Root is the fully assembled container to drop into the window/split.
	Root *fyne.Container
	// Content is the VBox that message widgets get appended to.
	Content *fyne.Container
	// Scroll wraps Content; call ScrollToBottom() after appending.
	Scroll *container.Scroll
	// Entry is the text input.
	Entry *widgets.SingleLineRichEntry
	// SendButton is the send button; attach OnTapped in main.
	SendButton *widget.Button
}

// NewChatPanel builds the chat message panel: a scrollable message area on top
// and an input row (entry + send button) pinned to the bottom. It wires no
// behavior — the caller sets SendButton.OnTapped and Entry.OnSubmitted, and
// renders messages by appending to Content.
func NewChatPanel() *ChatPanel {
	content := container.NewVBox(widget.NewLabel("Chat"))

	// Local theme override keeps the message area compact regardless of the
	// app-wide theme.
	compact := container.NewThemeOverride(
		content,
		assets.CompactTheme{Theme: theme.DefaultTheme()},
	)
	scroll := container.NewVScroll(compact)

	entry := widgets.NewSingleLineRichEntry()
	entry.SetPlaceHolder("Enter message...")

	sendButton := widget.NewButton("Send", nil) // OnTapped attached by caller

	inputRow := container.NewBorder(nil, nil, nil, sendButton, entry)
	root := container.NewBorder(nil, inputRow, nil, nil, scroll)

	return &ChatPanel{
		Root:       root,
		Content:    content,
		Scroll:     scroll,
		Entry:      entry,
		SendButton: sendButton,
	}
}
