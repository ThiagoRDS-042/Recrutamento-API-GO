package services_test

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	repositoriesFake "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/fake"
	addressService "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/address_service"
	contractEventService "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/contract_event_service"
	contractService "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/contract_service"
	pointService "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/point_service"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/utils"
	"github.com/stretchr/testify/require"
)

var (
	// Fake Databases
	dbClient        = repositoriesFake.DBClient
	dbAddress       = repositoriesFake.DBAddress
	dbPoint         = repositoriesFake.DBPoint
	dbContract      = repositoriesFake.DBContract
	dbContractEvent = repositoriesFake.DBContractEvent

	// Fake Repositories
	clientRepositoryFake        = repositoriesFake.NewClientRepositoryFake(dbClient)
	addressRepositoryFake       = repositoriesFake.NewAddressRepositoryFake(dbAddress)
	pointRepositoryFake         = repositoriesFake.NewPointRepositoryFake(dbPoint, dbClient, dbAddress)
	contractRepositoryFake      = repositoriesFake.NewContractRepositoryFake(dbContract, dbClient, dbAddress, dbPoint)
	contractEventRepositoryFake = repositoriesFake.NewContractEventRepositoryFake(dbContractEvent)

	// Services Tests
	contractEventServiceTest = contractEventService.NewContractEventService(contractEventRepositoryFake, contractRepositoryFake)
	contractServiceTest      = contractService.NewContractService(contractRepositoryFake, pointRepositoryFake, contractEventServiceTest)
	pointServiceTest         = pointService.NewPointService(pointRepositoryFake, clientRepositoryFake, addressRepositoryFake, contractServiceTest)
	addressServiceTest       = addressService.NewAddressService(addressRepositoryFake, pointServiceTest)
)

// TestCreateAddress testa se é possivel criar um novo endereço.
func TestCreateAddress(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 1.0",
		Bairro:     "BairroTest 1.0",
		Numero:     1,
	}

	address, responseError := addressServiceTest.CreateAddress(addressDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, address)
	require.NotEqual(t, "", address.ID)
	require.False(t, address.DataRemocao.Valid)
}

// TestCreateAddressWithAddressExistent testa se não é possivel criar um novo endrereço com os dados de um endereço ja existente.
func TestCreateAddressWithAddressExistent(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 2.0",
		Bairro:     "BairroTest 2.0",
		Numero:     2,
	}

	addressServiceTest.CreateAddress(addressDTO)
	address, responseError := addressServiceTest.CreateAddress(addressDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.AddressAlreadyExists, responseError.Message)
	require.Equal(t, http.StatusConflict, responseError.StatusCode)

	require.Empty(t, address)
}

// TestCreateAddressWithDeletedAtValid testa se é possivel atualizar um endereço de removido para ativo.
func TestCreateAddressWithDeletedAtValid(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 3.0",
		Bairro:     "BairroTest 3.0",
		Numero:     3,
	}

	address, _ := addressServiceTest.CreateAddress(addressDTO)
	addressServiceTest.DeleteAddressByID(address.ID)
	address, responseError := addressServiceTest.CreateAddress(addressDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, address)
	require.NotEqual(t, "", address.ID)
	require.False(t, address.DataRemocao.Valid)
}

// TestUpdateAddress testa se é possivel atualizar um endereço existente.
func TestUpdateAddress(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 4.0",
		Bairro:     "BairroTest 4.0",
		Numero:     4,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	newStreet := "LogradouroTest 4.1"
	newNeightbohood := "BairroTest 4.1"
	newNumber := 4
	addressUpdateDTO := dtos.AddressUpdateDTO{
		Base: dtos.Base{
			ID: address.ID,
		},
		Logradouro: newStreet,
		Bairro:     newNeightbohood,
		Numero:     newNumber,
	}
	addressUpdated, responseError := addressServiceTest.UpdateAddress(addressUpdateDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, addressUpdated)
	require.Equal(t, address.ID, addressUpdated.ID)
	require.Equal(t, newStreet, addressUpdated.Logradouro)
	require.Equal(t, newNeightbohood, addressUpdated.Bairro)
	require.Equal(t, newNumber, addressUpdated.Numero)
	require.False(t, addressUpdated.DataRemocao.Valid)
}

