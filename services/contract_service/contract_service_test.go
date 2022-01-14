package services_test

import (
	"net/http"
	"testing"
	"time"

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

// TestCreateContract testa se é possivel criar um novo contrato.
func TestCreateContract(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 52.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 55.0",
		Bairro:     "BairroTest 55.0",
		Numero:     55,
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
	contract, responseError := contractServiceTest.CreateContract(contractDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, contract)
	require.NotEqual(t, "", contract.ID)
	require.False(t, contract.DataRemocao.Valid)
}

// TestCreateContract testa se não é possivel criar um novo contrato com dados ja existentes.
func TestCreateContractWithContractExists(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 53.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 56.0",
		Bairro:     "BairroTest 56.0",
		Numero:     56,
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
	contractServiceTest.CreateContract(contractDTO)
	contract, responseError := contractServiceTest.CreateContract(contractDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.ContractAlreadyExists, responseError.Message)
	require.NotEmpty(t, http.StatusConflict, responseError.StatusCode)

	require.Empty(t, contract)
}

// TestCreateContractWithInvalidPointID testa se não é possivel criar um novo contrato a partir do id do ponto invalido.
func TestCreateContractWithInvalidPointID(t *testing.T) {
	contractDTO := dtos.ContractCreateDTO{
		Estado:  entities.VIGOR,
		PontoID: "",
	}
	contract, responseError := contractServiceTest.CreateContract(contractDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.PointNotFound, responseError.Message)
	require.NotEmpty(t, http.StatusNotFound, responseError.StatusCode)

	require.Empty(t, contract)
}

// TestCreateContractWithDeletedAtValid testa se é possivel atualizar um contrato removido para um ativo.
func TestCreateContractWithDeletedAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 54.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 57.0",
		Bairro:     "BairroTest 57.0",
		Numero:     57,
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
	contract, responseError := contractServiceTest.CreateContract(contractDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, contract)
	require.NotEqual(t, "", contract.ID)
	require.False(t, contract.DataRemocao.Valid)
}

// TestUpdateContract testa se é possivel atualizar o estado do contrato.
func TestUpdateContract(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 55.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 58.0",
		Bairro:     "BairroTest 58.0",
		Numero:     58,
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

	newState := entities.DESATIVADO
	contractUpdateDTO := dtos.ContractUpdateDTO{
		Base: dtos.Base{
			ID: contract.ID,
		},
		Estado: newState,
	}
	contractUpdated, responseError := contractServiceTest.UpdateContract(contractUpdateDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, contractUpdated)
	require.Equal(t, contract.ID, contractUpdated.ID)
	require.Equal(t, newState, contractUpdated.Estado)
	require.False(t, contractUpdated.DataRemocao.Valid)
}

// TestUpdateContractWithInvalidState testa se não é possivel atualizar o estado do contrato com um estado invalido.
func TestUpdateContractWithInvalidState(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 56.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 59.0",
		Bairro:     "BairroTest 59.0",
		Numero:     59,
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

	contractUpdateDTO := dtos.ContractUpdateDTO{
		Base: dtos.Base{
			ID: contract.ID,
		},
		Estado: entities.CANCELADO,
	}
	contractUpdated, responseError := contractServiceTest.UpdateContract(contractUpdateDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.Unathorized, responseError.Message)
	require.Equal(t, http.StatusUnauthorized, responseError.StatusCode)

	require.Empty(t, contractUpdated)
}

// TestUpdateContractWithInvalidID testa se não é possivel atualizar o estado do contrato a partir de um ID invalido.
func TestUpdateContractWithInvalidID(t *testing.T) {
	contractUpdateDTO := dtos.ContractUpdateDTO{
		Base: dtos.Base{
			ID: "",
		},
		Estado: entities.DESATIVADO,
	}
	contractUpdated, responseError := contractServiceTest.UpdateContract(contractUpdateDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.ContractNotFound, responseError.Message)
	require.Equal(t, http.StatusNotFound, responseError.StatusCode)

	require.Empty(t, contractUpdated)
}

// TestFindContractByID testa se é possivel buscar um contrato não removido a partir do ID.
func TestFindContractByID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 57.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 60.0",
		Bairro:     "BairroTest 60.0",
		Numero:     60,
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

	contractFound := contractServiceTest.FindContractByID(contract.ID)

	require.NotEmpty(t, contractFound)
	require.Equal(t, client.ID, contractFound.Ponto.Cliente.ID)
	require.Equal(t, address.ID, contractFound.Ponto.Endereco.ID)
}

// TestFindContractByIDWithInvalidID testa se não é possivel buscar um contrato não removido a partir do ID invalido.
func TestFindContractByIDWithInvalidID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 58.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 61.0",
		Bairro:     "BairroTest 61.0",
		Numero:     61,
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
	contractServiceTest.CreateContract(contractDTO)
	contractFound := contractServiceTest.FindContractByID("")

	require.Empty(t, contractFound)
}

// TestFindContractByIDWithDeletedAtValid testa se não é possivel buscar um contrato removido a partir do ID.
func TestFindContractByIDWithDeletedAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 59.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 62.0",
		Bairro:     "BairroTest 62.0",
		Numero:     62,
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
	contractFound := contractServiceTest.FindContractByID(contract.ID)

	require.Empty(t, contractFound)
}

// TestFindContractByPontoID testa se é possivel buscar um contrato não removido a partir do ID do ponto.
func TestFindContractByPontoID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 60.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 63.0",
		Bairro:     "BairroTest 63.0",
		Numero:     63,
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
	contractServiceTest.CreateContract(contractDTO)
	contractFound := contractServiceTest.FindContractByPontoID(point.ID)

	require.NotEmpty(t, contractFound)
}

