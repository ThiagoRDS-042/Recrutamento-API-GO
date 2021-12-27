package main

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/database"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/server"
)

func main() {
	database.ConnectDB()
	defer database.CloseDB()

	server := server.NewServer()

	server.Run()
}
