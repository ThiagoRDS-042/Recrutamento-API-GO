package repositories

import (
	"strings"
	"time"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	repositories "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/postgres"
	"github.com/gofrs/uuid"
)

// DBClient banco de dados fake de clientes para os testes
var DBClient = &[]entities.Cliente{}

type clientConnectionFake struct {
	connection *[]entities.Cliente
}

func (db *clientConnectionFake) CreateClient(client entities.Cliente) (entities.Cliente, error) {
	clientID, _ := uuid.NewV4()

	client.ID = clientID.String()
	client.DataCriacao = time.Now()
	client.DataAtualizacao = time.Now()

	*db.connection = append(*db.connection, client)

	return client, nil
}

func (db *clientConnectionFake) UpdateClient(client entities.Cliente) (entities.Cliente, error) {
	client.DataAtualizacao = time.Now()
	client.DataRemocao.Valid = false

	for i, clientValue := range *db.connection {
		if clientValue.ID == client.ID {
			(*db.connection)[i] = client
		}
	}

	return client, nil
}

func (db *clientConnectionFake) FindClientByID(clientID string) entities.Cliente {
	client := entities.Cliente{}

	for _, clientValue := range *db.connection {
		if clientValue.ID == clientID && !clientValue.DataRemocao.Valid {
			client = clientValue
		}
	}

	return client
}

func (db *clientConnectionFake) FindClientByName(name string) entities.Cliente {
	client := entities.Cliente{}

	for _, clientValue := range *db.connection {
		if clientValue.Nome == name {
			client = clientValue
		}
	}

	return client
}

func (db *clientConnectionFake) DeleteClient(client entities.Cliente) error {
	for i, clientValue := range *db.connection {
		if clientValue.ID == client.ID {
			(*db.connection)[i].DataRemocao.Scan(time.Now())
		}
	}

	return nil
}

func (db *clientConnectionFake) FindClients(clientName string, clientType entities.ClientType) []entities.Cliente {
	clients := []entities.Cliente{}

	if clientName != "" && clientType != entities.ClientType("") {
		for _, clientValue := range *db.connection {
			if strings.Contains(clientValue.Nome, clientName) && clientValue.Tipo == clientType &&
				!clientValue.DataRemocao.Valid {
				clients = append(clients, clientValue)
			}
		}
	} else if clientName != "" {
		for _, clientValue := range *db.connection {
			if strings.Contains(clientValue.Nome, clientName) && !clientValue.DataRemocao.Valid {
				clients = append(clients, clientValue)
			}
		}
	} else if clientType != "" {
		for _, clientValue := range *db.connection {
			if clientValue.Tipo == clientType && !clientValue.DataRemocao.Valid {
				clients = append(clients, clientValue)
			}
		}
	} else {
		for _, clientValue := range *db.connection {
			if !clientValue.DataRemocao.Valid {
				clients = append(clients, clientValue)
			}
		}
	}

	return clients
}

// NewClientRepositoryFake cria uma nova instancia de ClientRepository para os testes.
func NewClientRepositoryFake(database *[]entities.Cliente) repositories.ClientRepository {
	return &clientConnectionFake{
		connection: database,
	}
}
