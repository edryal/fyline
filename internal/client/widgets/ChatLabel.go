package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type ChatLabel struct {
	widget.Label
}

func NewChatLabel(text string) *ChatLabel {
	chat := &ChatLabel{}
	chat.ExtendBaseWidget(chat)
	chat.Selectable = true
	chat.Text = text
	chat.Refresh()
	return chat
}

func NewChatLabelWithData(data binding.String) *ChatLabel {
	chat := &ChatLabel{}
	chat.ExtendBaseWidget(chat)
	chat.Selectable = true
	chat.Bind(data)
	return chat
}

// Creates a Bold ChatLabel with a binding
func NewChatUsernameLabel(data binding.String) *ChatLabel {
	chat := &ChatLabel{}
	chat.ExtendBaseWidget(chat)
	chat.Selectable = true
	chat.Bind(data)
	chat.TextStyle = fyne.TextStyle{Bold: true}
	return chat
}

// Creates a dimmed ChatLabel
func NewChatTimestampLabel(text string) *ChatLabel {
	chat := NewChatLabel(text)
	chat.Importance = widget.LowImportance
	return chat
}

// Creates a ChatLabel with word-wrap enabled
func NewChatBodyLabel(text string) *ChatLabel {
	chat := NewChatLabel(text)
	chat.Wrapping = fyne.TextWrapWord
	return chat
}
