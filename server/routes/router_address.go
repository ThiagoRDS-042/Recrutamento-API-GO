package routes

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/controllers"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/services"
	"github.com/gin-gonic/gin"
)

// Address
var (
	addressRepository = repositories.NewAddressRepository(db)
	addressService    = services.NewAddressService(addressRepository)
)

// AddressRouterConfig define as configurações das rotas dos endereços.
func AddressRouterConfig(router *gin.RouterGroup) {
	addressController := controllers.NewAddressController(addressService, pointService)

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
