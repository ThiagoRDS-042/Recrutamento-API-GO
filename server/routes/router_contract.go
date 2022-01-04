package routes

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/controllers"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/services"
	"github.com/gin-gonic/gin"
)

// ContractEvent
var (
	contractEventRepository = repositories.NewContractEventRepository(db)
	contractEventService    = services.NewContractEventService(contractEventRepository)
)

// ContractRouterConfig define as configurações das rotas dos contratos.
func ContractRouterConfig(router *gin.RouterGroup) {
	contractController := controllers.NewContractController(contractService, contractEventService, pointService)

	clients := router.Group("contratos")
	{
		clients.POST("/", contractController.CreateContract)
		clients.GET("/", contractController.FindContracts)
	}

	client := router.Group("contrato")
	{
		client.PUT("/:id", contractController.UpdateContract)
		client.GET("/:id", contractController.FindContractByID)
		client.DELETE("/:id", contractController.DeleteContract)
	}
}
