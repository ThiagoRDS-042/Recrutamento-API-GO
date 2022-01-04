package routes

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/controllers"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/services"
	"github.com/gin-gonic/gin"
)

// Point
var (
	pointRepository = repositories.NewPointRepository(db)
	pointService    = services.NewPointService(pointRepository, contractService)
)

// PointRouterConfig define as configurações das rotas dos pontos.
func PointRouterConfig(router *gin.RouterGroup) {
	pointController := controllers.NewPointController(pointService, clientService, addressService, contractService)

	points := router.Group("pontos")
	{
		points.POST("/", pointController.CreatePoint)
		points.GET("/", pointController.FindPoints)
	}

	point := router.Group("ponto")
	{
		point.DELETE("/:id", pointController.DeletePoint)
	}
}
