package routes

import (
	"github.com/gin-gonic/gin"
)

// ConfigRoutes define as configurações das rotas.
func ConfigRoutes(router *gin.Engine) *gin.Engine {
	router.SetTrustedProxies([]string{"192.168.1.2"})
	main := router.Group("api/v1")
	{
		ClientRouterConfig(main)
		AddressRouterConfig(main)
		PointRouterConfig(main)
		ContractRouterConfig(main)
		ContractEventRouterConfig(main)
		SwaggerRouterConfig(main)
	}

	return router
}
