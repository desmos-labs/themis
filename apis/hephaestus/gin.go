package hephaestus

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewSaveDataRequestGinHandler returns a new function that can be used as a Gin handler
// to handle verification data saving requests
func NewSaveDataRequestGinHandler(handler *Handler) func(c *gin.Context) {
	return func(c *gin.Context) {
		req, err := handler.ParseSaveDataRequest(c.Request)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = handler.HandleSaveData(req)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

// NewGetVerificationDataGinHandler returns a new function that can be used as a Gin handler
// to handle verification data retrieval requests
func NewGetVerificationDataGinHandler(handler *Handler) func(c *gin.Context) {
	return func(c *gin.Context) {
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
	}
}
