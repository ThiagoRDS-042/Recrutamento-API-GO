package main

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/database"
	_ "github.com/ThiagoRDS-042/Recrutamento-API-GO/docs"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/server"
)

func main() {

	// @title API Recrutamento
	// @version 1.0
	// @description Uma API básica para cadastro de clientes, endereços e contratos.
	// @termsOfService http://swagger.io/terms/

	// @contact.name Thiago Rodrigues
	// @contact.url http://thiagords042/support
	// @contact.email thiagords042@gmail.com

	// @license.name Apache 2.0
	// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

	// @host localhost:2222
	// @BasePath /api/v1

	database.ConnectDB()
	defer database.CloseDB()

	server := server.NewServer()

	server.Run()
}
