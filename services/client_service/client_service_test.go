package services_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	repositoriesFake "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/fake"
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
)

// TestCreateClient testa se é possivel criar um novo cliente.
func TestCreateClient(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 1.0",
		Tipo: entities.FISICO,
	}

	client, responseError := clientServiceTest.CreateClient(clientDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, client)
	require.NotEqual(t, "", client.ID)
	require.False(t, client.DataRemocao.Valid)
}

// TestCreateClient testa se não é possivel criar um novo cliente com um nome ja existente.
func TestCreateClientWithNameExistent(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 2.0",
		Tipo: entities.JURIDICO,
	}

	clientServiceTest.CreateClient(clientDTO)
	client, responseError := clientServiceTest.CreateClient(clientDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.NameAlreadyExists, responseError.Message)
	require.Equal(t, http.StatusConflict, responseError.StatusCode)

	require.Empty(t, client)
}

// TestCreateClient testa se é possivel atualizar um cliente de removido para ativo.
func TestCreateClientWithDeletedAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 3.0",
		Tipo: entities.ESPECIAL,
	}

	client, _ := clientServiceTest.CreateClient(clientDTO)
	clientServiceTest.DeleteClientByID(client.ID)
	client, responseError := clientServiceTest.CreateClient(clientDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, client)
	require.NotEqual(t, "", client.ID)
	require.False(t, client.DataRemocao.Valid)
}

// TestUpdateClient testa se é possivel atualizar um cliente existente.
func TestUpdateClient(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 4.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	newName := "Test 4.1"
	newType := entities.FISICO
	clientUpdateDTO := dtos.ClientUpdateDTO{
		Base: dtos.Base{
			ID: client.ID,
		},
		Nome: newName,
		Tipo: newType,
	}
	clientUpdated, responseError := clientServiceTest.UpdateClient(clientUpdateDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, clientUpdated)
	require.Equal(t, client.ID, clientUpdated.ID)
	require.Equal(t, newName, clientUpdated.Nome)
	require.Equal(t, newType, clientUpdated.Tipo)
	require.False(t, clientUpdated.DataRemocao.Valid)
}

// TestUpdateClientWithoutName testa se é possivel atualizar um cliente sem passar o nome.
func TestUpdateClientWithoutName(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 5.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	newType := entities.JURIDICO
	clientUpdateDTO := dtos.ClientUpdateDTO{
		Base: dtos.Base{
			ID: client.ID,
		},
		Tipo: newType,
	}
	clientUpdated, responseError := clientServiceTest.UpdateClient(clientUpdateDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, clientUpdated)
	require.Equal(t, client.ID, clientUpdated.ID)
	require.Equal(t, client.Nome, clientUpdated.Nome)
	require.Equal(t, newType, clientUpdated.Tipo)
	require.False(t, clientUpdated.DataRemocao.Valid)
}

// TestUpdateClientWithNameExistent testa se não é possivel atualizar um cliente com um nome ja existente.
func TestUpdateClientWithNameExistent(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 6.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	clientDTO2 := dtos.ClientCreateDTO{
		Nome: "Test 7.0",
		Tipo: entities.JURIDICO,
	}
	client2, _ := clientServiceTest.CreateClient(clientDTO2)

	clientUpdateDTO := dtos.ClientUpdateDTO{
		Base: dtos.Base{
			ID: client.ID,
		},
		Nome: client2.Nome,
		Tipo: client2.Tipo,
	}
	clientUpdated, responseError := clientServiceTest.UpdateClient(clientUpdateDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.NameAlreadyExists, responseError.Message)
	require.Equal(t, http.StatusConflict, responseError.StatusCode)

	require.Empty(t, clientUpdated)
}

// TestUpdateClientWithInvalidName testa se não é possivel atualizar um cliente com um nome invalido.
func TestUpdateClientWithInvalidName(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 8.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	clientUpdateDTO := dtos.ClientUpdateDTO{
		Base: dtos.Base{
			ID: client.ID,
		},
		Nome: "Te",
		Tipo: client.Tipo,
	}
	clientUpdated, responseError := clientServiceTest.UpdateClient(clientUpdateDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, "nome: "+utils.InvalidNumberOfCaracter, responseError.Message)
	require.Equal(t, http.StatusBadRequest, responseError.StatusCode)

	require.Empty(t, clientUpdated)
}

