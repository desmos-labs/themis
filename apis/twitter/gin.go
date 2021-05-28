package twitter

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterGinHandler registers the proper handlers inside the given gin engine
func RegisterGinHandler(r *gin.Engine, cfg *Config) {
	twitterHandler := NewHandler(cfg.CacheFilePath, NewAPI(cfg.Bearer))
	r.Group("/twitter").
		GET("/tweets/:id", func(c *gin.Context) {
			tweet, err := twitterHandler.GetTweet(c.Param("id"))
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			c.JSON(http.StatusOK, &tweet)
		}).
		GET("/users/:username", func(c *gin.Context) {
			user, err := twitterHandler.GetUser(c.Param("username"))
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			c.JSON(http.StatusOK, &user)
		})
}
