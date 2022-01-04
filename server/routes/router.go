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

	// Repositories
	clientRepository := repositories.NewClientRepository(db)
	addressRepository := repositories.NewAddressRepository(db)
	pointRepository := repositories.NewPointRepository(db)
	contractRepository := repositories.NewContractRepository(db)
	contractEventRepository := repositories.NewContractEventRepository(db)

	// Services
	clientService := services.NewClientService(clientRepository)
	addressService := services.NewAddressService(addressRepository)
	contractEventService := services.NewContractEventService(contractEventRepository)
	contractService := services.NewContractService(contractRepository, pointRepository, contractEventService)
	pointService := services.NewPointService(pointRepository, clientRepository, addressRepository, contractService)

	// Controllers
	clientController := controllers.NewClientController(clientService, pointService)
	addressController := controllers.NewAddressController(addressService, pointService)
	pointController := controllers.NewPointController(pointService, contractService)
	contractController := controllers.NewContractController(contractService)
	contractEventController := controllers.NewContractEventController(contractEventService)

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
