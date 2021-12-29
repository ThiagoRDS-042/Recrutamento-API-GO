package repositories

import (
	"log"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/database"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"gorm.io/gorm"
)

// ContractEventRepository representa o contracto de ContractEventRepository.
type ContractEventRepository interface {
	CreateContractEvent(contractEvent entities.ContratoEvento) (entities.ContratoEvento, error)
	FindContractEventsByContractID(contractID string) []entities.ContratoEvento
}

type contractEventConnection struct {
	connection *gorm.DB
}

func (db *contractEventConnection) CreateContractEvent(contractEvent entities.ContratoEvento) (entities.ContratoEvento, error) {
	err := db.connection.Create(&contractEvent).Error
	if err != nil {
		return contractEvent, err
	}

	log.Println(contractEvent.DataCriacao)

	return contractEvent, nil
}

func (db *contractEventConnection) FindContractEventsByContractID(contractID string) []entities.ContratoEvento {
	contractEvents := []entities.ContratoEvento{}

	err := db.connection.Find(&contractEvents, "contrato_id = ?", contractID).Error
	if err != nil {
		log.Println(err.Error())
	}

	return contractEvents
}

// NewContractEventRepository cria uma nova instancia de ContractEventRepository.
func NewContractEventRepository() ContractEventRepository {
	return &contractEventConnection{
		connection: database.GetDB(),
	}
}
