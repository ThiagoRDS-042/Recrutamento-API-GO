package routes

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/controllers"
	"github.com/gin-gonic/gin"
)

// AddressRouterConfig define as configurações das rotas dos endereços.
func AddressRouterConfig(router *gin.RouterGroup, addressController controllers.AddressController) {
	addresses := router.Group("enderecos")
	{
		addresses.POST("/", addressController.CreateAddress)
		addresses.GET("/", addressController.FindAddress)
	}

	address := router.Group("endereco")
	{
		address.PUT("/:id", addressController.UpdateAddress)
		address.GET("/:id", addressController.FindAddressByID)
		address.DELETE("/:id", addressController.DeleteAddress)
	}
}
