package twitch

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterGinHandler registers the proper handlers inside the given gin engine
func RegisterGinHandler(r *gin.Engine, cfg *Config) {
	twitterHandler := NewHandler(NewAPI(cfg.ClientID, cfg.ClientSecret))
	r.Group("/twitch").
		GET("/users/:username", func(c *gin.Context) {
			user, err := twitterHandler.GetUser(c.Param("username"))
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			c.JSON(http.StatusOK, &user)
		})
}
