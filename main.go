package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/desmos-labs/themis/twitter"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// config contains the data that should be present inside the configuration file
type config struct {
	Twitter *twitter.Config
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

	// Build handlers
	twitterHandler := twitter.NewHandler(cfg.Twitter.CacheFilePath, twitter.NewApi(cfg.Twitter.Bearer))

	// Setup the rest server
	r := gin.Default()
	r.Use(gin.Recovery()) // Handles all panics writing 500

	r.Group("/twitter").
		GET("/tweets/:id", func(c *gin.Context) {
			tweet, err := twitterHandler.GetTweet(c.Param("id"))
			if err != nil {
				panic(err)
			}

			c.JSON(http.StatusOK, &tweet)
		}).
		GET("/users/:username", func(c *gin.Context) {
			user, err := twitterHandler.GetUser(c.Param("username"))
			if err != nil {
				panic(err)
			}

			c.JSON(http.StatusOK, &user)
		})

	// Run the server
	r.Run()
}
