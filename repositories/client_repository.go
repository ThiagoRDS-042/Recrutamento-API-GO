package repositories

import (
	"fmt"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/database"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"gorm.io/gorm"
)

// ClientRepository representa o contracto de ClientRepository.
type ClientRepository interface {
	CreateClient(client entities.Cliente) (entities.Cliente, error)
	UpdateClient(client entities.Cliente) (entities.Cliente, error)
	FindClientByID(clientID string) (entities.Cliente, error)
	FindClientByName(name string) (entities.Cliente, error)
	DeleteClient(client entities.Cliente) error
	FindClients(clientName string, clientType string) ([]entities.Cliente, error)
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

func (db *clientConnection) FindClientByID(clientID string) (entities.Cliente, error) {
	client := entities.Cliente{}

	err := db.connection.First(&client, "id = ?", clientID).Error
	if err != nil {
		return client, err
	}

	return client, nil
}

func (db *clientConnection) FindClientByName(name string) (entities.Cliente, error) {
	client := entities.Cliente{}

	err := db.connection.Unscoped().First(&client, "nome = ?", name).Error
	if err != nil {
		return client, err
	}

	return client, nil
}

func (db *clientConnection) DeleteClient(client entities.Cliente) error {
	err := db.connection.Delete(&client).Error
	if err != nil {
		return err
	}

	return nil
}

func (db *clientConnection) FindClients(clientName string, clientType string) ([]entities.Cliente, error) {
	clients := []entities.Cliente{}

	clientName = fmt.Sprint("%", clientName, "%")

	var err error

	switch {
	case clientName != "" && clientType != "":
		err = db.connection.Find(&clients, "nome LIKE ? AND tipo = ?", clientName, clientType).Error
	case clientName != "":
		err = db.connection.Find(&clients, "nome LIKE ?", clientName).Error
	case clientType != "":
		err = db.connection.Find(&clients, "tipo = ?", clientType).Error
	default:
		err = db.connection.Find(&clients).Error
	}

	if err != nil {
		return clients, err
	}

	return clients, nil
}

// NewClientRepository cria uma nova instancia de ClientRepository.
func NewClientRepository() ClientRepository {
	return &clientConnection{
		connection: database.GetDB(),
	}
}