// TestUpdateClientWithoutType testa se é possivel atualizar um cliente sem passar o tipo.
func TestUpdateClientWithoutType(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 9.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	newName := "Test 9.1"
	clientUpdateDTO := dtos.ClientUpdateDTO{
		Base: dtos.Base{
			ID: client.ID,
		},
		Nome: newName,
	}
	clientUpdated, responseError := clientServiceTest.UpdateClient(clientUpdateDTO)

	require.Empty(t, responseError)

	require.NotEmpty(t, clientUpdated)
	require.Equal(t, client.ID, clientUpdated.ID)
	require.Equal(t, newName, clientUpdated.Nome)
	require.Equal(t, client.Tipo, clientUpdated.Tipo)
	require.False(t, clientUpdated.DataRemocao.Valid)
}

// TestUpdateClientWithInvalidType testa se não é possivel atualizar um cliente com um tipo invalido.
func TestUpdateClientWithInvalidType(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 10.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	clientUpdateDTO := dtos.ClientUpdateDTO{
		Base: dtos.Base{
			ID: client.ID,
		},
		Nome: client.Nome,
		Tipo: "newType",
	}
	clientUpdated, responseError := clientServiceTest.UpdateClient(clientUpdateDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, "tipo: "+utils.InvalidClientType, responseError.Message)
	require.Equal(t, http.StatusBadRequest, responseError.StatusCode)

	require.Empty(t, clientUpdated)
}

// TestUpdateClientWithInvalidID testa se não é possivel atualizar um cliente com um ID invalido.
func TestUpdateClientWithInvalidID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 11.0",
		Tipo: entities.ESPECIAL,
	}
	clientServiceTest.CreateClient(clientDTO)

	clientUpdateDTO := dtos.ClientUpdateDTO{
		Base: dtos.Base{
			ID: "",
		},
		Nome: "Test 12.0",
		Tipo: entities.JURIDICO,
	}
	clientUpdated, responseError := clientServiceTest.UpdateClient(clientUpdateDTO)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.ClientNotFound, responseError.Message)
	require.Equal(t, http.StatusNotFound, responseError.StatusCode)

	require.Empty(t, clientUpdated)
}

// TestFindClientByID testa se é possivel buscar um cliente a partir do ID.
func TestFindClientByID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 13.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	clientFound := clientServiceTest.FindClientByID(client.ID)

	require.NotEmpty(t, clientFound)
	require.Equal(t, client, clientFound)
}

// TestFindClientByIDWithInvalidID testa se não é possivel buscar um cliente a partir de um ID invalido.
func TestFindClientByIDWithInvalidID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 14.0",
		Tipo: entities.JURIDICO,
	}
	clientServiceTest.CreateClient(clientDTO)

	clientFound := clientServiceTest.FindClientByID("")

	require.Empty(t, clientFound)
}

// TestFindClientByIDWithDeletedAtValid testa se não é possivel buscar um cliente removido a partir do ID.
func TestFindClientByIDWithDeletedAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 15.0",
		Tipo: entities.JURIDICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)
	clientServiceTest.DeleteClientByID(client.ID)

	clientFound := clientServiceTest.FindClientByID(client.ID)

	require.Empty(t, clientFound)
}

// TestFindClientByName testa se é possivel buscar um cliente a partir do nome.
func TestFindClientByName(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 16.0",
		Tipo: entities.JURIDICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	clientFound := clientServiceTest.FindClientByName(client.Nome)

	require.NotEmpty(t, clientFound)
	require.Equal(t, client, clientFound)
}

// TestFindClientByNameWithInvalidName testa se não é possivel buscar um cliente a partir de um nome invalido.
func TestFindClientByNameWithInvalidName(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 16.0",
		Tipo: entities.ESPECIAL,
	}
	clientServiceTest.CreateClient(clientDTO)

	clientFound := clientServiceTest.FindClientByName("")

	require.Empty(t, clientFound)
}

