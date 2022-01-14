package services

import (
	"fmt"
	"net/http"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	repositories "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/postgres"
	services "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/contract_event_service"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/utils"
	"github.com/mashingan/smapping"
)

// ContractService representa a interface de contractService.
type ContractService interface {
	CreateContract(contractDTO dtos.ContractCreateDTO) (entities.Contrato, utils.ResponseError)
	UpdateContract(contractDTO dtos.ContractUpdateDTO) (entities.Contrato, utils.ResponseError)
	FindContractByID(contractID string) entities.Contrato
	FindContractByPontoID(pontoID string) entities.Contrato
	DeleteContractByID(contractID string) utils.ResponseError
	DeleteContractByPontoID(pontoID string) utils.ResponseError
	FindContracts(clientID string, addressID string) []entities.Contrato
}

type contractService struct {
	contractRepository   repositories.ContractRepository
	pointRepository      repositories.PointRepository
	contractEventService services.ContractEventService
}

func (service *contractService) CreateContract(contractDTO dtos.ContractCreateDTO) (entities.Contrato, utils.ResponseError) {
	contract := entities.Contrato{}

	err := smapping.FillStruct(&contract, smapping.MapFields(&contractDTO))
	if err != nil {
		return entities.Contrato{},
			utils.NewResponseError(fmt.Sprintf("failed to map: %v", err), http.StatusInternalServerError)
	}

	pontoExists := service.pointRepository.FindPointByID(contract.PontoID)
	if pontoExists == (entities.Ponto{}) {
		return entities.Contrato{}, utils.NewResponseError(utils.PointNotFound, http.StatusNotFound)
	}

	contractAlreadyExists := service.contractRepository.FindContractByPontoID(contract.PontoID)

	switch {
	case contractAlreadyExists.DataRemocao.Valid:
		contract.ID = contractAlreadyExists.ID

		contract, err := service.contractRepository.UpdateContract(contract)
		if err != nil {
			return entities.Contrato{}, utils.NewResponseError(err.Error(), http.StatusInternalServerError)
		}

		contractEventDTO := dtos.ContratoEventCreateDTO{
			ContratoID:      contract.ID,
			EstadoAnterior:  contractAlreadyExists.Estado,
			EstadoPosterior: contract.Estado,
		}

		_, responseError := service.contractEventService.CreateContractEvent(contractEventDTO)
		if len(responseError.Message) != 0 {
			return entities.Contrato{}, utils.NewResponseError(responseError.Message, responseError.StatusCode)
		}

		return contract, utils.ResponseError{}

	case (contractAlreadyExists != entities.Contrato{}):
		return entities.Contrato{}, utils.NewResponseError(utils.ContractAlreadyExists, http.StatusConflict)

	default:
		contract, err := service.contractRepository.CreateContract(contract)
		if err != nil {
			return entities.Contrato{}, utils.NewResponseError(err.Error(), http.StatusInternalServerError)
		}

		contractEventDTO := dtos.ContratoEventCreateDTO{
			ContratoID:      contract.ID,
			EstadoAnterior:  contract.Estado,
			EstadoPosterior: contract.Estado,
		}

		_, responseError := service.contractEventService.CreateContractEvent(contractEventDTO)
		if len(responseError.Message) != 0 {
			return entities.Contrato{}, utils.NewResponseError(responseError.Message, responseError.StatusCode)
		}

		return contract, utils.ResponseError{}
	}
}

func (service *contractService) UpdateContract(contractDTO dtos.ContractUpdateDTO) (entities.Contrato, utils.ResponseError) {
	contract := entities.Contrato{}

	err := smapping.FillStruct(&contract, smapping.MapFields(&contractDTO))
	if err != nil {
		return entities.Contrato{},
			utils.NewResponseError(fmt.Sprintf("failed to map: %v", err), http.StatusInternalServerError)
	}

	contractFound := service.contractRepository.FindContractByID(contract.ID)
	if contractFound == (entities.Contrato{}) {
		return entities.Contrato{}, utils.NewResponseError(utils.ContractNotFound, http.StatusNotFound)
	}

	if !dtos.IsAuthorized(contractFound.Estado, contractDTO.Estado) {
		return entities.Contrato{}, utils.NewResponseError(utils.Unathorized, http.StatusUnauthorized)
	}

	contract.PontoID = contractFound.PontoID
	contract.DataRemocao.Scan(nil)
	contract, err = service.contractRepository.UpdateContract(contract)
	if err != nil {
		return entities.Contrato{}, utils.NewResponseError(err.Error(), http.StatusInternalServerError)
	}

	contractEventDTO := dtos.ContratoEventCreateDTO{
		ContratoID:      contract.ID,
		EstadoAnterior:  contractFound.Estado,
		EstadoPosterior: contract.Estado,
	}

	_, responseError := service.contractEventService.CreateContractEvent(contractEventDTO)
	if len(responseError.Message) != 0 {
		return entities.Contrato{}, utils.NewResponseError(responseError.Message, responseError.StatusCode)
	}

	return contract, utils.ResponseError{}
}

func (service *contractService) FindContractByID(contractID string) entities.Contrato {
	return service.contractRepository.FindContractByID(contractID)
}

func (service *contractService) FindContractByPontoID(pontoID string) entities.Contrato {
	return service.contractRepository.FindContractByPontoID(pontoID)
}

func (service *contractService) DeleteContractByID(contractID string) utils.ResponseError {
	contractFound := service.contractRepository.FindContractByID(contractID)

	if contractFound == (entities.Contrato{}) {
		return utils.NewResponseError(utils.ContractNotFound, http.StatusNotFound)
	}

	err := service.contractRepository.DeleteContract(contractFound)
	if err != nil {
		return utils.NewResponseError(err.Error(), http.StatusInternalServerError)
	}

	return utils.ResponseError{}
}

func (service *contractService) DeleteContractByPontoID(pontoID string) utils.ResponseError {
	contract := service.contractRepository.FindContractByPontoID(pontoID)
	if contract == (entities.Contrato{}) {
		return utils.ResponseError{}
	}

	err := service.contractRepository.DeleteContract(contract)
	if err != nil {
		return utils.NewResponseError(err.Error(), http.StatusInternalServerError)
	}

	return utils.ResponseError{}
}

func (service *contractService) FindContracts(clientID string, addressID string) []entities.Contrato {
	return service.contractRepository.FindContracts(clientID, addressID)
}

// NewContractService cria uma nova instancia de ContractService.
func NewContractService(contractRepository repositories.ContractRepository, pointRepository repositories.PointRepository, contractEventService services.ContractEventService) ContractService {
	return &contractService{
		contractRepository:   contractRepository,
		pointRepository:      pointRepository,
		contractEventService: contractEventService,
	}
}
