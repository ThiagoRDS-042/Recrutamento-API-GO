package services

import (
	"log"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories"
	"github.com/mashingan/smapping"
)

// ContractService representa a interface de contractService.
type ContractService interface {
	CreateContract(contractDTO dtos.ContractCreateDTO) (entities.Contrato, error)
	UpdateContract(contractDTO dtos.ContractUpdateDTO) (entities.Contrato, error)
	FindContractByID(contractID string) entities.Contrato
	FindContractByPontoID(pontoID string) entities.Contrato
	DeleteContract(contract entities.Contrato) error
	DeleteContractByPontoID(pontoID string) error
	FindContracts(clientID string, addressID string) []entities.Contrato
}

type contractService struct {
	contractRepository repositories.ContractRepository
}

func (service *contractService) CreateContract(contractDTO dtos.ContractCreateDTO) (entities.Contrato, error) {
	contract := entities.Contrato{}

	err := smapping.FillStruct(&contract, smapping.MapFields(&contractDTO))
	if err != nil {
		log.Fatalf("failed to map: %v", err)
	}

	contract, err = service.contractRepository.CreateContract(contract)
	if err != nil {
		return contract, err
	}

	return contract, nil
}

func (service *contractService) UpdateContract(contractDTO dtos.ContractUpdateDTO) (entities.Contrato, error) {
	contract := entities.Contrato{}

	err := smapping.FillStruct(&contract, smapping.MapFields(&contractDTO))
	if err != nil {
		log.Fatalf("failed to map: %v", err)
	}

	contract.DataRemocao.Scan(nil)

	contract, err = service.contractRepository.UpdateContract(contract)
	if err != nil {
		return contract, err
	}

	return contract, nil
}

func (service *contractService) FindContractByID(contractID string) entities.Contrato {
	return service.contractRepository.FindContractByID(contractID)
}

func (service *contractService) FindContractByPontoID(pontoID string) entities.Contrato {
	return service.contractRepository.FindContractByPontoID(pontoID)
}

func (service *contractService) DeleteContract(contract entities.Contrato) error {
	return service.contractRepository.DeleteContract(contract)
}

func (service *contractService) DeleteContractByPontoID(pontoID string) error {
	contract := service.contractRepository.FindContractByPontoID(pontoID)

	return service.contractRepository.DeleteContract(contract)
}

func (service *contractService) FindContracts(clientID string, addressID string) []entities.Contrato {
	return service.contractRepository.FindContracts(clientID, addressID)
}

// NewContractService cria uma nova instancia de ContractService.
func NewContractService(contractRepository repositories.ContractRepository) ContractService {
	return &contractService{
		contractRepository: contractRepository,
	}
}