// TestUpdateAddressWithAddressExistent testa se não é possivel atualizar um endereço com dados ja existentes.
func TestUpdateAddressWithAddressExistent(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 5.0",
		Bairro:     "BairroTest 5.0",
		Numero:     5,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	addressDTO2 := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 6.0",
		Bairro:     "BairroTest 6.0",
		Numero:     6,
	}
	address2, _ := addressServiceTest.CreateAddress(addressDTO2)

	addressUpdateDTO := dtos.AddressUpdateDTO{
		Base: dtos.Base{
			ID: address.ID,
		},
		Logradouro: address2.Logradouro,
		Bairro:     address2.Bairro,
		Numero:     address2.Numero,
	}
	addressUpdated, responseError := addressServiceTest.UpdateAddress(addressUpdateDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.AddressAlreadyExists, responseError.Message)
	require.NotEmpty(t, http.StatusConflict, responseError.StatusCode)

	require.Empty(t, addressUpdated)
}

// TestUpdateAddressWithoutStreet testa se é possivel atualizar um endereço sem passar o logradouro.
func TestUpdateAddressWithoutStreet(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 7.0",
		Bairro:     "BairroTest 7.0",
		Numero:     7,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	newNeightbohood := "BairroTest 7.1"
	newNumber := 7
	addressUpdateDTO := dtos.AddressUpdateDTO{
		Base: dtos.Base{
			ID: address.ID,
		},
		Bairro: newNeightbohood,
		Numero: newNumber,
	}
	addressUpdated, responseError := addressServiceTest.UpdateAddress(addressUpdateDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, addressUpdated)
	require.Equal(t, address.ID, addressUpdated.ID)
	require.Equal(t, address.Logradouro, addressUpdated.Logradouro)
	require.Equal(t, newNeightbohood, addressUpdated.Bairro)
	require.Equal(t, newNumber, addressUpdated.Numero)
	require.False(t, addressUpdated.DataRemocao.Valid)
}

// TestUpdateAddressWithoutNeighborhood testa se é possivel atualizar um endereço sem passar o bairro.
func TestUpdateAddressWithoutNeighborhood(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 8.0",
		Bairro:     "BairroTest 8.0",
		Numero:     8,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	newStreet := "LogradouroTest 8.1"
	newNumber := 8
	addressUpdateDTO := dtos.AddressUpdateDTO{
		Base: dtos.Base{
			ID: address.ID,
		},
		Logradouro: newStreet,
		Numero:     newNumber,
	}
	addressUpdated, responseError := addressServiceTest.UpdateAddress(addressUpdateDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, addressUpdated)
	require.Equal(t, address.ID, addressUpdated.ID)
	require.Equal(t, newStreet, addressUpdated.Logradouro)
	require.Equal(t, address.Bairro, addressUpdated.Bairro)
	require.Equal(t, newNumber, addressUpdated.Numero)
	require.False(t, addressUpdated.DataRemocao.Valid)
}

// TestUpdateAddressWithoutNumber testa se é possivel atualizar um endereço sem passar o numero.
func TestUpdateAddressWithoutNumber(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 9.0",
		Bairro:     "BairroTest 9.0",
		Numero:     9,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	newStreet := "LogradouroTest 9.1"
	newNeightbohood := "BairroTest 9.1"
	addressUpdateDTO := dtos.AddressUpdateDTO{
		Base: dtos.Base{
			ID: address.ID,
		},
		Logradouro: newStreet,
		Bairro:     newNeightbohood,
	}
	addressUpdated, responseError := addressServiceTest.UpdateAddress(addressUpdateDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, addressUpdated)
	require.Equal(t, address.ID, addressUpdated.ID)
	require.Equal(t, newStreet, addressUpdated.Logradouro)
	require.Equal(t, newNeightbohood, addressUpdated.Bairro)
	require.Equal(t, address.Numero, addressUpdated.Numero)
	require.False(t, addressUpdated.DataRemocao.Valid)
}

