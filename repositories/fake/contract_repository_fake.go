package repositories

import (
	"time"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	repositories "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/postgres"
	"github.com/gofrs/uuid"
)

// DBContract banco de dados fake de contratos para os testes
var DBContract = &[]entities.Contrato{}

type contractConnectionFake struct {
	connection        *[]entities.Contrato
	connectionClient  *[]entities.Cliente
	connectionAddress *[]entities.Endereco
	connectionPoint   *[]entities.Ponto
}

func (db *contractConnectionFake) CreateContract(contract entities.Contrato) (entities.Contrato, error) {
	contractID, _ := uuid.NewV4()

	contract.ID = contractID.String()
	contract.DataCriacao = time.Now()
	contract.DataAtualizacao = time.Now()

	*db.connection = append(*db.connection, contract)

	return contract, nil
}

func (db *contractConnectionFake) UpdateContract(contract entities.Contrato) (entities.Contrato, error) {
	contract.DataAtualizacao = time.Now()
	contract.DataRemocao.Valid = false

	for i, contractValue := range *db.connection {
		if contractValue.ID == contract.ID {
			(*db.connection)[i] = contract
		}
	}

	return contract, nil
}

func (db *contractConnectionFake) FindContractByID(contractID string) entities.Contrato {
	contract := entities.Contrato{}

	for _, contractValue := range *db.connection {
		if contractValue.ID == contractID && !contractValue.DataRemocao.Valid {
			contract = contractValue
		}
	}

	for _, point := range *db.connectionPoint {
		if contract.PontoID == point.ID {

			for _, client := range *db.connectionClient {
				if point.ClienteID == client.ID {
					contract.Ponto.Cliente = client
				}
			}

			for _, address := range *db.connectionAddress {
				if point.EnderecoID == address.ID {
					contract.Ponto.Endereco = address
				}
			}
		}
	}

	return contract
}

func (db *contractConnectionFake) FindContractByPontoID(pontoID string) entities.Contrato {
	contract := entities.Contrato{}

	for _, contractValue := range *db.connection {
		if contractValue.PontoID == pontoID {
			contract = contractValue
		}
	}

	return contract
}

func (db *contractConnectionFake) DeleteContract(contract entities.Contrato) error {
	for i, contractValue := range *db.connection {
		if contractValue.ID == contract.ID {
			(*db.connection)[i].DataRemocao.Scan(time.Now())
		}
	}

	return nil
}

func (db *contractConnectionFake) FindContracts(clientID string, addressID string) []entities.Contrato {
	contracts := []entities.Contrato{}

	if clientID != "" && addressID != "" {
		for _, contractValue := range *db.connection {
			if contractValue.Ponto.ClienteID == clientID && contractValue.Ponto.EnderecoID == addressID &&
				!contractValue.DataRemocao.Valid {
				contracts = append(contracts, contractValue)
			}
		}
	} else if clientID != "" {
		for _, contractValue := range *db.connection {
			if contractValue.Ponto.ClienteID == clientID && !contractValue.DataRemocao.Valid {
				contracts = append(contracts, contractValue)
			}
		}
	} else if addressID != "" {
		for _, contractValue := range *db.connection {
			if contractValue.Ponto.EnderecoID == addressID && !contractValue.DataRemocao.Valid {
				contracts = append(contracts, contractValue)
			}
		}
	} else {
		for _, contractValue := range *db.connection {
			if !contractValue.DataRemocao.Valid {
				contracts = append(contracts, contractValue)
			}
		}
	}

	if len(contracts) != 0 {
		for i, contract := range contracts {
			for _, point := range *db.connectionPoint {
				if contract.PontoID == point.ID {

					for _, client := range *db.connectionClient {
						if point.ClienteID == client.ID {
							contracts[i].Ponto.Cliente = client
						}
					}

					for _, address := range *db.connectionAddress {
						if point.EnderecoID == address.ID {
							contracts[i].Ponto.Endereco = address
						}
					}
				}
			}
		}

	}

	return contracts
}

// NewContractRepositoryFake cria uma nova instancia de ContractRepository para os testes.
func NewContractRepositoryFake(database *[]entities.Contrato, connectionClient *[]entities.Cliente, connectionAddress *[]entities.Endereco, connectionPoint *[]entities.Ponto) repositories.ContractRepository {
	return &contractConnectionFake{
		connection:        database,
		connectionClient:  connectionClient,
		connectionAddress: connectionAddress,
		connectionPoint:   connectionPoint,
	}
}
