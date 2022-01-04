package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SwaggerRouterConfig define as configurações da rota de documentação da API.
func SwaggerRouterConfig(router *gin.RouterGroup) {
	swagger := router.Group("swagger")
	{
		swagger.GET("/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}
