package services

import (
	"log"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories"
	"github.com/mashingan/smapping"
)

// PointService representa a interface de pointService.
type PointService interface {
	CreatePoint(pointDTO dtos.PointCreateDTO) (entities.Ponto, error)
	UpdatePoint(pointDTO dtos.PointUpdateDTO) (entities.Ponto, error)
	FindPointByID(pointID string) entities.Ponto
	FindPointByClientIDAndAddressID(clientID string, addressID string) entities.Ponto
	DeletePoint(point entities.Ponto) error
	DeletePointsByClientID(clientID string) error
	DeletePointsByAddressID(addressID string) error
	FindPoints(clientID string, addressID string) []entities.Ponto
}

type pointService struct {
	pointRepository repositories.PointRepository
	contractService ContractService
}

func (service *pointService) CreatePoint(pointDTO dtos.PointCreateDTO) (entities.Ponto, error) {
	point := entities.Ponto{}

	err := smapping.FillStruct(&point, smapping.MapFields(&pointDTO))
	if err != nil {
		log.Fatalf("failed to map: %v", err)
	}

	point, err = service.pointRepository.CreatePoint(point)
	if err != nil {
		return point, err
	}

	return point, nil
}

func (service *pointService) UpdatePoint(pointDTO dtos.PointUpdateDTO) (entities.Ponto, error) {
	point := entities.Ponto{}

	err := smapping.FillStruct(&point, smapping.MapFields(&pointDTO))
	if err != nil {
		log.Fatalf("failed to map: %v", err)
	}

	point.DataRemocao.Scan(nil)

	point, err = service.pointRepository.UpdatePoint(point)
	if err != nil {
		return point, err
	}

	return point, nil
}

func (service *pointService) FindPointByID(pointID string) entities.Ponto {
	return service.pointRepository.FindPointByID(pointID)
}

func (service *pointService) FindPointByClientIDAndAddressID(clientID string, addressID string) entities.Ponto {
	return service.pointRepository.FindPointByClientIDAndAddressID(clientID, addressID)
}

func (service *pointService) DeletePoint(point entities.Ponto) error {
	return service.pointRepository.DeletePoint(point)
}

func (service *pointService) DeletePointsByClientID(clientID string) error {
	points := service.pointRepository.FindPointsByClientID(clientID)

	var err error

	for _, point := range points {
		err = service.pointRepository.DeletePoint(point)
		if err != nil {
			return err
		}

		err = service.contractService.DeleteContractByPontoID(point.ID)
		if err != nil {
			return err
		}
	}

	return err
}

func (service *pointService) DeletePointsByAddressID(addressID string) error {
	points := service.pointRepository.FindPointsByAddressID(addressID)

	var err error

	for _, point := range points {
		err = service.pointRepository.DeletePoint(point)
		if err != nil {
			return err
		}

		err = service.contractService.DeleteContractByPontoID(point.ID)
		if err != nil {
			return err
		}
	}

	return err
}

func (service *pointService) FindPoints(clientID string, addressID string) []entities.Ponto {
	return service.pointRepository.FindPoints(clientID, addressID)
}

// NewPointService cria uma nova instancia de PointService.
func NewPointService(pointRepository repositories.PointRepository, contractService ContractService) PointService {
	return &pointService{
		pointRepository: pointRepository,
		contractService: contractService,
	}
}
