package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/utils"
	"github.com/mashingan/smapping"
)

// AddressService representa a interface de addressService.
type AddressService interface {
	CreateAddress(addressDTO dtos.AddressCreateDTO) (entities.Endereco, utils.ResponseError)
	UpdateAddress(addressDTO dtos.AddressUpdateDTO) (entities.Endereco, error)
	FindAddressByID(addressID string) entities.Endereco
	FindAddressByFields(street string, neighborhood string, number int) entities.Endereco
	DeleteAddress(address entities.Endereco) error
	FindAddresses(street string, neighborhood string, number string) []entities.Endereco
}

type addressService struct {
	addressRepository repositories.AddressRepository
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

func (service *addressService) UpdateAddress(addressDTO dtos.AddressUpdateDTO) (entities.Endereco, error) {
	address := entities.Endereco{}

	err := smapping.FillStruct(&address, smapping.MapFields(&addressDTO))
	if err != nil {
		log.Fatalf("failed to map: %v", err)
	}

	address.DataRemocao.Scan(nil)

	address, err = service.addressRepository.UpdateAddress(address)
	if err != nil {
		return address, err
	}

	return address, nil
}

func (service *addressService) FindAddressByID(addressID string) entities.Endereco {
	return service.addressRepository.FindAddressByID(addressID)
}

func (service *addressService) FindAddressByFields(street string, neighborhood string, number int) entities.Endereco {
	return service.addressRepository.FindAddressByFields(street, neighborhood, number)
}

func (service *addressService) DeleteAddress(address entities.Endereco) error {
	return service.addressRepository.DeleteAddress(address)
}

func (service *addressService) FindAddresses(street string, neighborhood string, number string) []entities.Endereco {
	return service.addressRepository.FindAddresses(street, neighborhood, number)
}

// NewAddressService cria uma nova instancia de AddressService.
func NewAddressService(addressRepository repositories.AddressRepository) AddressService {
	return &addressService{
		addressRepository: addressRepository,
	}
}
