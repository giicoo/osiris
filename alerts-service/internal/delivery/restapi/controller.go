package restapi

import (
	"github.com/giicoo/osiris/alerts-service/internal/config"
	"github.com/giicoo/osiris/alerts-service/internal/services"
)

type Controller struct {
	cfg      *config.Config
	services *services.Services
}

func NewController(cfg *config.Config, services *services.Services) *Controller {
	return &Controller{
		cfg:      cfg,
		services: services,
	}
}
