package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/edryal/fyline/internal/client/assets"
	netclient "github.com/edryal/fyline/internal/client/net"
	"github.com/edryal/fyline/internal/client/ui/components"
	"github.com/edryal/fyline/internal/client/ui/helpers"
	"github.com/edryal/fyline/internal/client/ui/widgets"
	"github.com/edryal/fyline/internal/protocol"
)

const ApplicationName = "Fyline"
const Username = "Catalin"
const ChannelsPanelOffset = 0.2
const ServerURL = "ws://localhost:8080/ws"

func main() {
	application := app.New()
	application.Settings().SetTheme(assets.CompactTheme{Theme: theme.DefaultTheme()})
	applicationWindow := application.NewWindow(ApplicationName)
	applicationWindow.Resize(fyne.NewSize(1280, 720))

	// channels sidebar
	channels := protocol.SeedChannels
	sidebar := components.NewChannelsSidebar(channels)

	// per-channel chat view
	chatView := helpers.NewChannelChatView()
	for _, ch := range channels {
		chatView.EnsureChannel(ch.ID)
	}

	// message input row
	entry := widgets.NewSingleLineRichEntry()
	entry.SetPlaceHolder("Enter message...")
	sendButton := widget.NewButton("Send", nil)
	inputRow := container.NewBorder(nil, nil, nil, sendButton, entry)
	chatArea := container.NewBorder(nil, inputRow, nil, nil, chatView.Root())

	// create the split
	bundledPanels := container.NewHSplit(sidebar.Root(), chatArea)
	bundledPanels.SetOffset(ChannelsPanelOffset)
	applicationWindow.SetContent(bundledPanels)

	// switching channels swaps the visible message area
	sidebar.OnSelect = func(channelID string) {
		chatView.Show(channelID)
	}

	// activate the first channel so something is visible on launch
	if len(channels) > 0 {
		sidebar.Select(channels[0].ID)
		chatView.Show(channels[0].ID)
	}

	// connect client to the server
	client := netclient.New(ServerURL, Username, netclient.Handlers{
		OnChat: func(msg protocol.ChatMessage) {
			fmt.Printf("[OnChat] channelID=%q user=%q body=%q\n", msg.ChannelID, msg.Username, msg.Body)
			fyne.Do(func() {
				renderChatMessage(chatView, msg)
			})
		},
		OnSystem: func(sys protocol.System) {
			fyne.Do(func() {
				renderSystemMessage(chatView, sidebar.ActiveID(), sys)
			})
		},
		OnStatus: func(s netclient.Status) {
			fyne.Do(func() {
				applicationWindow.SetTitle(ApplicationName + " - " + s.String())
			})
		},
	})
	client.Start()
	defer client.Close()

	send := func() {
		input := strings.TrimSpace(entry.Text)
		if input == "" {
			return
		}
		client.Send(sidebar.ActiveID(), input)
		entry.SetText("")
	}

	sendButton.OnTapped = send
	entry.OnSubmitted = func(string) {
		send()
	}

	applicationWindow.ShowAndRun()
}

// routes an incoming message to the right channel's area by ChannelID needs to run on the UI thread
func renderChatMessage(v *helpers.ChannelChatView, msg protocol.ChatMessage) {
	ts := msg.SentAt.Local().Format("3:04 PM")

	header := container.NewHBox(
		widgets.NewChatUsernameLabelText(msg.Username),
		widgets.NewChatTimestampLabel(ts),
	)

	body := container.NewVBox(
		header,
		widgets.NewChatBodyLabel(msg.Body),
	)

	block := widgets.NewMessageBlock(body)
	spaced := container.New(layout.NewCustomPaddedLayout(0, 0, 0, 0), block)

	v.AppendTo(msg.ChannelID, spaced)
}

// appends a dimmed server notice to the active channel. needs to run on the UI thread
func renderSystemMessage(v *helpers.ChannelChatView, activeID string, sys protocol.System) {
	label := widgets.NewChatTimestampLabel(sys.Text)
	v.AppendTo(activeID, label)
}
