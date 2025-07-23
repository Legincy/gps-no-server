package interfaces

import "github.com/gin-gonic/gin"

type StationHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type ClusterHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type RangingHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type MQTTStationHandler interface {
	HandleStationData(payload []byte) error
}

type MQTTRangingHandler interface {
	HandleRangingData(payload []byte) error
}
