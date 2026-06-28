package netclient

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/edryal/fyline/internal/debug"
	"github.com/edryal/fyline/internal/protocol"
)

// describes the connection state, reported via OnStatus
type Status int

const (
	StatusDisconnected Status = iota
	StatusConnecting
	StatusConnected
)

func (s Status) String() string {
	switch s {
	case StatusConnecting:
		return "connecting"
	case StatusConnected:
		return "connected"
	default:
		return "disconnected"
	}
}

type Handlers struct {
	OnChat   func(protocol.ChatMessage)
	OnSystem func(protocol.System)
	OnStatus func(Status)
}

type Client struct {
	url      string
	username string
	handlers Handlers

	mu   sync.Mutex
	conn *websocket.Conn

	// out is the outbound queue; Send() pushes here, the write loop drains it
	out chan protocol.Envelope

	ctx    context.Context
	cancel context.CancelFunc
}

// creates a client. url is like "ws://localhost:8080/ws"
func New(url, username string, handlers Handlers) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		url:      url,
		username: username,
		handlers: handlers,
		out:      make(chan protocol.Envelope, 32),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// begins connecting in the background and reconnects on drop until Close
func (c *Client) Start() {
	go c.connectLoop()
}

// tears down the client permanently
func (c *Client) Close() {
	c.cancel()
	c.mu.Lock()
	if c.conn != nil {
		c.conn.Close(websocket.StatusNormalClosure, "client closing")
	}
	c.mu.Unlock()
}

// queues a chat message for delivery. safe to call from the UI thread
// returns immediately, actual sending happens on the write goroutine
func (c *Client) Send(channelID, body string) {
	msg := protocol.ChatMessage{
		ChannelID: channelID,
		Username:  c.username, // server overrides this authoritatively
		Body:      body,
	}

	env, err := protocol.Encode(protocol.KindChat, msg)
	if err != nil {
		return
	}

	select {
	case c.out <- env:
	case <-c.ctx.Done():
	}
}

// dials, runs the session, and retries with backoff on failure
func (c *Client) connectLoop() {
	backoff := time.Second
	for {
		if c.ctx.Err() != nil {
			return
		}
		c.report(StatusConnecting)

		if err := c.session(); err != nil {
			debug.Debug("session ended", "err", err)
		}
		c.report(StatusDisconnected)

		// wait before reconnecting, but bail immediately if we're shutting down
		select {
		case <-c.ctx.Done():
			return
		case <-time.After(backoff):
		}
		if backoff < 15*time.Second {
			backoff *= 2
		}
	}
}

// session runs one full connection lifecycle: dial, hello, then read+write loops
func (c *Client) session() error {
	dialCtx, cancel := context.WithTimeout(c.ctx, 10*time.Second)
	conn, _, err := websocket.Dial(dialCtx, c.url, nil)
	cancel()
	if err != nil {
		return err
	}
	defer conn.CloseNow()

	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()

	// announce ourselves first, the server expects a hello as the first frame
	hello, _ := protocol.Encode(protocol.KindHello, protocol.Hello{Username: c.username})
	hctx, hcancel := context.WithTimeout(c.ctx, 5*time.Second)
	err = wsjson.Write(hctx, conn, hello)
	hcancel()
	if err != nil {
		return err
	}

	c.report(StatusConnected)

	writeDone := make(chan struct{})
	go func() {
		defer close(writeDone)
		for {
			select {
			case <-c.ctx.Done():
				return
			case env := <-c.out:
				wctx, wcancel := context.WithTimeout(c.ctx, 5*time.Second)
				werr := wsjson.Write(wctx, conn, env)
				wcancel()
				if werr != nil {
					return
				}
			}
		}
	}()

	// runs in this goroutine until the connection drops
	readErr := c.readLoop(conn)

	// tear down the write loop by closing the conn (its' Write will error out)
	conn.Close(websocket.StatusNormalClosure, "")
	<-writeDone
	return readErr
}

func (c *Client) readLoop(conn *websocket.Conn) error {
	for {
		var env protocol.Envelope
		if err := wsjson.Read(c.ctx, conn, &env); err != nil {
			return err
		}
		c.dispatch(env)
	}
}

// decodes a frame and fires the matching handler
func (c *Client) dispatch(env protocol.Envelope) {
	switch env.Type {
	case protocol.KindChat:
		var msg protocol.ChatMessage
		if err := json.Unmarshal(env.Data, &msg); err != nil {
			return
		}
		if c.handlers.OnChat != nil {
			c.handlers.OnChat(msg)
		}
	case protocol.KindSystem:
		var sys protocol.System
		if err := json.Unmarshal(env.Data, &sys); err != nil {
			return
		}
		if c.handlers.OnSystem != nil {
			c.handlers.OnSystem(sys)
		}
	}
}

func (c *Client) report(s Status) {
	if c.handlers.OnStatus != nil {
		c.handlers.OnStatus(s)
	}
}
