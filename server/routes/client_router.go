package routes

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/controllers"
	"github.com/gin-gonic/gin"
)

// ClientRouterConfig define as configurações das rotas dos clientes.
func ClientRouterConfig(router *gin.RouterGroup) {
	clientController := controllers.NewClientController()

	clients := router.Group("clientes")
	{
		clients.POST("/", clientController.CreateClient)
		clients.GET("/", clientController.FindClients)
	}

	client := router.Group("cliente")
	{
		client.PUT("/:id", clientController.UpdateClient)
		client.GET("/:id", clientController.FindClientByID)
		client.DELETE("/:id", clientController.DeleteClient)
	}
}