// TestFindClientByNameWithDeletedAtValid testa se é possivel buscar um cliente removido a partir do nome.
func TestFindClientByNameWithDeletedAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 18.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)
	clientServiceTest.DeleteClientByID(client.ID)

	clientFound := clientServiceTest.FindClientByName(client.Nome)

	client.DataRemocao.Scan(clientFound.DataRemocao.Time)

	require.NotEmpty(t, clientFound)
	require.Equal(t, client, clientFound)
}

// TestDeleteClientByID testa se é possivel "excluir"(solfdelete) um cliente a partir do ID.
func TestDeleteClientByID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 19.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	responseError := clientServiceTest.DeleteClientByID(client.ID)

	clientFound := clientServiceTest.FindClientByID(client.ID)

	require.Empty(t, responseError)
	require.Empty(t, clientFound)
}

// TestDeleteClientByIDWithInvalidID testa se não é possivel "excluir"(solfdelete) um cliente a partir de um ID invalido.
func TestDeleteClientByIDWithInvalidID(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 20.0",
		Tipo: entities.ESPECIAL,
	}
	clientServiceTest.CreateClient(clientDTO)

	responseError := clientServiceTest.DeleteClientByID("")

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.ClientNotFound, responseError.Message)
	require.NotEmpty(t, http.StatusNotFound, responseError)
}

// TestDeleteClientByIDWithDeletedAtValid testa se não é possivel "excluir"(solfdelete) um cliente removido a partir do ID.
func TestDeleteClientByIDWithDeletedAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 21.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	clientServiceTest.DeleteClientByID(client.ID)
	responseError := clientServiceTest.DeleteClientByID(client.ID)

	require.NotEmpty(t, responseError)
	require.Equal(t, utils.ClientNotFound, responseError.Message)
	require.NotEmpty(t, http.StatusNotFound, responseError)
}

// TestFindClientsByNameAndType testa se é possivel listar todos os clientes não removidos, a partir do nome e tipo.
func TestFindClientsByNameAndType(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 22.0",
		Tipo: entities.FISICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	clients := clientServiceTest.FindClients(client.Nome, client.Tipo)

	require.NotEmpty(t, clients)
	require.Greater(t, len(clients), 0)
}

// TestFindClientsByName testa se é possivel listar todos os clientes não removidos, a partir do nome.
func TestFindClientsByName(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 23.0",
		Tipo: entities.JURIDICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	clients := clientServiceTest.FindClients(client.Nome, "")

	require.NotEmpty(t, clients)
	require.Greater(t, len(clients), 0)
}

// TestFindClientsByType testa se é possivel listar todos os clientes não removidos, a partir do tipo.
func TestFindClientsByType(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 24.0",
		Tipo: entities.JURIDICO,
	}
	client, _ := clientServiceTest.CreateClient(clientDTO)

	clients := clientServiceTest.FindClients("", client.Tipo)

	require.NotEmpty(t, clients)
	require.Greater(t, len(clients), 0)
}

// TestFindClients testa se é possivel listar todos os clientes não removidos.
func TestFindClients(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 25.0",
		Tipo: entities.ESPECIAL,
	}
	clientServiceTest.CreateClient(clientDTO)

	clients := clientServiceTest.FindClients("", "")

	require.NotEmpty(t, clients)
	require.Greater(t, len(clients), 0)
}

// TestFindClientsWithDeteletAtValid testa se não é possivel listar clientes removidos.
func TestFindClientsWithDeteletAtValid(t *testing.T) {
	clientDTO := dtos.ClientCreateDTO{
		Nome: "Test 26.0",
		Tipo: entities.ESPECIAL,
	}
	clientServiceTest.CreateClient(clientDTO)

	for i := range *dbClient {
		(*dbClient)[i].DataRemocao.Scan(time.Now())
	}

	clients := clientServiceTest.FindClients("", "")

	require.Empty(t, clients)
	require.Equal(t, len(clients), 0)
}
