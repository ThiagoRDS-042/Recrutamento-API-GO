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

// TestCreatePoint testa se é possivel criar um novo ponto.
func TestCreatePoint(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 27.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 30.0",
		Bairro:     "BairroTest 30.0",
		Numero:     30,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, responseError := pointServiceTest.CreatePoint(pointDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, point)
	require.NotEqual(t, "", point.ID)
	require.False(t, point.DataRemocao.Valid)
}

// TestCreatePointWithPointExistent testa se não é possivel criar um novo ponto com dados ja existentes.
func TestCreatePointWithPointExistent(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 28.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 31.0",
		Bairro:     "BairroTest 31.0",
		Numero:     31,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	pointServiceTest.CreatePoint(pointDTO)
	point, responseError := pointServiceTest.CreatePoint(pointDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.PointAlreadyExists, responseError.Message)
	require.Equal(t, http.StatusConflict, responseError.StatusCode)

	require.Empty(t, point)
}

// TestCreatePointwithInvalidClientID testa se não é possivel criar um novo ponto a partir de um ID do cliente invalido.
func TestCreatePointwithInvalidClientID(t *testing.T) {
	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 32.0",
		Bairro:     "BairroTest 32.0",
		Numero:     32,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  "",
		EnderecoID: address.ID,
	}
	point, responseError := pointServiceTest.CreatePoint(pointDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.ClientNotFound, responseError.Message)
	require.Equal(t, http.StatusNotFound, responseError.StatusCode)

	require.Empty(t, point)
}

// TestCreatePointwithClientDeleted testa se não é possivel criar um novo ponto a partir de um ID do cliente removido.
func TestCreatePointwithClientDeleted(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 29.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)
	clientServiceTest.DeleteClientByID(client.ID)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 33.0",
		Bairro:     "BairroTest 33.0",
		Numero:     33,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, responseError := pointServiceTest.CreatePoint(pointDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.ClientNotFound, responseError.Message)
	require.Equal(t, http.StatusNotFound, responseError.StatusCode)

	require.Empty(t, point)
}

// TestCreatePointwithInvalidAddressID testa se não é possivel criar um novo ponto a partir de um ID do endereço invalido.
func TestCreatePointwithInvalidAddressID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 30.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: "",
	}
	point, responseError := pointServiceTest.CreatePoint(pointDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.AddressNotFound, responseError.Message)
	require.Equal(t, http.StatusNotFound, responseError.StatusCode)

	require.Empty(t, point)
}

// TestCreatePointwithAddressDeleted testa se não é possivel criar um novo ponto a partir de um ID do endereço removido.
func TestCreatePointwithAddressDeleted(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 31.0",
		Tipo: entities.JURIDICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 34.0",
		Bairro:     "BairroTest 34.0",
		Numero:     34,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)
	addressServiceTest.DeleteAddressByID(address.ID)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, responseError := pointServiceTest.CreatePoint(pointDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.AddressNotFound, responseError.Message)
	require.Equal(t, http.StatusNotFound, responseError.StatusCode)

	require.Empty(t, point)
}

// TestCreatePointwithDeletedAtValid testa se é possivel atualizar um ponto removido para ativo.
func TestCreatePointwithDeletedAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 32.0",
		Tipo: entities.ESPECIAL,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 35.0",
		Bairro:     "BairroTest 35.0",
		Numero:     35,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, _ := pointServiceTest.CreatePoint(pointDTO)
	pointServiceTest.DeletePointByID(point.ID)
	point, responseError := pointServiceTest.CreatePoint(pointDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, point)
	require.NotEqual(t, "", point.ID)
	require.False(t, point.DataRemocao.Valid)
}

// TestFindPointByID testa se é possivel buscar um ponto a partir do ID.
func TestFindPointByID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 33.0",
		Tipo: entities.ESPECIAL,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 36.0",
		Bairro:     "BairroTest 36.0",
		Numero:     36,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, _ := pointServiceTest.CreatePoint(pointDTO)
	pointFound := pointServiceTest.FindPointByID(point.ID)

	require.NotEmpty(t, pointFound)
	require.Equal(t, point, pointFound)
}

// TestFindPointByIDWithInvalidID testa se não é possivel buscar um ponto a partir do ID invalido.
func TestFindPointByIDWithInvalidID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 34.0",
		Tipo: entities.ESPECIAL,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 37.0",
		Bairro:     "BairroTest 37.0",
		Numero:     37,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	pointServiceTest.CreatePoint(pointDTO)
	pointFound := pointServiceTest.FindPointByID("")

	require.Empty(t, pointFound)
}

