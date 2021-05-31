package discord

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

// RegisterGinHandler registers the proper handlers inside the given gin engine
func RegisterGinHandler(r *gin.Engine, cfg *Config) {
	handler := NewHandler(cfg)

	r.Group("/discord").
		POST("/data", func(c *gin.Context) {
			err := handleSaveDataReq(c, handler)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			c.Status(http.StatusOK)
		}).
		GET("/:user", func(c *gin.Context) {
			data, err := handler.GetVerificationDataForUser(c.Param("user"))
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			if data == nil {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			c.JSON(http.StatusOK, data)
		})
}

// handleSaveDataReq handles the request that is done when saving some data
func handleSaveDataReq(c *gin.Context, handler *Handler) error {
	jsonBz, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	var data SaveDataReq
	err = json.Unmarshal(jsonBz, &data)
	if err != nil {
		return err
	}

	return handler.HandleSaveData(data)
}
