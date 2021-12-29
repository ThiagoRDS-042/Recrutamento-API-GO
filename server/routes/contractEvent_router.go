package routes

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/controllers"
	"github.com/gin-gonic/gin"
)

// ContractEventRouterConfig define as configurações das rotas dos historicos dos contratos.
func ContractEventRouterConfig(router *gin.RouterGroup) {
	contractEventController := controllers.NewContractEventController()

	hitorico := router.Group("contrato")
	{
		hitorico.GET("/:id/historico", contractEventController.FindContractEventsByContractID)
	}
}
