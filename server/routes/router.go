package routes

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/controllers"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/database"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/services"
	"github.com/gin-gonic/gin"
)

// ConfigRoutes define as configurações das rotas.
func ConfigRoutes(router *gin.Engine) *gin.Engine {
	// Database
	db := database.GetDB()

	// Contract
	contractRepository := repositories.NewContractRepository(db)
	contractService := services.NewContractService(contractRepository)

	// Point
	pointRepository := repositories.NewPointRepository(db)
	pointService := services.NewPointService(pointRepository, contractService)

	// Cient
	clientRepository := repositories.NewClientRepository(db)
	clientService := services.NewClientService(clientRepository)

	// Address
	addressRepository := repositories.NewAddressRepository(db)
	addressService := services.NewAddressService(addressRepository, pointService)

	// ContractEvent
	contractEventRepository := repositories.NewContractEventRepository(db)
	contractEventService := services.NewContractEventService(contractEventRepository)

	// Controllers
	clientController := controllers.NewClientController(clientService, pointService)
	addressController := controllers.NewAddressController(addressService)
	pointController := controllers.NewPointController(pointService, clientService, addressService, contractService)
	contractController := controllers.NewContractController(contractService, contractEventService, pointService)
	contractEventController := controllers.NewContractEventController(contractEventService, contractService)

	router.SetTrustedProxies([]string{"192.168.1.2"})
	main := router.Group("api/v1")
	{
		ClientRouterConfig(main, clientController)
		AddressRouterConfig(main, addressController)
		PointRouterConfig(main, pointController)
		ContractRouterConfig(main, contractController)
		ContractEventRouterConfig(main, contractEventController)
	}
	SwaggerRouterConfig(router.Group(""))

	return router
}
