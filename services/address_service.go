package services

import (
	"log"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories"
	"github.com/mashingan/smapping"
)

// AddressService representa a interface de addressService.
type AddressService interface {
	CreateAddress(addressDTO dtos.AddressCreateDTO) (entities.Endereco, error)
	UpdateAddress(addressDTO dtos.AddressUpdateDTO) (entities.Endereco, error)
	FindAddressByID(addressID string) entities.Endereco
	FindAddressByFields(street string, neighborhood string, number int) entities.Endereco
	DeleteAddress(address entities.Endereco) error
	FindAddresses(street string, neighborhood string, number string) []entities.Endereco
}

type addressService struct {
	addressRepository repositories.AddressRepository
}

func (service *addressService) CreateAddress(addressDTO dtos.AddressCreateDTO) (entities.Endereco, error) {
	address := entities.Endereco{}

	err := smapping.FillStruct(&address, smapping.MapFields(&addressDTO))
	if err != nil {
		log.Fatalf("failed to map: %v", err)
	}

	address, err = service.addressRepository.CreateAddress(address)
	if err != nil {
		return address, err
	}

	return address, nil
}

func (service *addressService) UpdateAddress(addressDTO dtos.AddressUpdateDTO) (entities.Endereco, error) {
	address := entities.Endereco{}

	err := smapping.FillStruct(&address, smapping.MapFields(&addressDTO))
	if err != nil {
		log.Fatalf("failed to map: %v", err)
	}

	address.ID = addressDTO.ID
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
func NewAddressService() AddressService {
	return &addressService{
		addressRepository: repositories.NewAddressRepository(),
	}
}
