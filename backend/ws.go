package main

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

// Conn is a connection to the server.
type Conn struct {
	ws      *websocket.Conn
	closed  bool
	limiter *rate.Limiter
	mu      *sync.Mutex

	team     int
	username string
}

// handleWebsocket handles a new websocket connection.
func (g *GameState) handleWebsocket(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "This endpoint is GET only.", http.StatusMethodNotAllowed)
		return
	}
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %s\n", err)
		return
	}
	conn := &Conn{
		ws:      ws,
		limiter: rate.NewLimiter(rate.Every(time.Second), 5),
		mu:      new(sync.Mutex),
		team:    -1,
	}
	g.Register(conn)
	g.pushState(conn)
	go conn.readPump(g)
	go func() {
		for range time.NewTicker(pingInterval).C {
			conn.mu.Lock()
			if conn.closed {
				return
			}
			conn.ws.WriteMessage(websocket.PingMessage, []byte{})
			conn.mu.Unlock()
		}
	}()
}

// Close closes the underlying connections and channels.
func (conn *Conn) Close() {
	conn.closed = true
	conn.ws.Close()
}

// Write writes the message to the connection.
func (conn *Conn) Write(msg string) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.closed {
		// If the connection is closed, just drop the message.
		return
	}
	err := conn.ws.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("Error writing to connection:", err)
	}
}

// readPump copies the commands provided by the client to the game.
func (conn *Conn) readPump(g *GameState) {
	defer conn.Close()
	conn.ws.SetReadLimit(maxMessageSize)
	conn.ws.SetReadDeadline(time.Now().Add(pongWait))
	conn.ws.SetPongHandler(func(string) error {
		conn.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := conn.ws.ReadMessage()
		switch {
		case err == nil:
			// Process the message below.
		case websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway):
			// Unexpected close.
			log.Printf("Unexpected close: %s\n", err)
			return
		default:
			return
		}
		if conn.limiter.Allow() {
			g.ProcessCommand(conn, strings.Fields(string(message)))
		} else {
			conn.Write("error You have hit a rate limit")
		}
	}
}
