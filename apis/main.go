package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/themis/apis/discord"
	"github.com/desmos-labs/themis/apis/hephaestus"
	"github.com/desmos-labs/themis/apis/nslookup"
	"github.com/desmos-labs/themis/apis/telegram"
	"github.com/desmos-labs/themis/apis/twitch"
	"github.com/desmos-labs/themis/apis/twitter"
	"github.com/desmos-labs/themis/apis/youtube"
)

const (
	ConfigPathEnv = "CONFIG_PATH"
)

// config contains the data that should be present inside the configuration file
type config struct {
	Apis struct {
		Port     uint   `yaml:"port" toml:"port"`
		Address  string `yaml:"address" toml:"address"`
		LogLevel string `yaml:"log_level" toml:"log_level"`
	} `yaml:"apis" toml:"apis"`

	Twitter    *twitter.Config    `yaml:"twitter" toml:"twitter"`
	Discord    *discord.Config    `yaml:"discord" toml:"discord"`
	Telegram   *telegram.Config   `yaml:"telegram" toml:"telegram"`
	Twitch     *twitch.Config     `yaml:"twitch" toml:"twitch"`
	Hephaestus *hephaestus.Config `yaml:"hephaestus" toml:"hephaestus"`
	Youtube    *youtube.Config    `yaml:"youtube" toml:"youtube"`
}

func getConfigPath() (string, error) {
	// Try reading from env variable
	if configPath := os.Getenv(ConfigPathEnv); configPath != "" {
		return configPath, nil
	}

	if len(os.Args) < 2 {
		return "", fmt.Errorf("no config path found: use either the env variable %s or the command argument", ConfigPathEnv)
	}

	return os.Args[1], nil
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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	cfgPath, err := getConfigPath()
	if err != nil {
		panic(err)
	}

	// Read the config
	log.Debug().Msg("reading config")
	cfg, err := readConfig(cfgPath)
	if err != nil {
		panic(err)
	}

	// Setup the rest server
	r := gin.Default()
	r.Use(gin.Recovery()) // Handles all panics writing 500
	r.Use(cors.Default()) // Allows all origins

	// Register the handlers
	log.Debug().Msg("registering handlers")
	twitter.RegisterGinHandler(r, cfg.Twitter)
	discord.RegisterGinHandler(r, cfg.Hephaestus, cfg.Discord)
	twitch.RegisterGinHandler(r, cfg.Twitch)
	nslookup.RegisterGinHandler(r)
	telegram.RegisterGinHandler(r, cfg.Hephaestus, cfg.Telegram)
	youtube.RegisterGinHandler(r, cfg.Youtube)

	// Run the server
	port := cfg.Apis.Port
	if port == 0 {
		port = 8080
	}
	address := fmt.Sprintf("%s:%d", cfg.Apis.Address, port)
	log.Info().Msgf("running on address %s", address)
	r.Run(address)
}
