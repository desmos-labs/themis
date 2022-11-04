package telegram

import (
	"github.com/gin-gonic/gin"

	"github.com/desmos-labs/themis/apis/hephaestus"
)

// RegisterGinHandler registers the proper handlers inside the given gin engine
func RegisterGinHandler(r *gin.Engine, hephaestsusCfg *hephaestus.Config, cfg *Config) {
	handler := NewHandler(cfg, hephaestsusCfg)
	r.Group("/telegram").
		POST("/data", hephaestus.NewSaveDataRequestGinHandler(handler.Handler)).
		GET("/:user", hephaestus.NewGetVerificationDataGinHandler(handler.Handler))
}
