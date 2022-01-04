package routes

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/controllers"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/services"
	"github.com/gin-gonic/gin"
)

// Contract
var (
	contractRepository = repositories.NewContractRepository(db)
	contractService    = services.NewContractService(contractRepository)
)

// ContractEventRouterConfig define as configurações das rotas dos historicos dos contratos.
func ContractEventRouterConfig(router *gin.RouterGroup) {
	contractEventController := controllers.NewContractEventController(contractEventService, contractService)

	hitorico := router.Group("contrato")
	{
		hitorico.GET("/:id/historico", contractEventController.FindContractEventsByContractID)
	}
}
