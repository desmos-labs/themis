package telegram

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/desmos-labs/themis/apis/utils/bot"
	"github.com/gin-gonic/gin"
)

// RegisterGinHandler registers the proper handlers inside the given gin engine
func RegisterGinHandler(r *gin.Engine, cfg *bot.Config) {
	handler := bot.NewHandler(cfg)

	r.Group("/telegram").
		POST("/data", func(c *gin.Context) {
			err := handleSaveDataReq(c, handler)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			c.Status(http.StatusCreated)
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
func handleSaveDataReq(c *gin.Context, handler *bot.Handler) error {
	jsonBz, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	var data bot.SaveDataReq
	err = json.Unmarshal(jsonBz, &data)
	if err != nil {
		return err
	}

	return handler.HandleSaveData(data)
}
