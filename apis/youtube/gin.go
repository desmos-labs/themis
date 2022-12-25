package youtube

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterGinHandler registers the proper handlers inside the given gin engine
func RegisterGinHandler(r *gin.Engine, cfg *Config) {
	handler := NewHandler(NewAPI(cfg.APIKey))
	r.Group("/youtube").
		GET("/users/:user", func(c *gin.Context) {
			channel, err := handler.GetChannel(c.Param("user"))
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			// Return a 404 if the channel is not found
			if channel == nil {
				c.String(http.StatusNotFound, "Channel not found")
				return
			}

			c.JSON(http.StatusOK, &channel)
		})
}
