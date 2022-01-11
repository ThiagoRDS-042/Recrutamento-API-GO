package services

import (
	"fmt"
	"net/http"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	repositories "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/postgres"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/utils"
	"github.com/mashingan/smapping"
)

// AddressService representa a interface de addressService.
type AddressService interface {
	CreateAddress(addressDTO dtos.AddressCreateDTO) (entities.Endereco, utils.ResponseError)
	UpdateAddress(addressDTO dtos.AddressUpdateDTO) (entities.Endereco, utils.ResponseError)
	FindAddressByID(addressID string) entities.Endereco
	FindAddressByFields(street string, neighborhood string, number int) entities.Endereco
	DeleteAddressByID(addressID string) utils.ResponseError
	FindAddresses(street string, neighborhood string, number string) []entities.Endereco
}

type addressService struct {
	addressRepository repositories.AddressRepository
	pointService      PointService
}

func (service *addressService) CreateAddress(addressDTO dtos.AddressCreateDTO) (entities.Endereco, utils.ResponseError) {
	address := entities.Endereco{}

	err := smapping.FillStruct(&address, smapping.MapFields(&addressDTO))
	if err != nil {
		return entities.Endereco{},
			utils.NewResponseError(fmt.Sprintf("failed to map: %v", err), http.StatusInternalServerError)
	}

	addressAlreadyExists := service.FindAddressByFields(
		address.Logradouro, address.Bairro, address.Numero)

	switch {
	case addressAlreadyExists.DataRemocao.Valid:
		address.ID = addressAlreadyExists.ID

		address, err := service.addressRepository.UpdateAddress(address)
		if err != nil {
			return entities.Endereco{}, utils.NewResponseError(err.Error(), http.StatusInternalServerError)
		}

		return address, utils.ResponseError{}

	case (addressAlreadyExists != entities.Endereco{}):
		return entities.Endereco{}, utils.NewResponseError(utils.AddressAlreadyExists, http.StatusConflict)

	default:
		address, err := service.addressRepository.CreateAddress(address)
		if err != nil {
			return entities.Endereco{}, utils.NewResponseError(err.Error(), http.StatusInternalServerError)
		}

		return address, utils.ResponseError{}
	}
}

func (service *addressService) UpdateAddress(addressDTO dtos.AddressUpdateDTO) (entities.Endereco, utils.ResponseError) {
	address := entities.Endereco{}

	err := smapping.FillStruct(&address, smapping.MapFields(&addressDTO))
	if err != nil {
		return entities.Endereco{},
			utils.NewResponseError(fmt.Sprintf("failed to map: %v", err), http.StatusInternalServerError)
	}

	addressFound := service.addressRepository.FindAddressByID(address.ID)

	if addressFound == (entities.Endereco{}) {
		return entities.Endereco{}, utils.NewResponseError(utils.AddressNotFound, http.StatusNotFound)

	}

	if address.Logradouro == "" {
		address.Logradouro = addressFound.Logradouro
	} else {
		if !dtos.IsValidTextLenght(address.Logradouro) {
			return entities.Endereco{}, utils.NewResponseError("logradouro: "+utils.InvalidNumberOfCaracter, http.StatusBadRequest)
		}
	}

	if address.Bairro == "" {
		address.Bairro = addressFound.Bairro
	} else {
		if !dtos.IsValidTextLenght(address.Bairro) {
			return entities.Endereco{}, utils.NewResponseError("bairro: "+utils.InvalidNumberOfCaracter, http.StatusBadRequest)
		}
	}

	if address.Numero == 0 {
		address.Numero = addressFound.Numero
	}

	addressAlreadyExists := service.addressRepository.FindAddressByFields(
		address.Logradouro, address.Bairro, address.Numero)

	if (addressAlreadyExists != entities.Endereco{}) && (addressFound.ID != addressAlreadyExists.ID) {
		return entities.Endereco{}, utils.NewResponseError(utils.AddressAlreadyExists, http.StatusConflict)
	}

	address.DataRemocao.Scan(nil)
	address, err = service.addressRepository.UpdateAddress(address)
	if err != nil {
		return entities.Endereco{}, utils.NewResponseError(err.Error(), http.StatusInternalServerError)
	}

	return address, utils.ResponseError{}
}

func (service *addressService) FindAddressByID(addressID string) entities.Endereco {
	return service.addressRepository.FindAddressByID(addressID)
}

func (service *addressService) FindAddressByFields(street string, neighborhood string, number int) entities.Endereco {
	return service.addressRepository.FindAddressByFields(street, neighborhood, number)
}

func (service *addressService) DeleteAddressByID(addressID string) utils.ResponseError {

	addressFound := service.addressRepository.FindAddressByID(addressID)

	if addressFound == (entities.Endereco{}) {
		return utils.NewResponseError(utils.AddressNotFound, http.StatusNotFound)
	}

	err := service.addressRepository.DeleteAddress(addressFound)
	if err != nil {
		return utils.NewResponseError(err.Error(), http.StatusInternalServerError)
	}

	responseError := service.pointService.DeletePointsByAddressID(addressID)
	if len(responseError.Message) != 0 {
		return utils.NewResponseError(responseError.Message, responseError.StatusCode)
	}

	return utils.ResponseError{}
}

func (service *addressService) FindAddresses(street string, neighborhood string, number string) []entities.Endereco {
	return service.addressRepository.FindAddresses(street, neighborhood, number)
}

// NewAddressService cria uma nova instancia de AddressService.
func NewAddressService(addressRepository repositories.AddressRepository, pointService PointService) AddressService {
	return &addressService{
		addressRepository: addressRepository,
		pointService:      pointService,
	}
}
