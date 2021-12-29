package routes

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/controllers"
	"github.com/gin-gonic/gin"
)

// ContractRouterConfig define as configurações das rotas dos contratos.
func ContractRouterConfig(router *gin.RouterGroup) {
	contractController := controllers.NewContractController()

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
