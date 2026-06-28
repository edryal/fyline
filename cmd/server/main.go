package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/edryal/fyline/internal/protocol"
	"github.com/google/uuid"
)

const writeTimeout = 5 * time.Second

// client is a connected user from the hub's pov
type client struct {
	username string

	// send is the outbound queue. The hub pushes envelopes here; the client's
	// write goroutine drains it. Buffered so a slow client doesn't block the hub.
	send chan protocol.Envelope
}

// Hub owns the set of clients and serializes all access through channels
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
			log.Printf("client connected: %s (%d online)", c.username, len(h.clients))
			h.announce(c.username + " joined")

		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
				log.Printf("client disconnected: %s (%d online)", c.username, len(h.clients))
				h.announce(c.username + " left")
			}

		case env := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- env:
				default:
					// Client's queue is full — it's too slow. Drop it rather than
					// block the whole hub. The read goroutine will clean up.
					delete(h.clients, c)
					close(c.send)
				}
			}
		}
	}
}

// announce builds a system notice and pushes it onto the broadcast channel.
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
		// For local dev on LAN only.
		// when connecting over the internet - set OriginPatterns to known hosts instead.
		InsecureSkipVerify: true,
	})

	if err != nil {
		log.Printf("accept error: %v", err)
		return
	}
	defer conn.CloseNow()

	ctx := r.Context()

	// First frame must be a Hello so we learn the username.
	var hello protocol.Hello
	if name, ok := readHello(ctx, conn); ok {
		hello.Username = name
	} else {
		conn.Close(websocket.StatusPolicyViolation, "expected hello")
		return
	}

	c := &client{
		username: hello.Username,
		send:     make(chan protocol.Envelope, 16),
	}
	h.register <- c
	defer func() { h.unregister <- c }()

	// Write goroutine: drain c.send to the socket.
	go func() {
		for env := range c.send {
			wctx, cancel := context.WithTimeout(ctx, writeTimeout)
			err := wsjson.Write(wctx, conn, env)
			cancel()
			if err != nil {
				return
			}
		}
	}()

	// Read loop: decode frames and forward chat to the hub.
	for {
		var env protocol.Envelope
		if err := wsjson.Read(ctx, conn, &env); err != nil {
			return // connection closed or errored; defer unregisters us
		}
		h.handleFrame(c, env)
	}
}

// readHello reads exactly one frame and returns the username if it's a valid hello.
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

// handleFrame processes one decoded frame from a client.
func (h *Hub) handleFrame(c *client, env protocol.Envelope) {
	switch env.Type {
	case protocol.KindChat:
		var msg protocol.ChatMessage
		if err := json.Unmarshal(env.Data, &msg); err != nil {
			return
		}

		// Server is authoritative for identity, id, and timestamp — never trust
		// the client's claimed username or time.
		msg.Username = c.username
		msg.ID = uuid.NewString()
		msg.SentAt = time.Now().UTC()

		out, err := protocol.Encode(protocol.KindChat, msg)
		if err != nil {
			return
		}
		h.broadcast <- out

	default:
		// Unknown or unsupported-from-client type; ignore.
	}
}

func main() {
	hub := newHub()
	go hub.run()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.serveWS)

	addr := ":8080"
	log.Printf("Fyline server listening on %s (ws://%s/ws)", addr, addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
