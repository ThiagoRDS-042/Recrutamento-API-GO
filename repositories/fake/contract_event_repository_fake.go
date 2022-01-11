package repositories

import (
	"time"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	repositories "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/postgres"
)

// DBContractEvent banco de dados fake de eventos de contratos para os testes
var DBContractEvent = &[]entities.ContratoEvento{}

type contractEventConnectionFake struct {
	connection *[]entities.ContratoEvento
}

func (db *contractEventConnectionFake) CreateContractEvent(contractEvent entities.ContratoEvento) (entities.ContratoEvento, error) {
	contractEvent.DataCriacao = time.Now()
	contractEvent.DataAtualizacao = time.Now()

	*db.connection = append(*db.connection, contractEvent)

	return contractEvent, nil
}

func (db *contractEventConnectionFake) FindContractEventsByContractID(contractID string) []entities.ContratoEvento {
	contractsEvent := []entities.ContratoEvento{}

	if contractID != "" {
		for _, v := range *db.connection {
			if v.ContratoID == contractID {
				contractsEvent = append(contractsEvent, v)
			}
		}
	} else {
		contractsEvent = append(contractsEvent, *db.connection...)
	}

	return contractsEvent
}

// NewContractEventRepositoryFake cria uma nova instancia de ContractEventRepository para os testes.
func NewContractEventRepositoryFake(database *[]entities.ContratoEvento) repositories.ContractEventRepository {
	return &contractEventConnectionFake{
		connection: database,
	}
}
