package services

import (
	"log"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories"
	"github.com/mashingan/smapping"
)

// ClientService representa a interface de clientService.
type ClientService interface {
	CreateClient(clientDTO dtos.ClientCreateDTO) (entities.Cliente, error)
	UpdateClient(clientDTO dtos.ClientUpdateDTO) (entities.Cliente, error)
	FindClientByID(clientID string) entities.Cliente
	FindClientByName(name string) entities.Cliente
	DeleteClient(client entities.Cliente) error
	FindClients(clientName string, clientType string) []entities.Cliente
}

type clientService struct {
	clientRepository repositories.ClientRepository
}

func (service *clientService) CreateClient(clientDTO dtos.ClientCreateDTO) (entities.Cliente, error) {
	client := entities.Cliente{}

	err := smapping.FillStruct(&client, smapping.MapFields(&clientDTO))
	if err != nil {
		log.Fatalf("failed to map: %v", err)
	}

	client, err = service.clientRepository.CreateClient(client)
	if err != nil {
		return client, err
	}

	return client, nil
}

func (service *clientService) UpdateClient(clientDTO dtos.ClientUpdateDTO) (entities.Cliente, error) {
	client := entities.Cliente{}

	err := smapping.FillStruct(&client, smapping.MapFields(&clientDTO))
	if err != nil {
		log.Fatalf("failed to map: %v", err)
	}

	client.DataRemocao.Scan(nil)
	client, err = service.clientRepository.UpdateClient(client)
	if err != nil {
		return client, err
	}

	return client, nil
}

func (service *clientService) FindClientByID(clientID string) entities.Cliente {
	return service.clientRepository.FindClientByID(clientID)
}

func (service *clientService) FindClientByName(name string) entities.Cliente {
	return service.clientRepository.FindClientByName(name)
}

func (service *clientService) DeleteClient(client entities.Cliente) error {
	return service.clientRepository.DeleteClient(client)
}

func (service *clientService) FindClients(clientName string, clientType string) []entities.Cliente {
	return service.clientRepository.FindClients(clientName, clientType)
}

// NewClientService cria uma nova instancia de ClientService.
func NewClientService() ClientService {
	return &clientService{
		clientRepository: repositories.NewClientRepository(),
	}
}