// TestFindPointByIDWithDeletedAtValid testa se não é possivel buscar um ponto removido a partir do ID.
func TestFindPointByIDWithDeletedAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 35.0",
		Tipo: entities.ESPECIAL,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 38.0",
		Bairro:     "BairroTest 38.0",
		Numero:     38,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, _ := pointServiceTest.CreatePoint(pointDTO)
	pointServiceTest.DeletePointByID(point.ID)
	pointFound := pointServiceTest.FindPointByID(point.ID)

	require.Empty(t, pointFound)
}

// TestFindPointByClientIDAndAddressID testa se é possivel buscar um ponto a partir do ID do cliente e do endereço.
func TestFindPointByClientIDAndAddressID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 36.0",
		Tipo: entities.ESPECIAL,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 39.0",
		Bairro:     "BairroTest 39.0",
		Numero:     39,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, _ := pointServiceTest.CreatePoint(pointDTO)

	pointFound := pointServiceTest.FindPointByClientIDAndAddressID(client.ID, address.ID)

	require.NotEmpty(t, pointFound)
	require.Equal(t, point, pointFound)
}

// TestFindPointByClientIDAndAddressIDWithInvalidClientID testa se não é possivel buscar um ponto a partir do ID invalido do cliente e do ID valido do endereço.
func TestFindPointByClientIDAndAddressIDWithInvalidClientID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 37.0",
		Tipo: entities.ESPECIAL,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 40.0",
		Bairro:     "BairroTest 40.0",
		Numero:     40,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	pointServiceTest.CreatePoint(pointDTO)

	pointFound := pointServiceTest.FindPointByClientIDAndAddressID("", address.ID)

	require.Empty(t, pointFound)
}

// TestFindPointByClientIDAndAddressIDWithInvalidAddressID testa se não é possivel buscar um ponto a partir do ID valido do cliente e do ID invalido do endereço.
func TestFindPointByClientIDAndAddressIDWithInvalidAddressID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 38.0",
		Tipo: entities.ESPECIAL,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 41.0",
		Bairro:     "BairroTest 41.0",
		Numero:     41,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	pointServiceTest.CreatePoint(pointDTO)

	pointFound := pointServiceTest.FindPointByClientIDAndAddressID(client.ID, "")

	require.Empty(t, pointFound)
}

// TestFindPointByClientIDAndAddressIDWithDeletedAtValid testa se é possivel buscar um ponto removido a partir do ID do cliente e do endereço.
func TestFindPointByClientIDAndAddressIDWithDeletedAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 39.0",
		Tipo: entities.ESPECIAL,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 42.0",
		Bairro:     "BairroTest 42.0",
		Numero:     42,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, _ := pointServiceTest.CreatePoint(pointDTO)
	pointServiceTest.DeletePointByID(point.ID)

	pointFound := pointServiceTest.FindPointByClientIDAndAddressID(client.ID, address.ID)
	point.DataRemocao.Scan(pointFound.DataRemocao.Time)

	require.NotEmpty(t, pointFound)
	require.Equal(t, point, pointFound)
}

// TestDeletePointByID testa se é possivel remover um ponto a partir do ID.
func TestDeletePointByID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 40.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 43.0",
		Bairro:     "BairroTest 43.0",
		Numero:     43,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, _ := pointServiceTest.CreatePoint(pointDTO)
	responseError := pointServiceTest.DeletePointByID(point.ID)
	pointFound := pointServiceTest.FindPointByID(point.ID)

	require.Empty(t, responseError)
	require.Empty(t, pointFound)
}

// TestDeletePointByIDWithInvalidID testa se não é possivel remover um ponto a partir do ID invalido.
func TestDeletePointByIDWithInvalidID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 41.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 44.0",
		Bairro:     "BairroTest 44.0",
		Numero:     44,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	pointServiceTest.CreatePoint(pointDTO)
	responseError := pointServiceTest.DeletePointByID("")

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.PointNotFound, responseError.Message)
	require.NotEmpty(t, http.StatusNotFound, responseError.StatusCode)
}

// TestDeletePointByIDWithDeletedAtValid testa se não é possivel remover um ponto removido a partir do ID.
func TestDeletePointByIDWithDeletedAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 42.0",
		Tipo: entities.ESPECIAL,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 45.0",
		Bairro:     "BairroTest 45.0",
		Numero:     45,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, _ := pointServiceTest.CreatePoint(pointDTO)
	pointServiceTest.DeletePointByID(point.ID)
	responseError := pointServiceTest.DeletePointByID(point.ID)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.PointNotFound, responseError.Message)
	require.NotEmpty(t, http.StatusNotFound, responseError.StatusCode)
}

// TestDeletePointsByClientID testa se é possivel remover um ponto a partir do ID do cliente.
func TestDeletePointsByClientID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 43.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 46.0",
		Bairro:     "BairroTest 46.0",
		Numero:     46,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, _ := pointServiceTest.CreatePoint(pointDTO)
	responseError := pointServiceTest.DeletePointsByClientID(client.ID)
	pointFound := pointServiceTest.FindPointByID(point.ID)

	require.Empty(t, responseError)
	require.Empty(t, pointFound)
}

