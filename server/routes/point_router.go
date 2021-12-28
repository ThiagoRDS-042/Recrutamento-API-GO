package routes

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/controllers"
	"github.com/gin-gonic/gin"
)

// PointRouterConfig define as configurações das rotas dos pontos.
func PointRouterConfig(router *gin.RouterGroup) {
	pointController := controllers.NewPointController()

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
