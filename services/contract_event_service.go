package services

import (
	"log"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories"
	"github.com/mashingan/smapping"
)

// ContractEventService representa a interface de ContractEventService.
type ContractEventService interface {
	CreateContractEvent(contractEventDTO dtos.ContratoEventCreateDTO) (entities.ContratoEvento, error)
	FindContractEventsByContractID(contractID string) []entities.ContratoEvento
}

type contractEventService struct {
	contractEventRepository repositories.ContractEventRepository
}

func (service *contractEventService) CreateContractEvent(contractEventDTO dtos.ContratoEventCreateDTO) (entities.ContratoEvento, error) {
	contractEvent := entities.ContratoEvento{}

	err := smapping.FillStruct(&contractEvent, smapping.MapFields(&contractEventDTO))
	if err != nil {
		log.Fatalf("failed to map: %v", err)
	}

	log.Println(contractEvent.DataCriacao)

	contractEvent, err = service.contractEventRepository.CreateContractEvent(contractEvent)
	if err != nil {
		return contractEvent, err
	}

	return contractEvent, nil
}

func (service *contractEventService) FindContractEventsByContractID(contractID string) []entities.ContratoEvento {
	return service.contractEventRepository.FindContractEventsByContractID(contractID)
}

// NewContractEventService cria uma nova instancia de ContractEventService.
func NewContractEventService(contractEventRepository repositories.ContractEventRepository) ContractEventService {
	return &contractEventService{
		contractEventRepository: contractEventRepository,
	}
}
