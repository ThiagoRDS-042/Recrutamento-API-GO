package services

import (
	"fmt"
	"net/http"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/utils"
	"github.com/mashingan/smapping"
)

// ClientService representa a interface de clientService.
type ClientService interface {
	CreateClient(clientDTO dtos.ClientCreateDTO) (entities.Cliente, utils.ResponseError)
	UpdateClient(clientDTO dtos.ClientUpdateDTO) (entities.Cliente, utils.ResponseError)
	FindClientByID(clientID string) entities.Cliente
	FindClientByName(name string) entities.Cliente
	DeleteClient(clientID string) utils.ResponseError
	FindClients(clientName string, clientType string) []entities.Cliente
}

type clientService struct {
	clientRepository repositories.ClientRepository
	pointService     PointService
}

func (service *clientService) CreateClient(clientDTO dtos.ClientCreateDTO) (entities.Cliente, utils.ResponseError) {
	client := entities.Cliente{}

	err := smapping.FillStruct(&client, smapping.MapFields(&clientDTO))
	if err != nil {
		return entities.Cliente{},
			utils.NewResponseError(fmt.Sprintf("failed to map: %v", err), http.StatusInternalServerError)
	}

	clientAlreadyExists := service.clientRepository.FindClientByName(clientDTO.Nome)

	switch {
	case clientAlreadyExists.DataRemocao.Valid:
		client.ID = clientAlreadyExists.ID

		client, err := service.clientRepository.UpdateClient(client)
		if err != nil {
			return entities.Cliente{}, utils.NewResponseError(err.Error(), http.StatusInternalServerError)
		}

		return client, utils.ResponseError{}

	case (clientAlreadyExists != entities.Cliente{}):
		return entities.Cliente{}, utils.NewResponseError(utils.NameAlreadyExists, http.StatusConflict)

	default:
		client, err := service.clientRepository.CreateClient(client)
		if err != nil {
			return entities.Cliente{}, utils.NewResponseError(err.Error(), http.StatusInternalServerError)
		}

		return client, utils.ResponseError{}
	}
}

func (service *clientService) UpdateClient(clientDTO dtos.ClientUpdateDTO) (entities.Cliente, utils.ResponseError) {
	client := entities.Cliente{}

	err := smapping.FillStruct(&client, smapping.MapFields(&clientDTO))
	if err != nil {
		return entities.Cliente{},
			utils.NewResponseError(fmt.Sprintf("failed to map: %v", err), http.StatusInternalServerError)
	}

	clientFound := service.clientRepository.FindClientByID(client.ID)

	if clientFound == (entities.Cliente{}) {
		return entities.Cliente{}, utils.NewResponseError(utils.ClientNotFound, http.StatusNotFound)
	}

	if client.Nome == "" {
		client.Nome = clientFound.Nome
	} else {
		if !dtos.IsValidTextLenght(client.Nome) {
			return entities.Cliente{}, utils.NewResponseError("nome: "+utils.InvalidNumberOfCaracter, http.StatusBadRequest)
		}
	}

	if client.Tipo == "" {
		client.Tipo = clientFound.Tipo
	} else {
		if !dtos.IsValidClientType(client.Tipo) {
			return entities.Cliente{}, utils.NewResponseError("tipo: "+utils.InvalidClientType, http.StatusBadRequest)
		}
	}

	clientAlreadyExists := service.clientRepository.FindClientByName(client.Nome)

	if (clientAlreadyExists != entities.Cliente{}) && (clientFound.ID != clientAlreadyExists.ID) {
		return entities.Cliente{}, utils.NewResponseError(utils.NameAlreadyExists, http.StatusConflict)
	}

	client.DataRemocao.Scan(nil)
	client, err = service.clientRepository.UpdateClient(client)
	if err != nil {
		return entities.Cliente{}, utils.NewResponseError(err.Error(), http.StatusInternalServerError)
	}

	return client, utils.ResponseError{}
}

func (service *clientService) FindClientByID(clientID string) entities.Cliente {
	return service.clientRepository.FindClientByID(clientID)
}

func (service *clientService) FindClientByName(name string) entities.Cliente {
	return service.clientRepository.FindClientByName(name)
}

func (service *clientService) DeleteClient(clientID string) utils.ResponseError {
	clientFound := service.clientRepository.FindClientByID(clientID)

	if clientFound == (entities.Cliente{}) {
		return utils.NewResponseError(utils.ClientNotFound, http.StatusNotFound)
	}

	err := service.clientRepository.DeleteClient(clientFound)
	if err != nil {
		return utils.NewResponseError(err.Error(), http.StatusInternalServerError)
	}

	responseError := service.pointService.DeletePointsByClientID(clientID)
	if len(responseError.Message) != 0 {
		return utils.NewResponseError(responseError.Message, responseError.StatusCode)
	}

	return utils.ResponseError{}
}

func (service *clientService) FindClients(clientName string, clientType string) []entities.Cliente {
	return service.clientRepository.FindClients(clientName, clientType)
}

// NewClientService cria uma nova instancia de ClientService.
func NewClientService(clientRepository repositories.ClientRepository, pointService PointService) ClientService {
	return &clientService{
		clientRepository: clientRepository,
		pointService:     pointService,
	}
}
