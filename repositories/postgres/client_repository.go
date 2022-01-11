package repositories

import (
	"fmt"
	"log"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"gorm.io/gorm"
)

// ClientRepository representa o contracto de ClientRepository.
type ClientRepository interface {
	CreateClient(client entities.Cliente) (entities.Cliente, error)
	UpdateClient(client entities.Cliente) (entities.Cliente, error)
	FindClientByID(clientID string) entities.Cliente
	FindClientByName(name string) entities.Cliente
	DeleteClient(client entities.Cliente) error
	FindClients(clientName string, clientType entities.ClientType) []entities.Cliente
}

type clientConnection struct {
	connection *gorm.DB
}

func (db *clientConnection) CreateClient(client entities.Cliente) (entities.Cliente, error) {
	err := db.connection.Create(&client).Error
	if err != nil {
		return client, err
	}

	return client, nil
}

func (db *clientConnection) UpdateClient(client entities.Cliente) (entities.Cliente, error) {
	err := db.connection.Save(&client).Error
	if err != nil {
		return client, err
	}

	return client, nil
}

func (db *clientConnection) FindClientByID(clientID string) entities.Cliente {
	client := entities.Cliente{}

	err := db.connection.First(&client, "id = ?", clientID).Error
	if err != nil {
		log.Println(err.Error())
	}

	return client
}

func (db *clientConnection) FindClientByName(name string) entities.Cliente {
	client := entities.Cliente{}

	err := db.connection.Unscoped().First(&client, "nome = ?", name).Error
	if err != nil {
		log.Println(err.Error())
	}

	return client
}

func (db *clientConnection) DeleteClient(client entities.Cliente) error {

	err := db.connection.Delete(&client).Error
	if err != nil {
		return err
	}

	return nil
}

func (db *clientConnection) FindClients(clientName string, clientType entities.ClientType) []entities.Cliente {
	clients := []entities.Cliente{}

	var sqlQuery string

	if clientName != "" {
		clientName = fmt.Sprint("%", clientName, "%")
		sqlQuery += fmt.Sprintf("nome LIKE '%v' ", clientName)
	} else {
		sqlQuery += "NOT nome IS NULL "
	}

	if clientType != entities.ClientType("") {
		sqlQuery += fmt.Sprintf("AND tipo = '%v'", clientType)
	} else {
		sqlQuery += "AND NOT tipo IS NULL"
	}

	err := db.connection.Find(&clients, sqlQuery).Error
	if err != nil {
		log.Println(err.Error())
	}

	return clients
}

// NewClientRepository cria uma nova instancia de ClientRepository.
func NewClientRepository(database *gorm.DB) ClientRepository {
	return &clientConnection{
		connection: database,
	}
}
