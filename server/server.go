package server

import (
	"log"
	"os"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/server/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Server representa o contrato de servidor.
type Server interface {
	Run()
}

type server struct {
	port   string
	server *gin.Engine
}

// Run inicia o servidor.
func (server *server) Run() {
	router := routes.ConfigRoutes(server.server)

	log.Println("Server is running at port:", server.port)
	router.Run(":" + server.port)
}

// NewServer cria um novo servidor.
func NewServer() Server {
	godotenv.Load()
	port := os.Getenv("SERVER_PORT")

	return &server{
		port:   port,
		server: gin.Default(),
	}
}
