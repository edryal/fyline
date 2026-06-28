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

var (
	_ desktop.Hoverable = (*ChannelButton)(nil)
	_ fyne.Tappable     = (*ChannelButton)(nil)
)

type ChannelButton struct {
	widget.BaseWidget

	background *canvas.Rectangle
	label      *widget.Label

	hovered bool
	active  bool

	// OnTapped fires when the row is clicked
	OnTapped func()
}

func NewChannelButton(name string) *ChannelButton {
	b := &ChannelButton{
		background: canvas.NewRectangle(color.Transparent),
		label:      widget.NewLabel(name),
	}
	b.ExtendBaseWidget(b)
	return b
}

// updates the displayed channel name (used when a channel is renamed)
func (b *ChannelButton) SetText(name string) {
	b.label.SetText(name)
}

// marks this row as the selected channel (or not) and repaint
func (b *ChannelButton) SetActive(active bool) {
	b.active = active
	b.applyState()
}

func (b *ChannelButton) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewStack(b.background, b.label)
	return widget.NewSimpleRenderer(c)
}

// --- interface stuff ---

func (b *ChannelButton) Tapped(_ *fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *ChannelButton) MouseIn(_ *desktop.MouseEvent) {
	b.hovered = true
	b.applyState()
}

func (b *ChannelButton) MouseMoved(_ *desktop.MouseEvent) {
}

func (b *ChannelButton) MouseOut() {
	b.hovered = false
	b.applyState()
}

func (b *ChannelButton) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

// picks the background color for the current state
func (b *ChannelButton) applyState() {
	switch {
	case b.active:
		b.background.FillColor = theme.Color(theme.ColorNameSelection)
	case b.hovered:
		b.background.FillColor = theme.Color(theme.ColorNameHover)
	default:
		b.background.FillColor = color.Transparent
	}
	b.background.Refresh()
}
