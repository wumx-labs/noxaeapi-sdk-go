// Package ws provides the optional realtime console/event socket for the
// NoxAeApi Go SDK. It is a separate module (with its own go.mod) so that
// the core github.com/NoxlyDev/noxaeapi-sdk-go module stays free of the
// gorilla/websocket dependency for users who don't need it.
//
//	go get github.com/NoxlyDev/noxaeapi-sdk-go/ws
package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

// EventType distinguishes the kind of Event delivered on the channel.
type EventType int

const (
	EventOpen EventType = iota
	EventMessage
	EventConsole
	EventGame
	EventError
	EventClose
)

// Event is a single item delivered from the socket's Events channel.
type Event struct {
	Type    EventType
	Payload json.RawMessage // for EventMessage/EventConsole/EventGame
	Err     error           // for EventError
}

// Options configures a Socket.
type Options struct {
	BaseURL string // e.g. "http://localhost:8080" (same base URL as the REST client)
	APIKey  string
	Route   string // e.g. "console" or "events" — defaults to "events"

	// AutoReconnect enables automatic reconnection with exponential
	// backoff on unexpected disconnects. Default true.
	AutoReconnect *bool
	// MaxReconnectDelay caps the backoff delay. Default 30s.
	MaxReconnectDelay time.Duration
	// Logger receives internal connection diagnostics. Defaults to
	// discarding output; pass log.Default() to see reconnect attempts, etc.
	Logger *log.Logger
}

// Socket is a thin wrapper around the server's WebSocket endpoints
// (console tail and event broadcasts). It runs the connection on a
// background goroutine and delivers events through the Events channel,
// reconnecting automatically (with exponential backoff) on unexpected
// disconnects.
//
//	socket, err := ws.Connect(ws.Options{BaseURL: client.BaseURL(), APIKey: client.APIKey(), Route: "console"})
//	for event := range socket.Events {
//	    switch event.Type {
//	    case ws.EventConsole:
//	        fmt.Println(string(event.Payload))
//	    }
//	}
type Socket struct {
	Events chan Event

	url      string
	closed   atomic.Bool
	sendCh   chan []byte
	wg       sync.WaitGroup
	autoRec  bool
	maxDelay time.Duration
	logger   *log.Logger
}

// Connect opens a WebSocket connection and starts the background read loop.
func Connect(opts Options) (*Socket, error) {
	route := opts.Route
	if route == "" {
		route = "events"
	}
	wsURL, err := toWSURL(opts.BaseURL, route, opts.APIKey)
	if err != nil {
		return nil, err
	}

	autoRec := true
	if opts.AutoReconnect != nil {
		autoRec = *opts.AutoReconnect
	}
	maxDelay := opts.MaxReconnectDelay
	if maxDelay <= 0 {
		maxDelay = 30 * time.Second
	}
	logger := opts.Logger
	if logger == nil {
		logger = log.New(discardWriter{}, "", 0)
	}

	s := &Socket{
		Events:   make(chan Event, 32),
		url:      wsURL,
		sendCh:   make(chan []byte, 16),
		autoRec:  autoRec,
		maxDelay: maxDelay,
		logger:   logger,
	}

	s.wg.Add(1)
	go s.run()

	return s, nil
}

func toWSURL(baseURL, route, apiKey string) (string, error) {
	trimmed := strings.TrimRight(baseURL, "/")
	route = strings.Trim(route, "/")
	raw := fmt.Sprintf("%s/v1/ws/%s", trimmed, route)
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}
	if apiKey != "" {
		q := u.Query()
		q.Set("key", apiKey)
		u.RawQuery = q.Encode()
	}
	return u.String(), nil
}

func (s *Socket) run() {
	defer s.wg.Done()
	attempt := 0

	for {
		if s.closed.Load() {
			return
		}

		conn, _, err := websocket.DefaultDialer.Dial(s.url, nil)
		if err != nil {
			s.Events <- Event{Type: EventError, Err: err}
		} else {
			attempt = 0
			s.Events <- Event{Type: EventOpen}
			s.readLoop(conn)
		}

		if s.closed.Load() {
			return
		}

		delay := backoffDelay(attempt, s.maxDelay)
		attempt++
		s.logger.Printf("noxaeapi/ws: reconnecting in %s (attempt %d)", delay, attempt)
		time.Sleep(delay)

		if !s.autoRec {
			return
		}
	}
}

func (s *Socket) readLoop(conn *websocket.Conn) {
	defer conn.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case msg := <-s.sendCh:
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					return
				}
			case <-done:
				return
			}
		}
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			s.Events <- Event{Type: EventError, Err: err}
			s.Events <- Event{Type: EventClose}
			return
		}

		var parsed struct {
			Type string `json:"type"`
		}
		_ = json.Unmarshal(data, &parsed)

		switch parsed.Type {
		case "console":
			s.Events <- Event{Type: EventConsole, Payload: data}
		case "event":
			s.Events <- Event{Type: EventGame, Payload: data}
		}
		s.Events <- Event{Type: EventMessage, Payload: data}
	}
}

// SendText sends a raw text payload over the socket.
func (s *Socket) SendText(text string) {
	if s.closed.Load() {
		return
	}
	select {
	case s.sendCh <- []byte(text):
	default:
	}
}

// SendJSON marshals v and sends it as a text payload.
func (s *Socket) SendJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	s.SendText(string(data))
	return nil
}

// Close closes the socket and stops reconnecting.
func (s *Socket) Close() {
	s.closed.Store(true)
}

func backoffDelay(attempt int, max time.Duration) time.Duration {
	base := 500 * time.Millisecond
	exp := float64(base) * math.Pow(2, float64(attempt))
	if exp > float64(max) {
		exp = float64(max)
	}
	return time.Duration(rand.Float64() * exp)
}

type discardWriter struct{}

func (discardWriter) Write(p []byte) (int, error) { return len(p), nil }
