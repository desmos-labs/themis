package nslookup

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterGinHandler registers the proper handlers inside the given gin engine
func RegisterGinHandler(r *gin.Engine) {
	handler := NewHandler()

	r.GET("/nslookup/:domain", func(c *gin.Context) {
		txtRecords, err := handler.ReadTxtRecords(c.Param("domain"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, NewLookupResponse(txtRecords))
	})
}
