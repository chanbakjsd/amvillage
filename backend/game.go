package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// NewGameState returns a new game state based on the provided config.
func NewGameState(cfg Config) *GameState {
	balances := make([]Balance, len(cfg.Teams))
	for i := range balances {
		balances[i] = make([]int, len(cfg.Currencies))
	}
	return &GameState{
		game: Game{
			Config:   cfg.WithoutSecret(),
			Balances: balances,
			Locks:    make([]*Lock, len(cfg.Teams)),
		},
		config: cfg,
		conn:   make(map[*Conn]struct{}),
		update: false,
		mu:     new(sync.Mutex),
	}
}

// GameState is the state of the game server.
type GameState struct {
	game Game
	// config is the configuration of the game.
	config Config
	// conn is a hashset of existing connections.
	conn map[*Conn]struct{}
	// update is set to true when there are updates to be pushed to the server.
	update bool

	mu *sync.Mutex
}

// Game is the state of the active game.
type Game struct {
	// Config is the configuration of the game.
	Config `json:"config"`
	// Balances is the balance of the groups.
	Balances []Balance `json:"balances"`
	// Locks is a list of Lock that are currently held by members of a group as
	// only one member of a group may trade at a time.
	Locks []*Lock `json:"locks"`
}

// Balance is the balance of a group.
type Balance []int

// Lock is a lock held by the member of a group.
type Lock struct {
	// Member is the name of the member holding the lock.
	Member string `json:"member"`
	// MillisLeft is the number of milliseconds before the lock expires.
	MillisLeft int `json:"millis_left"`
}

// Loop is the game loop responsible for pushing updates and removing expired connections.
func (g *GameState) Loop() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		g.tick(time.Second)
	}
}

// tick advances the game by the specified duration.
func (g *GameState) tick(delta time.Duration) {
	g.mu.Lock()
	defer g.mu.Unlock()
	toClose := make([]*Conn, 0)
	for conn := range g.conn {
		if conn.closed {
			toClose = append(toClose, conn)
		}
	}
	for _, conn := range toClose {
		delete(g.conn, conn)
	}
	for i, lock := range g.game.Locks {
		if lock == nil {
			continue
		}
		lock.MillisLeft -= int(delta.Milliseconds())
		if lock.MillisLeft <= 0 {
			g.game.Locks[i] = nil
			g.update = true
		}
	}
	if !g.update {
		return
	}
	g.update = false
	for conn := range g.conn {
		g.pushState(conn)
	}
}

// Register registers the connection to the game.
func (g *GameState) Register(c *Conn) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.conn[c] = struct{}{}
}

// ProcessCommand processes the command provided by the client.
func (g *GameState) ProcessCommand(c *Conn, args []string) {
	if len(args) == 0 {
		go c.Write("error Empty commands not allowed")
		return
	}
	switch args[0] {
	case "login":
	default:
		if c.team == -1 {
			go c.Write("error Login is required before proceeding")
			return
		}
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	switch args[0] {
	case "login":
		if len(args) <= 2 {
			go c.Write("error Expected username and password")
			return
		}
		username := strings.Join(args[2:], " ")
		password, err := base64.StdEncoding.DecodeString(args[1])
		if err != nil {
			go c.Write("error Invalid password")
			return
		}
		g.login(c, username, string(password))
		data, err := json.Marshal(g.game)
		if err != nil {
			log.Fatalln("Error marshalling game state:", err)
		}
		msg := "state " + string(data)
		go c.Write(msg)
	case "lock":
		if len(args) != 1 {
			go c.Write("error Unexpected argument")
			return
		}
		switch {
		case g.game.Locks[c.team] == nil:
			// Successfully acquire lock.
			g.update = true
			g.game.Locks[c.team] = &Lock{
				Member:     c.username,
				MillisLeft: lockTime,
			}
		case g.game.Locks[c.team].Member == c.username:
			// Desync.
			g.pushState(c)
			return
		default:
			// Trying to acquire teammate's lock.
			go c.Write("error Another player from your team is trading")
			return
		}
	case "cancel":
		if len(args) != 1 {
			go c.Write("error Unexpected argument")
			return
		}
		g.update = true
		g.releaseLock(c)
	case "trade":
		if len(args) != 2+len(g.game.Currencies) {
			go c.Write("error Unexpected number of arguments")
			return
		}
		nums := make([]int, 0, len(args)-1)
		for _, v := range args[1:] {
			num, err := strconv.Atoi(v)
			if err != nil {
				go c.Write("error Expected number for arguments")
				return
			}
			nums = append(nums, num)
		}
		g.trade(c, nums[0], nums[1:])
	default:
		go c.Write("error Unknown command " + args[0])
	}
}

// login attempts to log the player into a team.
func (g *GameState) login(c *Conn, username, password string) {
	team := -1
	for i, teamConfig := range g.config.Teams {
		if teamConfig.Password == password {
			team = i
			break
		}
	}
	if team == -1 {
		go c.Write("error Invalid password")
		return
	}
	c.username = username
	c.team = team
}

// releaseLock releases the lock if it is owned by the provided connection. It
// returns if the lock was released.
func (g *GameState) releaseLock(c *Conn) bool {
	if g.game.Locks[c.team] == nil {
		go c.Write("error You must acquire the trade lock first")
		return false
	}
	if g.game.Locks[c.team].Member != c.username {
		go c.Write("error Another player from your team is trading")
		return false
	}
	g.game.Locks[c.team] = nil
	return true
}

// trade trades the resources described to the target team.
func (g *GameState) trade(c *Conn, to int, amount []int) {
	// Release the lock.
	from := c.team
	if !g.game.Teams[from].IsAdmin && !g.releaseLock(c) {
		return
	}
	if to < 0 || to >= len(g.game.Teams) {
		go c.Write("error Invalid target team ID")
		return
	}
	// Non-admins' inputs are validated. Admins may cause balances to go
	// negative.
	if !g.game.Teams[from].IsAdmin {
		for i, v := range amount {
			if v < 0 {
				go c.Write("error Cannot send negative amount")
				return
			}
			if v > g.game.Balances[from][i] {
				go c.Write("error Cannot send more than you have")
				return
			}
		}
	}
	// Execute the trade.
	g.update = true
	for i, v := range amount {
		g.game.Balances[from][i] -= v
		g.game.Balances[to][i] += v
	}
}

// pushState pushes the state to the specified connection.
func (g *GameState) pushState(c *Conn) {
	data, err := json.Marshal(g.game)
	if err != nil {
		log.Fatalln("Error marshalling game state:", err)
	}
	msg := "state " + string(data)
	go c.Write(msg)
}
