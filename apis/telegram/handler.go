package telegram

import (
	"github.com/desmos-labs/themis/apis/hephaestus"
)

type Handler struct {
	*hephaestus.Handler
}

func NewHandler(cfg *Config, hephaestusCfg *hephaestus.Config) *Handler {
	return &Handler{
		Handler: hephaestus.NewHandler(cfg.StoreFolderPath, hephaestusCfg),
	}
}
