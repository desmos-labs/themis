package instagram

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterGinHandler registers the proper handlers inside the given gin engine
func RegisterGinHandler(r *gin.Engine, cfg *Config) {
	handler := NewHandler(cfg.CacheFilePath, NewAPI())
	r.Group("/instagram").
		GET("/users/:userID", func(c *gin.Context) {
			user, err := handler.GetUser(c.Param("userID"))
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			c.JSON(http.StatusOK, &user)
		})
}
