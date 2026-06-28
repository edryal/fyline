package main

import (
	"flag"
	"os"
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
	"github.com/edryal/fyline/internal/debug"
	"github.com/edryal/fyline/internal/protocol"
)

const ApplicationName = "Fyline"
const ChannelsPanelOffset = 0.2

// username comes from -user flag, then FYLINE_USER, then the default
func resolveUsername() string {
	var flagUser string
	flag.StringVar(&flagUser, "user", "", "username for this client")
	flag.Parse()

	if flagUser != "" {
		return flagUser
	}

	if env := os.Getenv("FYLINE_USER"); env != "" {
		return env
	}

	return "Catalin"
}

func resolveServerURL() string {
	if env := os.Getenv("FYLINE_SERVER"); env != "" {
		return env
	}
	return "ws://localhost:8080/ws"
}

func main() {
	username := resolveUsername()
	serverURL := resolveServerURL()
	debug.Info("starting client", "user", username, "server", serverURL)

	application := app.New()
	application.Settings().SetTheme(assets.CompactTheme{Theme: theme.DefaultTheme()})
	applicationWindow := application.NewWindow(ApplicationName + " - " + username)
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
		debug.Debug("channel switched", "channelID", channelID)
		chatView.Show(channelID)
	}

	// activate the first channel so something is visible on launch
	if len(channels) > 0 {
		sidebar.Select(channels[0].ID)
		chatView.Show(channels[0].ID)
	}

	// connect client to the server
	client := netclient.New(serverURL, username, netclient.Handlers{
		OnChat: func(msg protocol.ChatMessage) {
			debug.Debug("recv chat", "channelID", msg.ChannelID, "user", msg.Username, "body", msg.Body)
			fyne.Do(func() {
				renderChatMessage(chatView, msg)
			})
		},
		OnSystem: func(sys protocol.System) {
			debug.Debug("recv system", "text", sys.Text)
			fyne.Do(func() {
				renderSystemMessage(chatView, sidebar.ActiveID(), sys)
			})
		},
		OnStatus: func(s netclient.Status) {
			debug.Info("connection status", "status", s.String())
			fyne.Do(func() {
				applicationWindow.SetTitle(ApplicationName + " - " + username + " - " + s.String())
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
		debug.Debug("send chat", "channelID", sidebar.ActiveID(), "body", input)
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
