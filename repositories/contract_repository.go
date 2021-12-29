package repositories

import (
	"fmt"
	"log"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/database"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"gorm.io/gorm"
)

// ContractRepository representa o contracto de ContractRepository.
type ContractRepository interface {
	CreateContract(contract entities.Contrato) (entities.Contrato, error)
	UpdateContract(contract entities.Contrato) (entities.Contrato, error)
	FindContractByID(contractID string) entities.Contrato
	FindContractByPontoID(pontoID string) entities.Contrato
	DeleteContract(contract entities.Contrato) error
	FindContracts(clientID string, addressID string) []entities.Contrato
}

type contractConnection struct {
	connection *gorm.DB
}

func (db *contractConnection) CreateContract(contract entities.Contrato) (entities.Contrato, error) {
	err := db.connection.Create(&contract).Error
	if err != nil {
		return contract, err
	}

	return contract, nil
}

func (db *contractConnection) UpdateContract(contract entities.Contrato) (entities.Contrato, error) {
	err := db.connection.Save(&contract).Error
	if err != nil {
		return contract, err
	}

	return contract, nil
}

func (db *contractConnection) FindContractByID(contractID string) entities.Contrato {
	contract := entities.Contrato{}

	err := db.connection.Preload("Ponto.Cliente").Preload("Ponto.Endereco").First(&contract, "id = ?", contractID).Error
	if err != nil {
		log.Println(err.Error())
	}

	return contract
}

func (db *contractConnection) FindContractByPontoID(pontoID string) entities.Contrato {
	contract := entities.Contrato{}

	err := db.connection.Unscoped().First(&contract, "ponto_id = ?", pontoID).Error
	if err != nil {
		log.Println(err.Error())
	}

	return contract
}

func (db *contractConnection) DeleteContract(contract entities.Contrato) error {
	err := db.connection.Delete(&contract).Error
	if err != nil {
		return err
	}

	return nil
}

func (db *contractConnection) FindContracts(clientID string, addressID string) []entities.Contrato {
	contracts := []entities.Contrato{}

	var sqlQuery = "JOIN t_ponto ON t_ponto.id = t_contrato.ponto_id "

	if clientID != "" {
		sqlQuery += fmt.Sprintf("AND t_ponto.cliente_id = '%v' ", clientID)
	} else {
		sqlQuery += "AND NOT t_ponto.cliente_id IS NULL "
	}

	if addressID != "" {
		sqlQuery += fmt.Sprintf("AND t_ponto.endereco_id = '%v'", addressID)
	} else {
		sqlQuery += "AND NOT t_ponto.endereco_id IS NULL"
	}

	err := db.connection.Preload("Ponto.Cliente").Preload("Ponto.Endereco").
		Joins(sqlQuery).Find(&contracts).Error
	if err != nil {
		log.Println(err.Error())
	}

	return contracts
}

// NewContractRepository cria uma nova instancia de ContractRepository.
func NewContractRepository() ContractRepository {
	return &contractConnection{
		connection: database.GetDB(),
	}
}
