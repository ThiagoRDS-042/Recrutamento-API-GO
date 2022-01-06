package services

import (
	"fmt"
	"net/http"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/utils"
	"github.com/mashingan/smapping"
)

// ContractEventService representa a interface de ContractEventService.
type ContractEventService interface {
	CreateContractEvent(contractEventDTO dtos.ContratoEventCreateDTO) (entities.ContratoEvento, utils.ResponseError)
	FindContractEventsByContractID(contractID string) []entities.ContratoEvento
}

type contractEventService struct {
	contractEventRepository repositories.ContractEventRepository
	contractRepository      repositories.ContractRepository
}

func (service *contractEventService) CreateContractEvent(contractEventDTO dtos.ContratoEventCreateDTO) (entities.ContratoEvento, utils.ResponseError) {
	contractEvent := entities.ContratoEvento{}

	err := smapping.FillStruct(&contractEvent, smapping.MapFields(&contractEventDTO))
	if err != nil {
		return entities.ContratoEvento{},
			utils.NewResponseError(fmt.Sprintf("failed to map: %v", err), http.StatusInternalServerError)
	}

	contractFound := service.contractRepository.FindContractByID(contractEvent.ContratoID)
	if contractFound == (entities.Contrato{}) {
		return entities.ContratoEvento{},
			utils.NewResponseError(utils.ContractNotFound, http.StatusNotFound)
	}

	contractEvent, err = service.contractEventRepository.CreateContractEvent(contractEvent)
	if err != nil {
		return entities.ContratoEvento{}, utils.NewResponseError(err.Error(), http.StatusInternalServerError)
	}

	return contractEvent, utils.ResponseError{}
}

func (service *contractEventService) FindContractEventsByContractID(contractID string) []entities.ContratoEvento {
	return service.contractEventRepository.FindContractEventsByContractID(contractID)
}

// NewContractEventService cria uma nova instancia de ContractEventService.
func NewContractEventService(contractEventRepository repositories.ContractEventRepository, contractRepository repositories.ContractRepository) ContractEventService {
	return &contractEventService{
		contractEventRepository: contractEventRepository,
		contractRepository:      contractRepository,
	}
}
