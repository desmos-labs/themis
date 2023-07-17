package instagram

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterGinHandler registers the proper handlers inside the given gin engine
func RegisterGinHandler(r *gin.Engine, cfg *Config) {
	handler := NewHandler(cfg.CacheFilePath, NewAPI())
	r.Group("/instagram").
		GET("/users/:username", func(c *gin.Context) {
			user, err := handler.GetUser(c.Param("username"))
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			c.JSON(http.StatusOK, &user)
		}).
		POST("/users", func(c *gin.Context) {
			var request AddUserRequest
			err := c.BindJSON(&request)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			err = handler.RequestUser(request.AccessToken)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			c.JSON(http.StatusOK, nil)
		})
}
