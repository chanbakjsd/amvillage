package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	maxMessageSize = 1024
	pingInterval   = time.Minute
	pongWait       = pingInterval + 5*time.Second

	lockTime = 60000
)

var (
	//go:embed build/*
	frontend embed.FS
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(*http.Request) bool { return true },
	}
)

func main() {
	configFile := "config.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}
	var cfg Config
	if err := readConfig(configFile, &cfg); err != nil {
		log.Fatalln("Error reading config:", err)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatalln("Error validating config:", err)
	}
	frontend, err := fs.Sub(frontend, "build")
	if err != nil {
		log.Fatalln("Error retrieving frontend directory:", err)
	}
	g := NewGameState(cfg)
	go g.Loop()
	http.HandleFunc("/ws", g.handleWebsocket)
	http.Handle("/", http.FileServer(http.FS(frontend)))
	fmt.Println("Listening on", cfg.ListenAddr)
	if err := http.ListenAndServe(cfg.ListenAddr, nil); err != nil {
		log.Println("Error listening and serving:", err)
	}
}

// readConfig reads the configuration or returns an error.
func readConfig(filename string, cfg *Config) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("error decoding config: %w", err)
	}
	return nil
}
