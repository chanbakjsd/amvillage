package main

import "fmt"

// Config is the configuration for the program.
type Config struct {
	// ListenAddr is the address for the HTTP server to listen.
	ListenAddr string `json:"listen_addr,omitempty"`

	// Currencies is the list of currencies in play.
	Currencies []string `json:"currencies"`
	// Teams is the list of teams participating in the game.
	Teams []TeamConfig `json:"teams"`
}

// TeamConfig is the config for each team.
type TeamConfig struct {
	// Name is the name to be displayed for the team.
	Name string `json:"name"`
	// Password is the password used to join the team.
	Password string `json:"password,omitempty"`
	// IsAdmin allows the team to have negative balances and access to admin
	// tools.
	IsAdmin bool `json:"is_admin,omitempty"`
}

// WithoutSecret returns Config without server-specific configuration.
func (c Config) WithoutSecret() Config {
	teams := make([]TeamConfig, 0, len(c.Teams))
	for _, team := range c.Teams {
		teams = append(teams, TeamConfig{
			Name:    team.Name,
			IsAdmin: team.IsAdmin,
		})
	}
	return Config{
		Currencies: c.Currencies,
		Teams:      teams,
	}
}

// Validate returns any error validating the configuration.
func (c Config) Validate() error {
	password := make(map[string]bool, len(c.Teams))
	hasAdmin := false
	for _, team := range c.Teams {
		if password[team.Password] {
			return fmt.Errorf("password %q is reused", team.Password)
		}
		password[team.Password] = true
		hasAdmin = hasAdmin || team.IsAdmin
	}
	if !hasAdmin {
		return fmt.Errorf("there are no teams with admin permission")
	}
	return nil
}
