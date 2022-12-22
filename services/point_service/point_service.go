package services

import (
	"fmt"
	"net/http"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	repositories "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/postgres"
	services "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/contract_service"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/utils"
	"github.com/mashingan/smapping"
)

// PointService representa a interface de pointService.
type PointService interface {
	CreatePoint(pointDTO dtos.PointCreateDTO) (entities.Ponto, utils.ResponseError)
	FindPointByID(pointID string) entities.Ponto
	FindPointByClientIDAndAddressID(clientID string, addressID string) entities.Ponto
	DeletePointByID(pointID string) utils.ResponseError
	DeletePointsByClientID(clientID string) utils.ResponseError
	DeletePointsByAddressID(addressID string) utils.ResponseError
	FindPoints(clientID string, addressID string) []entities.Ponto
}

type pointService struct {
	pointRepository   repositories.PointRepository
	clientRepository  repositories.ClientRepository
	addressReporitory repositories.AddressRepository
	contractService   services.ContractService
}

func (service *pointService) CreatePoint(pointDTO dtos.PointCreateDTO) (entities.Ponto, utils.ResponseError) {
	point := entities.Ponto{}

	err := smapping.FillStruct(&point, smapping.MapFields(&pointDTO))
	if err != nil {
		return entities.Ponto{},
			utils.NewResponseError(fmt.Sprintf("failed to map: %v", err), http.StatusInternalServerError)
	}

	clientExists := service.clientRepository.FindClientByID(pointDTO.ClienteID)
	if clientExists == (entities.Cliente{}) {
		return entities.Ponto{}, utils.NewResponseError(utils.ClientNotFound, http.StatusNotFound)
	}

	addressExists := service.addressReporitory.FindAddressByID(pointDTO.EnderecoID)
	if addressExists == (entities.Endereco{}) {
		return entities.Ponto{}, utils.NewResponseError(utils.AddressNotFound, http.StatusNotFound)
	}

	pointAlreadyExists := service.pointRepository.FindPointByClientIDAndAddressID(
		point.ClienteID, point.EnderecoID)

	switch {
	case pointAlreadyExists.DataRemocao.Valid:
		point.ID = pointAlreadyExists.ID

		point, err := service.pointRepository.UpdatePoint(point)
		if err != nil {
			return entities.Ponto{}, utils.NewResponseError(err.Error(), http.StatusInternalServerError)
		}

		return point, utils.ResponseError{}

	case (pointAlreadyExists != entities.Ponto{}):
		return entities.Ponto{}, utils.NewResponseError(utils.PointAlreadyExists, http.StatusConflict)

	default:
		point, err := service.pointRepository.CreatePoint(point)
		if err != nil {
			return entities.Ponto{}, utils.NewResponseError(err.Error(), http.StatusInternalServerError)
		}

		return point, utils.ResponseError{}
	}
}

func (service *pointService) FindPointByID(pointID string) entities.Ponto {
	return service.pointRepository.FindPointByID(pointID)
}

func (service *pointService) FindPointByClientIDAndAddressID(clientID string, addressID string) entities.Ponto {
	return service.pointRepository.FindPointByClientIDAndAddressID(clientID, addressID)
}

func (service *pointService) DeletePointByID(pointID string) utils.ResponseError {
	pointFound := service.pointRepository.FindPointByID(pointID)

	if pointFound == (entities.Ponto{}) {
		return utils.NewResponseError(utils.PointNotFound, http.StatusNotFound)
	}

	err := service.pointRepository.DeletePoint(pointFound)
	if err != nil {
		return utils.NewResponseError(err.Error(), http.StatusInternalServerError)
	}

	responseError := service.contractService.DeleteContractByPontoID(pointID)
	if responseError != (utils.ResponseError{}) {
		return responseError
	}

	return utils.ResponseError{}
}

func (service *pointService) DeletePointsByClientID(clientID string) utils.ResponseError {
	points := service.pointRepository.FindPointsByClientID(clientID)

	if len(points) == 0 {
		return utils.ResponseError{}
	}

	var err error

	for _, point := range points {
		err = service.pointRepository.DeletePoint(point)
		if err != nil {
			return utils.NewResponseError(err.Error(), http.StatusInternalServerError)
		}

		responseError := service.contractService.DeleteContractByPontoID(point.ID)
		if responseError != (utils.ResponseError{}) {
			return responseError
		}
	}

	return utils.ResponseError{}
}

func (service *pointService) DeletePointsByAddressID(addressID string) utils.ResponseError {
	points := service.pointRepository.FindPointsByAddressID(addressID)

	if len(points) == 0 {
		return utils.ResponseError{}
	}

	var err error

	for _, point := range points {
		err = service.pointRepository.DeletePoint(point)
		if err != nil {
			return utils.NewResponseError(err.Error(), http.StatusInternalServerError)
		}

		responseError := service.contractService.DeleteContractByPontoID(point.ID)
		if responseError != (utils.ResponseError{}) {
			return responseError
		}
	}

	return utils.ResponseError{}
}

func (service *pointService) FindPoints(clientID string, addressID string) []entities.Ponto {
	return service.pointRepository.FindPoints(clientID, addressID)
}

// NewPointService cria uma nova instancia de PointService.
func NewPointService(pointRepository repositories.PointRepository, clientRepository repositories.ClientRepository, addressReporitory repositories.AddressRepository, contractService services.ContractService) PointService {
	return &pointService{
		pointRepository:   pointRepository,
		contractService:   contractService,
		clientRepository:  clientRepository,
		addressReporitory: addressReporitory,
	}
}
