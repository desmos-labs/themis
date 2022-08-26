package youtube

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterGinHandler registers the proper handlers inside the given gin engine
func RegisterGinHandler(r *gin.Engine, cfg *Config) {
	handler := NewHandler(NewAPI(cfg.APIKEY), cfg.CacheFilePath)
	r.Group("/youtube").
		GET("/users/:user_id", func(c *gin.Context) {
			user, err := handler.GetChannel(c.Param("user_id"))
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusOK, &user)
		})
}
