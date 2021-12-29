package routes

import (
	"github.com/gin-gonic/gin"
)

// ConfigRoutes define as configurações das rotas.
func ConfigRoutes(router *gin.Engine) *gin.Engine {
	main := router.Group("api/v1")
	{
		ClientRouterConfig(main)
		AddressRouterConfig(main)
		PointRouterConfig(main)
		ContractRouterConfig(main)
	}

	return router
}
