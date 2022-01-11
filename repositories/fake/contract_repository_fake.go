package repositories

import (
	"time"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	repositories "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/postgres"
)

// DBContract banco de dados fake de contratos para os testes
var DBContract = &[]entities.Contrato{}

type contractConnectionFake struct {
	connection *[]entities.Contrato
}

func (db *contractConnectionFake) CreateContract(contract entities.Contrato) (entities.Contrato, error) {
	contract.DataCriacao = time.Now()
	contract.DataAtualizacao = time.Now()

	*db.connection = append(*db.connection, contract)

	return contract, nil
}

func (db *contractConnectionFake) UpdateContract(contract entities.Contrato) (entities.Contrato, error) {
	contract.DataAtualizacao = time.Now()

	for i, v := range *db.connection {
		if v.ID == contract.ID {
			(*db.connection)[i] = contract
		}
	}

	return contract, nil
}

func (db *contractConnectionFake) FindContractByID(contractID string) entities.Contrato {
	contract := entities.Contrato{}

	for _, v := range *db.connection {
		if v.ID == contractID && !v.DataRemocao.Valid {
			contract = v
		}
	}

	return contract
	// err := db.connection.Preload("Ponto.Cliente").Preload("Ponto.Endereco").First(&contract, "id = ?", contractID).Error
}

func (db *contractConnectionFake) FindContractByPontoID(pontoID string) entities.Contrato {
	contract := entities.Contrato{}

	for _, v := range *db.connection {
		if v.PontoID == pontoID {
			contract = v
		}
	}

	return contract
}

func (db *contractConnectionFake) DeleteContract(contract entities.Contrato) error {
	for i, v := range *db.connection {
		if v.ID == contract.ID {
			(*db.connection)[i].DataRemocao.Scan(time.Now())
		}
	}

	return nil
}

func (db *contractConnectionFake) FindContracts(clientID string, addressID string) []entities.Contrato {
	contracts := []entities.Contrato{}

	if clientID != "" && addressID != "" {
		for _, v := range *db.connection {
			if v.Ponto.ClienteID == clientID && v.Ponto.EnderecoID == addressID && !v.DataRemocao.Valid {
				contracts = append(contracts, v)
			}
		}
	} else if clientID != "" {
		for _, v := range *db.connection {
			if v.Ponto.ClienteID == clientID && !v.DataRemocao.Valid {
				contracts = append(contracts, v)
			}
		}
	} else if addressID != "" {
		for _, v := range *db.connection {
			if v.Ponto.EnderecoID == addressID && !v.DataRemocao.Valid {
				contracts = append(contracts, v)
			}
		}
	} else {
		contracts = append(contracts, *db.connection...)
	}

	return contracts
	// err := db.connection.Preload("Ponto.Cliente").Preload("Ponto.Endereco").
}

// NewContractRepositoryFake cria uma nova instancia de ContractRepository para os testes.
func NewContractRepositoryFake(database *[]entities.Contrato) repositories.ContractRepository {
	return &contractConnectionFake{
		connection: database,
	}
}