// TestUpdateAddressWithInvalidStreet testa se não é possivel atualizar um endereço passando um logradouro invalido.
func TestUpdateAddressWithInvalidStreet(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 10.0",
		Bairro:     "BairroTest 10.0",
		Numero:     10,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	addressUpdateDTO := dtos.AddressUpdateDTO{
		Base: dtos.Base{
			ID: address.ID,
		},
		Logradouro: "Lo",
		Bairro:     address.Bairro,
		Numero:     address.Numero,
	}
	addressUpdated, responseError := addressServiceTest.UpdateAddress(addressUpdateDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, "logradouro: "+utils.InvalidNumberOfCaracter, responseError.Message)
	require.Equal(t, http.StatusBadRequest, responseError.StatusCode)

	require.Empty(t, addressUpdated)
}

// TestUpdateAddressWithInvalidNeighborhood testa se não é possivel atualizar um endereço passando um bairro invalido.
func TestUpdateAddressWithInvalidNeighborhood(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 11.0",
		Bairro:     "BairroTest 11.0",
		Numero:     11,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	addressUpdateDTO := dtos.AddressUpdateDTO{
		Base: dtos.Base{
			ID: address.ID,
		},
		Logradouro: address.Logradouro,
		Bairro:     "Ba",
		Numero:     address.Numero,
	}
	addressUpdated, responseError := addressServiceTest.UpdateAddress(addressUpdateDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, "bairro: "+utils.InvalidNumberOfCaracter, responseError.Message)
	require.Equal(t, http.StatusBadRequest, responseError.StatusCode)

	require.Empty(t, addressUpdated)
}

// TestUpdateAddressWithInvalidID testa se não é possivel atualizar os dadods do endereço a partir de um ID invalido.
func TestUpdateAddressWithInvalidID(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 12.0",
		Bairro:     "BairroTest 12.0",
		Numero:     12,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	addressUpdateDTO := dtos.AddressUpdateDTO{
		Base: dtos.Base{
			ID: "",
		},
		Logradouro: address.Logradouro,
		Bairro:     address.Bairro,
		Numero:     address.Numero,
	}
	addressUpdated, responseError := addressServiceTest.UpdateAddress(addressUpdateDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.AddressNotFound, responseError.Message)
	require.Equal(t, http.StatusNotFound, responseError.StatusCode)

	require.Empty(t, addressUpdated)
}

// TestFindAddressByID testa se é possivel buscar um endereço a partir do ID.
func TestFindAddressByID(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 13.0",
		Bairro:     "BairroTest 13.0",
		Numero:     13,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	addressFound := addressServiceTest.FindAddressByID(address.ID)

	require.NotEmpty(t, addressFound)
	require.Equal(t, address, addressFound)
}

// TestFindAddressByIDWithoutInvalidID testa se não é possivel buscar um endereço a partir de um ID invalido.
func TestFindAddressByIDWithoutInvalidID(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 14.0",
		Bairro:     "BairroTest 14.0",
		Numero:     14,
	}
	addressServiceTest.CreateAddress(addressDTO)

	addressFound := addressServiceTest.FindAddressByID("")

	require.Empty(t, addressFound)
}

// TestFindAddressByIDWithDeletedAtValid testa se não é possivel buscar um endereço removido a partir do ID.
func TestFindAddressByIDWithDeletedAtValid(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 15.0",
		Bairro:     "BairroTest 15.0",
		Numero:     15,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)
	addressServiceTest.DeleteAddressByID(address.ID)

	addressFound := addressServiceTest.FindAddressByID(address.ID)

	require.Empty(t, addressFound)
}

// TestFindAddressByFields testa se é possivel buscar um endereço a partir do logradouro, bairro e numero.
func TestFindAddressByFields(t *testing.T) {
	street := "LogradouroTest 16.0"
	neighborhood := "BairroTest 16.0"
	number := 16

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: street,
		Bairro:     neighborhood,
		Numero:     number,
	}
	addressServiceTest.CreateAddress(addressDTO)

	addressFound := addressServiceTest.FindAddressByFields(street, neighborhood, number)

	require.NotEmpty(t, addressFound)
	require.Equal(t, street, addressFound.Logradouro)
	require.NotEmpty(t, neighborhood, addressFound.Bairro)
	require.NotEmpty(t, number, addressFound.Numero)
}

