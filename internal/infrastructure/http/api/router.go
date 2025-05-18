package api

import (
	"github.com/gin-gonic/gin"
)

type RouterRegistry interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type API struct {
	registry []RouterRegistry
}

func NewAPI(registry ...RouterRegistry) *API {
	return &API{
		registry,
	}
}

func (api *API) RegisterRoutes(router *gin.Engine) {
	apiGroup := router.Group("/api/v1")

	for _, registry := range api.registry {
		registry.RegisterRoutes(apiGroup)
	}
}
