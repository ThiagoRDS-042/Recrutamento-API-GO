package routes

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/controllers"
	"github.com/gin-gonic/gin"
)

// ConfigRoutes define as configurações das rotas.
func ConfigRoutes(router *gin.Engine) *gin.Engine {
	clientController := controllers.NewClientController()

	main := router.Group("api/v1")
	{
		clients := main.Group("clientes")
		{
			clients.POST("/", clientController.CreateClient)
			clients.GET("/", clientController.FindClients)
		}

		client := main.Group("cliente")
		{
			client.PUT("/:id", clientController.UpdateClient)
			client.GET("/:id", clientController.FindClientByID)
			client.DELETE("/:id", clientController.DeleteClient)
		}
	}

	return router
}