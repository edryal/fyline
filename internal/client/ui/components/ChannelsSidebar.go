package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"github.com/edryal/fyline/internal/client/ui/widgets"
	"github.com/edryal/fyline/internal/protocol"
)

type ChannelsSidebar struct {
	// VBox of rows basically
	column *fyne.Container

	// channel ID -> row
	buttons map[string]*widgets.ChannelButton

	// channel IDs in display order
	order    []string
	activeID string

	// OnSelect fires with the channel ID when a channel is clicked.
	// The caller (usually main) uses this to switch the visible chat view.
	OnSelect func(channelID string)
}

// builds the sidebar from an initial set of channels.
func NewChannelsSidebar(channels []protocol.Channel) *ChannelsSidebar {
	s := &ChannelsSidebar{
		column:  container.NewVBox(),
		buttons: make(map[string]*widgets.ChannelButton),
	}
	s.SetChannels(channels)
	return s
}

// returns the container to place in the window/split.
func (s *ChannelsSidebar) Root() fyne.CanvasObject {
	return container.NewVScroll(s.column)
}

// rebuilds the rows from a fresh channel list. Called on initial
// load and whenever the server sends an updated set (added/removed channels).
func (s *ChannelsSidebar) SetChannels(channels []protocol.Channel) {
	s.column.RemoveAll()
	s.buttons = make(map[string]*widgets.ChannelButton)
	s.order = s.order[:0]

	for _, ch := range channels {
		id := ch.ID
		btn := widgets.NewChannelButton(ch.Name)
		btn.OnTapped = func() {
			s.selectByID(id)
		}

		s.buttons[id] = btn
		s.order = append(s.order, id)
		s.column.Add(btn)
	}

	// keep the previously active channel selected if it still exists; else pick
	// the first one so there's always a valid active channel.
	if _, ok := s.buttons[s.activeID]; !ok && len(s.order) > 0 {
		s.activeID = ""
		s.selectByID(s.order[0])
	} else if s.activeID != "" {
		s.buttons[s.activeID].SetActive(true)
	}
}

// rename updates a single channel's displayed name
// TODO: when the server broadcasts a rename, call this I guess
func (s *ChannelsSidebar) Rename(channelID, newName string) {
	if btn, ok := s.buttons[channelID]; ok {
		btn.SetText(newName)
	}
}

// returns the currently selected channel ID.
func (s *ChannelsSidebar) ActiveID() string {
	return s.activeID
}

// select programmatically activates a channel (e.g. on initial load).
func (s *ChannelsSidebar) Select(channelID string) {
	s.selectByID(channelID)
}

func (s *ChannelsSidebar) selectByID(channelID string) {
	if channelID == s.activeID {
		return
	}

	if prev, ok := s.buttons[s.activeID]; ok {
		prev.SetActive(false)
	}

	if btn, ok := s.buttons[channelID]; ok {
		btn.SetActive(true)
		s.activeID = channelID
		if s.OnSelect != nil {
			s.OnSelect(channelID)
		}
	}
}