// TestFindAddressByFieldsWithInvalidFields testa se não é possivel buscar um endereço a partir de um logradouro, bairro e numero invalidos.
func TestFindAddressByFieldsWithInvalidFields(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 17.0",
		Bairro:     "BairroTest 17.0",
		Numero:     17,
	}
	addressServiceTest.CreateAddress(addressDTO)

	addressFound := addressServiceTest.FindAddressByFields("", "", 0)

	require.Empty(t, addressFound)
}

// TestFindAddressByFieldsWithInvalidFields testa se é possivel buscar um endereço removido a partir do logradouro, bairro e numero.
func TestFindAddressByFieldsWithDeletedAtValid(t *testing.T) {
	street := "LogradouroTest 18.0"
	neighborhood := "BairroTest 18.0"
	number := 18

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: street,
		Bairro:     neighborhood,
		Numero:     number,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)
	addressServiceTest.DeleteAddressByID(address.ID)

	addressFound := addressServiceTest.FindAddressByFields(street, neighborhood, number)

	require.NotEmpty(t, addressFound)
	require.Equal(t, street, addressFound.Logradouro)
	require.NotEmpty(t, neighborhood, addressFound.Bairro)
	require.NotEmpty(t, number, addressFound.Numero)
}

// TestDeleteAddressByID testa se é possivel "excluir"(solfdelete) um endereço a partir do ID.
func TestDeleteAddressByID(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 19.0",
		Bairro:     "BairroTest 19.0",
		Numero:     19,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	responseError := addressServiceTest.DeleteAddressByID(address.ID)

	addressFound := addressServiceTest.FindAddressByID(address.ID)

	require.Empty(t, responseError)
	require.Empty(t, addressFound)
}

// TestDeleteAddressByIDWithInvalidID testa se não é possivel "excluir"(solfdelete) um endereço a partir de um ID invalido.
func TestDeleteAddressByIDWithInvalidID(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 20.0",
		Bairro:     "BairroTest 20.0",
		Numero:     20,
	}
	addressServiceTest.CreateAddress(addressDTO)

	responseError := addressServiceTest.DeleteAddressByID("")

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.AddressNotFound, responseError.Message)
	require.NotEmpty(t, http.StatusNotFound, responseError)
}

// TestDeleteAddressByIDWithDeletedAtValid testa se não é possivel "excluir"(solfdelete) um endereço removido a partir do ID.
func TestDeleteAddressByIDWithDeletedAtValid(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 21.0",
		Bairro:     "BairroTest 21.0",
		Numero:     21,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	addressServiceTest.DeleteAddressByID(address.ID)
	responseError := addressServiceTest.DeleteAddressByID(address.ID)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.AddressNotFound, responseError.Message)
	require.NotEmpty(t, http.StatusNotFound, responseError)
}

// TestFindAddresses testa se é possivel listar todos os endereços não removidos.
func TestFindAddresses(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 22.0",
		Bairro:     "BairroTest 22.0",
		Numero:     22,
	}
	addressServiceTest.CreateAddress(addressDTO)

	addresses := addressServiceTest.FindAddresses("", "", "")

	require.NotEmpty(t, addresses)
	require.Greater(t, len(addresses), 0)
}

// TestFindAddressesWithDeteletAtValid testa se não é possivel listar endereços removidos.
func TestFindAddressesWithDeteletAtValid(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 23.0",
		Bairro:     "BairroTest 23.0",
		Numero:     23,
	}
	addressServiceTest.CreateAddress(addressDTO)

	for i := range *dbAddress {
		(*dbAddress)[i].DataRemocao.Scan(time.Now())
	}

	addresses := addressServiceTest.FindAddresses("", "", "")

	require.Empty(t, addresses)
	require.Equal(t, len(addresses), 0)
}

