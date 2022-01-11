package repositories

import (
	"time"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	repositories "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/postgres"
)

// DBClient banco de dados fake de clientes para os testes
var DBClient = &[]entities.Cliente{}

type clientConnectionFake struct {
	connection *[]entities.Cliente
}

func (db *clientConnectionFake) CreateClient(client entities.Cliente) (entities.Cliente, error) {
	client.DataCriacao = time.Now()
	client.DataAtualizacao = time.Now()

	*db.connection = append(*db.connection, client)

	return client, nil
}

func (db *clientConnectionFake) UpdateClient(client entities.Cliente) (entities.Cliente, error) {
	client.DataAtualizacao = time.Now()

	for i, v := range *db.connection {
		if v.ID == client.ID {
			(*db.connection)[i] = client
		}
	}

	return client, nil
}

func (db *clientConnectionFake) FindClientByID(clientID string) entities.Cliente {
	client := entities.Cliente{}

	for _, v := range *db.connection {
		if v.ID == clientID && !v.DataRemocao.Valid {
			client = v
		}
	}

	return client
}

func (db *clientConnectionFake) FindClientByName(name string) entities.Cliente {
	client := entities.Cliente{}

	for _, v := range *db.connection {
		if v.Nome == name {
			client = v
		}
	}

	return client
}

func (db *clientConnectionFake) DeleteClient(client entities.Cliente) error {
	for i, v := range *db.connection {
		if v.ID == client.ID {
			(*db.connection)[i].DataRemocao.Scan(time.Now())
		}
	}

	return nil
}

func (db *clientConnectionFake) FindClients(clientName string, clientType string) []entities.Cliente {
	return *db.connection
}

// NewClientRepositoryFake cria uma nova instancia de ClientRepository para os testes.
func NewClientRepositoryFake(database *[]entities.Cliente) repositories.ClientRepository {
	return &clientConnectionFake{
		connection: database,
	}
}
