package services_test

import (
	"net/http"
	"testing"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	repositoriesFake "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/fake"
	addressService "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/address_service"
	clientService "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/client_service"
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
	clientServiceTest        = clientService.NewClientService(clientRepositoryFake, pointServiceTest)
	addressServiceTest       = addressService.NewAddressService(addressRepositoryFake, pointServiceTest)
)

// TestCreateContractEvent testa se é possivel criar um novo evento contrato.
func TestCreateContractEvent(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 73.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 76.0",
		Bairro:     "BairroTest 76.0",
		Numero:     76,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, _ := pointServiceTest.CreatePoint(pointDTO)

	contractDTO := dtos.ContractCreateDTO{
		PontoID: point.ID,
		Estado:  entities.VIGOR,
	}
	contract, _ := contractServiceTest.CreateContract(contractDTO)

	contractEventDTO := dtos.ContratoEventCreateDTO{
		EstadoAnterior:  entities.VIGOR,
		EstadoPosterior: entities.DESATIVADO,
		ContratoID:      contract.ID,
	}
	contractEvent, responseError := contractEventServiceTest.CreateContractEvent(contractEventDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, contractEvent)
	require.NotEqual(t, "", contractEvent.ID)
}

// TestCreateContractEventWithInvalidContractID testa se não é possivel criar um novo evento contrato a partir do ID invalido do contrato.
func TestCreateContractEventWithInvalidContractID(t *testing.T) {
	contractEventDTO := dtos.ContratoEventCreateDTO{
		EstadoAnterior:  entities.VIGOR,
		EstadoPosterior: entities.DESATIVADO,
		ContratoID:      "",
	}
	contractEvent, responseError := contractEventServiceTest.CreateContractEvent(contractEventDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.ContractNotFound, responseError.Message)
	require.Equal(t, http.StatusNotFound, responseError.StatusCode)

	require.Empty(t, contractEvent)
}

// TestFindContractEventsByContractID testa se é possivel listar os eventos de um contrato a partir do seu ID.
func TestFindContractEventsByContractID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 74.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 77.0",
		Bairro:     "BairroTest 77.0",
		Numero:     77,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, _ := pointServiceTest.CreatePoint(pointDTO)

	contractDTO := dtos.ContractCreateDTO{
		PontoID: point.ID,
		Estado:  entities.VIGOR,
	}
	contract, _ := contractServiceTest.CreateContract(contractDTO)

	contractEventDTO := dtos.ContratoEventCreateDTO{
		EstadoAnterior:  entities.VIGOR,
		EstadoPosterior: entities.DESATIVADO,
		ContratoID:      contract.ID,
	}
	contractEventServiceTest.CreateContractEvent(contractEventDTO)

	contractEventDTO2 := dtos.ContratoEventCreateDTO{
		EstadoAnterior:  entities.DESATIVADO,
		EstadoPosterior: entities.CANCELADO,
		ContratoID:      contract.ID,
	}
	contractEventServiceTest.CreateContractEvent(contractEventDTO2)

	contractEvents := contractEventServiceTest.FindContractEventsByContractID(contract.ID)

	require.NotEmpty(t, contractEvents)
	require.Greater(t, len(contractEvents), 0)
}

// TestFindContractEventsByContractIDWithInvalidID testa se não é possivel listar os eventos de um contrato a partir do seu ID invalido.
func TestFindContractEventsByContractIDWithInvalidID(t *testing.T) {
	contractEvents := contractEventServiceTest.FindContractEventsByContractID("")

	require.Empty(t, contractEvents)
}

// TestFindContractEventsByContractIDWithDeletedAtValid testa se é possivel listar os eventos de um contrato removido a partir do seu ID.
func TestFindContractEventsByContractIDWithDeletedAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 75.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 78.0",
		Bairro:     "BairroTest 78.0",
		Numero:     78,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, _ := pointServiceTest.CreatePoint(pointDTO)

	contractDTO := dtos.ContractCreateDTO{
		PontoID: point.ID,
		Estado:  entities.VIGOR,
	}
	contract, _ := contractServiceTest.CreateContract(contractDTO)
	contractServiceTest.DeleteContractByID(contract.ID)

	contractEventDTO := dtos.ContratoEventCreateDTO{
		EstadoAnterior:  entities.VIGOR,
		EstadoPosterior: entities.DESATIVADO,
		ContratoID:      contract.ID,
	}
	contractEventServiceTest.CreateContractEvent(contractEventDTO)

	contractEventDTO2 := dtos.ContratoEventCreateDTO{
		EstadoAnterior:  entities.DESATIVADO,
		EstadoPosterior: entities.CANCELADO,
		ContratoID:      contract.ID,
	}
	contractEventServiceTest.CreateContractEvent(contractEventDTO2)

	contractEvents := contractEventServiceTest.FindContractEventsByContractID(contract.ID)

	require.NotEmpty(t, contractEvents)
	require.Greater(t, len(contractEvents), 0)
}
