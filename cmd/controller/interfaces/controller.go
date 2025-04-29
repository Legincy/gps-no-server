package interfaces

import "github.com/gin-gonic/gin"

type Controller interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type API struct {
	Controllers []Controller
}

func NewAPI(controllers ...Controller) *API {
	return &API{
		controllers,
	}
}

func (api *API) RegisterRoutes(router *gin.RouterGroup) {
	apiGroup := router.Group("/api/v1")

	for _, controller := range api.Controllers {
		controller.RegisterRoutes(apiGroup)
	}
}