// TestFindAddressesByFields testa se é possivel listar todos os endereços não removidos a partir do logradouro, bairro e numero.
func TestFindAddressesByFields(t *testing.T) {
	street := "LogradouroTest 23.0"
	neighborhood := "BairroTest 23.0"
	number := 23

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: street,
		Bairro:     neighborhood,
		Numero:     number,
	}
	addressServiceTest.CreateAddress(addressDTO)

	addresses := addressServiceTest.FindAddresses(street, neighborhood, strconv.Itoa(number))

	require.NotEmpty(t, addresses)
	require.Greater(t, len(addresses), 0)
}

// TestFindAddressesByStreetAndNeighborhood testa se é possivel listar todos os endereços não removidos a partir do logradouro e bairro.
func TestFindAddressesByStreetAndNeighborhood(t *testing.T) {
	street := "LogradouroTest 24.0"
	neighborhood := "BairroTest 24.0"
	number := 24

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: street,
		Bairro:     neighborhood,
		Numero:     number,
	}
	addressServiceTest.CreateAddress(addressDTO)

	addresses := addressServiceTest.FindAddresses(street, neighborhood, "")

	require.NotEmpty(t, addresses)
	require.Greater(t, len(addresses), 0)
}

// TestFindAddressesByStreetAndNumber testa se é possivel listar todos os endereços não removidos a partir do logradouro e numero.
func TestFindAddressesByStreetAndNumber(t *testing.T) {
	street := "LogradouroTest 25.0"
	neighborhood := "BairroTest 25.0"
	number := 25

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: street,
		Bairro:     neighborhood,
		Numero:     number,
	}
	addressServiceTest.CreateAddress(addressDTO)

	addresses := addressServiceTest.FindAddresses(street, "", strconv.Itoa(number))

	require.NotEmpty(t, addresses)
	require.Greater(t, len(addresses), 0)
}

// TestFindAddressesByNeighborhoodAndNumber testa se é possivel listar todos os endereços não removidos a partir do bairro e numero.
func TestFindAddressesByNeighborhoodAndNumber(t *testing.T) {
	street := "LogradouroTest 26.0"
	neighborhood := "BairroTest 26.0"
	number := 26

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: street,
		Bairro:     neighborhood,
		Numero:     number,
	}
	addressServiceTest.CreateAddress(addressDTO)

	addresses := addressServiceTest.FindAddresses("", neighborhood, strconv.Itoa(number))

	require.NotEmpty(t, addresses)
	require.Greater(t, len(addresses), 0)
}

// TestFindAddressesByStreet testa se é possivel listar todos os endereços não removidos a partir do logradouro.
func TestFindAddressesByStreet(t *testing.T) {
	street := "LogradouroTest 27.0"
	neighborhood := "BairroTest 27.0"
	number := 27

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: street,
		Bairro:     neighborhood,
		Numero:     number,
	}
	addressServiceTest.CreateAddress(addressDTO)

	addresses := addressServiceTest.FindAddresses(street, "", "")

	require.NotEmpty(t, addresses)
	require.Greater(t, len(addresses), 0)
}

// TestFindAddressesByNeighborhood testa se é possivel listar todos os endereços não removidos a partir do bairro.
func TestFindAddressesByNeighborhood(t *testing.T) {
	street := "LogradouroTest 28.0"
	neighborhood := "BairroTest 28.0"
	number := 28

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: street,
		Bairro:     neighborhood,
		Numero:     number,
	}
	addressServiceTest.CreateAddress(addressDTO)

	addresses := addressServiceTest.FindAddresses("", neighborhood, "")

	require.NotEmpty(t, addresses)
	require.Greater(t, len(addresses), 0)
}

// TestFindAddressesByNumber testa se é possivel listar todos os endereços não removidos a partir do numero.
func TestFindAddressesByNumber(t *testing.T) {
	street := "LogradouroTest 29.0"
	neighborhood := "BairroTest 29.0"
	number := 29

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: street,
		Bairro:     neighborhood,
		Numero:     number,
	}
	addressServiceTest.CreateAddress(addressDTO)

	addresses := addressServiceTest.FindAddresses("", "", strconv.Itoa(number))

	require.NotEmpty(t, addresses)
	require.Greater(t, len(addresses), 0)
}