// TestDeletePointsByClientIDWithInvalidClientID testa se não é possivel remover um ponto a partir do ID invalido do cliente.
func TestDeletePointsByClientIDWithInvalidClientID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 44.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 47.0",
		Bairro:     "BairroTest 47.0",
		Numero:     47,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, _ := pointServiceTest.CreatePoint(pointDTO)
	pointServiceTest.DeletePointsByClientID("")
	pointFound := pointServiceTest.FindPointByID(point.ID)

	require.NotEmpty(t, pointFound)
	require.False(t, pointFound.DataRemocao.Valid)
}

// TestDeletePointsByAddressID testa se é possivel remover um ponto a partir do ID do endereço.
func TestDeletePointsByAddressID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 45.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 48.0",
		Bairro:     "BairroTest 48.0",
		Numero:     48,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, _ := pointServiceTest.CreatePoint(pointDTO)
	responseError := pointServiceTest.DeletePointsByAddressID(address.ID)
	pointFound := pointServiceTest.FindPointByID(point.ID)

	require.Empty(t, responseError)
	require.Empty(t, pointFound)
}

// TestDeletePointsByAddressIDWithInvalidAddressID testa se não é possivel remover um ponto a partir do ID invalido do endereço.
func TestDeletePointsByAddressIDWithInvalidAddressID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 46.0",
		Tipo: entities.ESPECIAL,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 49.0",
		Bairro:     "BairroTest 49.0",
		Numero:     49,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	point, _ := pointServiceTest.CreatePoint(pointDTO)
	pointServiceTest.DeletePointsByAddressID("")
	pointFound := pointServiceTest.FindPointByID(point.ID)

	require.NotEmpty(t, pointFound)
	require.False(t, pointFound.DataRemocao.Valid)
}

// TestFindPoints testa se é possivel listar todos os pontos não removidos.
func TestFindPoints(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 47.0",
		Tipo: entities.ESPECIAL,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 50.0",
		Bairro:     "BairroTest 50.0",
		Numero:     50,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	pointServiceTest.CreatePoint(pointDTO)
	points := pointServiceTest.FindPoints("", "")

	require.NotEmpty(t, points)
	require.Greater(t, len(points), 0)
}

// TestFindPointsByClienteIDAndAddressID testa se é possivel listar todos os pontos não removidos a partir do ID do cliente e endereço.
func TestFindPointsByClienteIDAndAddressID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 48.0",
		Tipo: entities.ESPECIAL,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 51.0",
		Bairro:     "BairroTest 51.0",
		Numero:     51,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	pointServiceTest.CreatePoint(pointDTO)
	points := pointServiceTest.FindPoints(client.ID, address.ID)

	require.NotEmpty(t, points)
	require.Greater(t, len(points), 0)
}

// TestFindPointsByClienteID testa se é possivel listar todos os pontos não removidos a partir do ID do cliente.
func TestFindPointsByClienteID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 49.0",
		Tipo: entities.ESPECIAL,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 52.0",
		Bairro:     "BairroTest 52.0",
		Numero:     52,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	pointServiceTest.CreatePoint(pointDTO)
	points := pointServiceTest.FindPoints(client.ID, "")

	require.NotEmpty(t, points)
	require.Greater(t, len(points), 0)
}

// TestFindPointsByAddressID testa se é possivel listar todos os pontos não removidos a partir do ID do endereço.
func TestFindPointsByAddressID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 50.0",
		Tipo: entities.ESPECIAL,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 53.0",
		Bairro:     "BairroTest 53.0",
		Numero:     53,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	pointServiceTest.CreatePoint(pointDTO)
	points := pointServiceTest.FindPoints("", address.ID)

	require.NotEmpty(t, points)
	require.Greater(t, len(points), 0)
}

// TestFindPointsWithDeteletAtValid testa se não é possivel listar pontos removidos.
func TestFindPointsWithDeteletAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 51.0",
		Tipo: entities.ESPECIAL,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	addressDTO := dtos.AddressCreateDTO{
		Logradouro: "LogradouroTest 54.0",
		Bairro:     "BairroTest 54.0",
		Numero:     54,
	}
	address, _ := addressServiceTest.CreateAddress(addressDTO)

	pointDTO := dtos.PointCreateDTO{
		ClienteID:  client.ID,
		EnderecoID: address.ID,
	}
	pointServiceTest.CreatePoint(pointDTO)

	for i := range *dbPoint {
		(*dbPoint)[i].DataRemocao.Scan(time.Now())
	}

	points := pointServiceTest.FindPoints("", "")

	require.Empty(t, points)
	require.Equal(t, len(points), 0)
}
