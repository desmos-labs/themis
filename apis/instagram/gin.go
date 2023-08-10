package instagram

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterGinHandler registers the proper handlers inside the given gin engine
func RegisterGinHandler(r *gin.Engine, cfg *Config) {
	handler := NewHandler(cfg.CacheFilePath, NewAPI())
	r.Group("/instagram").
		GET("/medias/:username", func(c *gin.Context) {
			media, err := handler.GetUserMedia(c.Param("username"))
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			c.JSON(http.StatusOK, &media)
		}).
		POST("/medias", func(c *gin.Context) {
			var request AddUserMediaRequest
			err := c.BindJSON(&request)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			err = handler.RequestUserMedia(request.AccessToken)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			c.JSON(http.StatusOK, nil)
		})
}
