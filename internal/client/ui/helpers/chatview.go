package helpers

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	"github.com/edryal/fyline/internal/client/assets"
)

type ChannelChatView struct {
	// holds every channel's scroll container (only the active one is shown)
	stack    *fyne.Container
	contents map[string]*fyne.Container
	scrolls  map[string]*container.Scroll
	active   string
}

// creates an empty view. add channels with EnsureChannel.
func NewChannelChatView() *ChannelChatView {
	return &ChannelChatView{
		stack:    container.NewStack(),
		contents: make(map[string]*fyne.Container),
		scrolls:  make(map[string]*container.Scroll),
	}
}

// returns the container to place in the layout
func (v *ChannelChatView) Root() fyne.CanvasObject {
	return v.stack
}

// makes sure a message area exists for the given channel ID.
// safe to call repeatedly (only creates one the first time)
func (v *ChannelChatView) EnsureChannel(channelID string) {
	if _, ok := v.contents[channelID]; ok {
		return
	}

	content := container.NewVBox()
	compact := container.NewThemeOverride(
		content,
		assets.CompactTheme{Theme: theme.DefaultTheme()},
	)
	scroll := container.NewVScroll(compact)
	scroll.Hide() // hidden until activated

	v.contents[channelID] = content
	v.scrolls[channelID] = scroll
	v.stack.Add(scroll)
}

// adds a pre-built message widget to a channel's area
// (creating the channel area if needed) and scrolls that channel to the bottom
func (v *ChannelChatView) AppendTo(channelID string, item fyne.CanvasObject) {
	v.EnsureChannel(channelID)
	v.contents[channelID].Add(item)

	// FIXME: investigate how to more reliably refresh the scroll that holds messages so that chatting is more fluid.
	// the Add() above also updates but doesn't reliably repaint the scroll, so messages only show after
	// the next event (focus, hover on something in the UI), if we were to exlude this explicit refresh
	v.scrolls[channelID].Refresh()
	v.scrolls[channelID].ScrollToBottom()
}

// makes the given channel's area the visible one
func (v *ChannelChatView) Show(channelID string) {
	v.EnsureChannel(channelID)
	if v.active == channelID {
		return
	}

	if s, ok := v.scrolls[v.active]; ok {
		s.Hide()
	}

	v.scrolls[channelID].Show()
	v.active = channelID
}
