package routes

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/controllers"
	"github.com/gin-gonic/gin"
)

// ConfigRoutes define as configurações das rotas.
func ConfigRoutes(router *gin.Engine) *gin.Engine {
	clientController := controllers.NewClientController()
	addressController := controllers.NewAddressController()

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

		addresses := main.Group("enderecos")
		{
			addresses.POST("/", addressController.CreateAddress)
			addresses.GET("/", addressController.FindAddress)
		}

		address := main.Group("endereco")
		{
			address.PUT("/:id", addressController.UpdateAddress)
			address.GET("/:id", addressController.FindAddressByID)
			address.DELETE("/:id", addressController.DeleteAddress)
		}
	}

	return router
}
