package controllers

import (
	"github.com/gin-gonic/gin"
)

type RouteRegistry interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type API struct {
	registry []RouteRegistry
}

func NewAPI(registry ...RouteRegistry) *API {
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
