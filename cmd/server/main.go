package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/edryal/fyline/internal/debug"
	"github.com/edryal/fyline/internal/protocol"
	"github.com/google/uuid"
)

const writeTimeout = 5 * time.Second

// a connected user from the hub's pov
type client struct {
	username string
	// outbound queue. the hub pushes here, the write goroutine drains it.
	// buffered so a slow client doesn't block the hub
	send chan protocol.Envelope
}

// owns the set of clients, only this goroutine touches the map (no locks)
type Hub struct {
	register   chan *client
	unregister chan *client
	broadcast  chan protocol.Envelope
	clients    map[*client]struct{}
}

func newHub() *Hub {
	return &Hub{
		register:   make(chan *client),
		unregister: make(chan *client),
		broadcast:  make(chan protocol.Envelope, 64),
		clients:    make(map[*client]struct{}),
	}
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = struct{}{}
			debug.Info("client connected", "user", c.username, "online", len(h.clients))
			h.announce(c.username + " joined")

		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
				debug.Info("client disconnected", "user", c.username, "online", len(h.clients))
				h.announce(c.username + " left")
			}

		case env := <-h.broadcast:
			debug.Debug("broadcasting", "type", env.Type, "recipients", len(h.clients))
			for c := range h.clients {
				select {
				case c.send <- env:
				default:
					// queue full, client is too slow. drop it instead of blocking everyone
					debug.Warn("dropping slow client", "user", c.username)
					delete(h.clients, c)
					close(c.send)
				}
			}
		}
	}
}

func (h *Hub) announce(text string) {
	envelop, err := protocol.Encode(protocol.KindSystem, protocol.System{Text: text})
	if err != nil {
		return
	}

	select {
	case h.broadcast <- envelop:
	default:
	}
}

func (h *Hub) serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// LAN dev only. set OriginPatterns to known hosts before exposing over the internet
		InsecureSkipVerify: true,
	})

	if err != nil {
		debug.Warn("accept error", "err", err)
		return
	}
	defer conn.CloseNow()

	ctx := r.Context()

	// first frame must be a hello so we learn the username
	var hello protocol.Hello
	if name, ok := readHello(ctx, conn); ok {
		hello.Username = name
	} else {
		debug.Warn("handshake rejected", "reason", "expected hello first")
		conn.Close(websocket.StatusPolicyViolation, "expected hello")
		return
	}

	c := &client{
		username: hello.Username,
		send:     make(chan protocol.Envelope, 16),
	}
	h.register <- c
	defer func() { h.unregister <- c }()

	// write goroutine: drain c.send to the socket
	go func() {
		for env := range c.send {
			wctx, cancel := context.WithTimeout(ctx, writeTimeout)
			err := wsjson.Write(wctx, conn, env)
			cancel()
			if err != nil {
				debug.Debug("write failed, closing writer", "user", c.username, "err", err)
				return
			}
		}
	}()

	// read loop: decode frames and forward to the hub. defer unregisters us on exit
	for {
		var env protocol.Envelope
		if err := wsjson.Read(ctx, conn, &env); err != nil {
			debug.Debug("read loop ended", "user", c.username, "err", err)
			return
		}
		h.handleFrame(c, env)
	}
}

func readHello(ctx context.Context, conn *websocket.Conn) (string, bool) {
	hctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var env protocol.Envelope
	if err := wsjson.Read(hctx, conn, &env); err != nil {
		return "", false
	}

	if env.Type != protocol.KindHello {
		return "", false
	}

	var hello protocol.Hello
	if err := json.Unmarshal(env.Data, &hello); err != nil || hello.Username == "" {
		return "", false
	}

	return hello.Username, true
}

func (h *Hub) handleFrame(c *client, env protocol.Envelope) {
	switch env.Type {
	case protocol.KindChat:
		var msg protocol.ChatMessage
		if err := json.Unmarshal(env.Data, &msg); err != nil {
			debug.Warn("bad chat frame", "user", c.username, "err", err)
			return
		}

		// server is authoritative, never trust the client's claimed identity or time
		msg.Username = c.username
		msg.ID = uuid.NewString()
		msg.SentAt = time.Now().UTC()

		debug.Debug("chat received", "user", c.username, "channelID", msg.ChannelID, "body", msg.Body)

		out, err := protocol.Encode(protocol.KindChat, msg)
		if err != nil {
			return
		}
		h.broadcast <- out

	default:
		debug.Debug("ignored frame", "type", env.Type, "user", c.username)
	}
}

func main() {
	hub := newHub()
	go hub.run()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.serveWS)

	addr := ":8080"
	if env := os.Getenv("FYLINE_ADDR"); env != "" {
		addr = env
	}

	debug.Info("server listening", "addr", addr, "url", "ws://localhost"+addr+"/ws")
	if err := http.ListenAndServe(addr, mux); err != nil {
		debug.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