// TestFindContractByPontoIDWithInvalidPointID testa se não é possivel buscar um contrato não removido a partir do ID invalido do ponto.
func TestFindContractByPontoIDWithInvalidPointID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 61.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 64.0",
		Bairro:     "BairroTest 64.0",
		Numero:     64,
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
	contractServiceTest.CreateContract(contractDTO)
	contractFound := contractServiceTest.FindContractByPontoID("")

	require.Empty(t, contractFound)
}

// TestFindContractByPontoIDWithDeletedAtValid testa se é possivel buscar um contrato removido a partir do ID do ponto.
func TestFindContractByPontoIDWithDeletedAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 62.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 65.0",
		Bairro:     "BairroTest 65.0",
		Numero:     65,
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
	contractFound := contractServiceTest.FindContractByPontoID(point.ID)

	require.NotEmpty(t, contractFound)
	require.True(t, contractFound.DataRemocao.Valid)
}

// TestDeleteContractByID testa se é possivel excluir um contrato a partir do ID.
func TestDeleteContractByID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 63.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 66.0",
		Bairro:     "BairroTest 66.0",
		Numero:     66,
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
	responseError := contractServiceTest.DeleteContractByID(contract.ID)
	contractFound := contractServiceTest.FindContractByID(contract.ID)

	require.Empty(t, responseError)
	require.Empty(t, contractFound)
}

// TestDeleteContractByIDWithInvalidID testa se não é possivel excluir um contrato a partir do ID invalido.
func TestDeleteContractByIDWithInvalidID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 64.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 67.0",
		Bairro:     "BairroTest 67.0",
		Numero:     67,
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
	responseError := contractServiceTest.DeleteContractByID("")
	contractFound := contractServiceTest.FindContractByID(contract.ID)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.ContractNotFound, responseError.Message)
	require.Equal(t, http.StatusNotFound, responseError.StatusCode)

	require.NotEmpty(t, contractFound)
}

// TestDeleteContractByIDWithDeletedAtValid testa se não é possivel excluir um contrato ja removido a partir do ID.
func TestDeleteContractByIDWithDeletedAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 65.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 68.0",
		Bairro:     "BairroTest 68.0",
		Numero:     68,
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
	responseError := contractServiceTest.DeleteContractByID(contract.ID)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.ContractNotFound, responseError.Message)
	require.Equal(t, http.StatusNotFound, responseError.StatusCode)
}

// TestDeleteContractByPontoID testa se é possivel excluir um contrato a partir do ID do ponto.
func TestDeleteContractByPontoID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 66.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 69.0",
		Bairro:     "BairroTest 69.0",
		Numero:     69,
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
	responseError := contractServiceTest.DeleteContractByPontoID(point.ID)
	contractFound := contractServiceTest.FindContractByID(contract.ID)

	require.Empty(t, responseError)

	require.Empty(t, contractFound)
}

// TestDeleteContractByPontoIDWithInvalidPointID testa se não é possivel excluir um contrato a partir do ID invalido do ponto.
func TestDeleteContractByPontoIDWithInvalidPointID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 67.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 70.0",
		Bairro:     "BairroTest 70.0",
		Numero:     70,
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
	contractServiceTest.DeleteContractByPontoID("")
	contractFound := contractServiceTest.FindContractByID(contract.ID)

	require.NotEmpty(t, contractFound)
	require.False(t, contractFound.DataRemocao.Valid)
}

// TestFindContracts testa se é possivel listar contratos não removidos.
func TestFindContracts(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 68.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 71.0",
		Bairro:     "BairroTest 71.0",
		Numero:     71,
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
	contractServiceTest.CreateContract(contractDTO)
	contracts := contractServiceTest.FindContracts("", "")

	require.NotEmpty(t, contracts)
	require.Greater(t, len(contracts), 0)
}

// TestFindContractsByClientIDAndAddressID testa se é possivel listar contratos não removidos a partir do ID do cliente e endereço.
func TestFindContractsByClientIDAndAddressID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 69.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 72.0",
		Bairro:     "BairroTest 72.0",
		Numero:     72,
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
	contractServiceTest.CreateContract(contractDTO)
	contracts := contractServiceTest.FindContracts(client.ID, address.ID)

	require.NotEmpty(t, contracts)
	require.Greater(t, len(contracts), 0)
}

// TestFindContractsByClientID testa se é possivel listar contratos não removidos a partir do ID do cliente.
func TestFindContractsByClientID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 70.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 73.0",
		Bairro:     "BairroTest 73.0",
		Numero:     73,
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
	contractServiceTest.CreateContract(contractDTO)
	contracts := contractServiceTest.FindContracts(client.ID, "")

	require.NotEmpty(t, contracts)
	require.Greater(t, len(contracts), 0)
}

// TestFindContractsByAddressID testa se é possivel listar contratos não removidos a partir do ID do endereço.
func TestFindContractsByAddressID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 71.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 74.0",
		Bairro:     "BairroTest 74.0",
		Numero:     74,
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
	contractServiceTest.CreateContract(contractDTO)
	contracts := contractServiceTest.FindContracts("", address.ID)

	require.NotEmpty(t, contracts)
	require.Greater(t, len(contracts), 0)
}

// TestFindContractsWithDeletedAtValid testa se não é possivel listar contratos removidos.
func TestFindContractsWithDeletedAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 72.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 75.0",
		Bairro:     "BairroTest 75.0",
		Numero:     75,
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
	contractServiceTest.CreateContract(contractDTO)

	for i := range *dbContract {
		(*dbContract)[i].DataRemocao.Scan(time.Now())
	}

	contracts := contractServiceTest.FindContracts("", "")

	require.Empty(t, contracts)
	require.Equal(t, len(contracts), 0)
}
