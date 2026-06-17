package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Make the compiler check if we actually match the interfaces
var (
	_ desktop.Hoverable = (*MessageBlock)(nil)
)

type MessageBlock struct {
	widget.BaseWidget
	background *canvas.Rectangle
	content    fyne.CanvasObject
	hovered    bool
}

func NewMessageBlock(content fyne.CanvasObject) *MessageBlock {
	m := &MessageBlock{
		background: canvas.NewRectangle(color.Transparent),
		content:    content,
	}
	m.ExtendBaseWidget(m)
	return m
}

// CreateRenderer stacks the background behind the content.
func (m *MessageBlock) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewStack(m.background, m.content)
	return widget.NewSimpleRenderer(c)
}

// Mouse is inside the MessageBlock, so create the hover effect
func (m *MessageBlock) MouseIn(_ *desktop.MouseEvent) {
	m.hovered = true
	m.background.FillColor = theme.Color(theme.ColorNameHover)
	m.background.Refresh()
}

// Don't care about MouseMoved, we just need to satisfy the Hoverable interface
func (m *MessageBlock) MouseMoved(_ *desktop.MouseEvent) {}

// Mouse is no longer on top of the MessageBlock, so we can remove hover effect
func (m *MessageBlock) MouseOut() {
	m.hovered = false
	m.background.FillColor = color.Transparent
	m.background.Refresh()
}
