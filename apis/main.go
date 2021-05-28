package main

import (
	"fmt"
	"github.com/desmos-labs/themis/discord"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"

	"github.com/desmos-labs/themis/twitter"
)

// config contains the .data that should be present inside the configuration file
type config struct {
	Apis struct {
		Port uint
	}

	Twitter *twitter.Config
	Discord *discord.Config
}

// readConfig parses the file present at the given path and returns a config object
func readConfig(path string) (*config, error) {
	var cfg config
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	if len(os.Args) < 2 {
		panic(fmt.Errorf("missing config argument"))
	}

	// Read the config
	cfg, err := readConfig(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Setup the rest server
	r := gin.Default()
	r.Use(gin.Recovery()) // Handles all panics writing 500

	// Register the handlers
	twitter.RegisterGinHandler(r, cfg.Twitter)
	discord.RegisterGinHandler(r, cfg.Discord)

	// Run the server
	port := cfg.Apis.Port
	if port == 0 {
		port = 8080
	}
	r.Run(fmt.Sprintf(":%d", port))
}
