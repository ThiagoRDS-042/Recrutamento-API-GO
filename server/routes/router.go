package routes

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/controllers"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/database"
	repositories "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/postgres"
	addressService "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/address_service"
	clientService "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/client_service"
	contractEventService "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/contract_event_service"
	contractService "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/contract_service"
	pointService "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/point_service"
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
	contractEventService := contractEventService.NewContractEventService(contractEventRepository, contractRepository)
	contractService := contractService.NewContractService(contractRepository, pointRepository, contractEventService)
	pointService := pointService.NewPointService(pointRepository, clientRepository, addressRepository, contractService)
	clientService := clientService.NewClientService(clientRepository, pointService)
	addressService := addressService.NewAddressService(addressRepository, pointService)

	// Controllers
	clientController := controllers.NewClientController(clientService)
	addressController := controllers.NewAddressController(addressService)
	pointController := controllers.NewPointController(pointService)
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
