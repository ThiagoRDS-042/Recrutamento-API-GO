package controllers

import (
	"log"
	"net/http"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/services"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/utils"
	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
)

// ClientController representa o contracto de ClientController.
type ClientController interface {
	CreateClient(ctx *gin.Context)
	UpdateClient(ctx *gin.Context)
	FindClientByID(ctx *gin.Context)
	DeleteClient(ctx *gin.Context)
	FindClients(ctx *gin.Context)
}

type clientController struct {
	clientService services.ClientService
}

func (controller *clientController) CreateClient(ctx *gin.Context) {
	clientDTO := dtos.ClientCreateDTO{}

	if err := ctx.ShouldBindJSON(&clientDTO); err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	clientAlreadyExists, err := controller.clientService.FindClientByName(clientDTO.Nome)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	switch {
	case clientAlreadyExists.DataRemocao.Valid:
		clientUpdateDTO := dtos.ClientUpdateDTO{}
		err := smapping.FillStruct(&clientUpdateDTO, smapping.MapFields(&clientDTO))
		if err != nil {
			log.Fatalf("failed to map: %v", err)
		}

		clientUpdateDTO.ID = clientAlreadyExists.ID

		log.Println(clientUpdateDTO.Tipo)

		client, err := controller.clientService.UpdateClient(clientUpdateDTO)
		if err != nil {
			response := utils.BuildErrorResponse(err.Error())
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		ctx.JSON(http.StatusCreated, client)

	case (clientAlreadyExists != entities.Cliente{}):
		response := utils.BuildErrorResponse(utils.NameAlreadyExists)
		ctx.JSON(http.StatusConflict, response)

	default:
		client, err := controller.clientService.CreateClient(clientDTO)
		if err != nil {
			response := utils.BuildErrorResponse(err.Error())
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		ctx.JSON(http.StatusCreated, client)
	}
}

func (controller *clientController) UpdateClient(ctx *gin.Context) {
	clientDTO := dtos.ClientUpdateDTO{}

	if err := ctx.ShouldBindJSON(&clientDTO); err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	clientID := ctx.Param("id")

	clientFound, err := controller.clientService.FindClientByID(clientID)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if clientFound == (entities.Cliente{}) {
		response := utils.BuildErrorResponse(utils.ClientNotFound)
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	clientDTO.ID = clientID

	if clientDTO.Nome == "" {
		clientDTO.Nome = clientFound.Nome
	} else {
		if !entities.IsValidClientName(clientDTO.Nome) {
			response := utils.BuildErrorResponse("nome: " + utils.InvalidNumberOfCaracter)
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
	}

	if clientDTO.Tipo == "" {
		clientDTO.Tipo = clientFound.Tipo
	} else {
		if !entities.IsValidClientType(clientDTO.Tipo) {
			response := utils.BuildErrorResponse("tipo: " + utils.InvalidClientType)
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
	}

	clientAlreadyExists, err := controller.clientService.FindClientByName(clientDTO.Nome)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if (clientAlreadyExists != entities.Cliente{}) && (clientFound.ID != clientAlreadyExists.ID) {
		response := utils.BuildErrorResponse(utils.NameAlreadyExists)
		ctx.JSON(http.StatusConflict, response)
		return
	}

	client, err := controller.clientService.UpdateClient(clientDTO)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	ctx.JSON(http.StatusOK, client)
}

func (controller *clientController) FindClientByID(ctx *gin.Context) {
	clientID := ctx.Param("id")

	clientFound, err := controller.clientService.FindClientByID(clientID)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if clientFound == (entities.Cliente{}) {
		response := utils.BuildErrorResponse(utils.ClientNotFound)
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, clientFound)
}

func (controller *clientController) DeleteClient(ctx *gin.Context) {
	clientID := ctx.Param("id")

	clientFound, err := controller.clientService.FindClientByID(clientID)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if clientFound == (entities.Cliente{}) {
		response := utils.BuildErrorResponse(utils.ClientNotFound)
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	controller.clientService.DeleteClient(clientFound)

	ctx.JSON(http.StatusNoContent, entities.Cliente{})
}

func (controller *clientController) FindClients(ctx *gin.Context) {
	clientType := ctx.Query("tipo")
	clientName := ctx.Query("nome")

	clients, err := controller.clientService.FindClients(clientName, clientType)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if len(clients) == 0 {
		response := utils.BuildErrorResponse(utils.ClientNotFound)
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	response := map[string][]entities.Cliente{
		"dados": clients,
	}

	ctx.JSON(http.StatusOK, response)
}

// NewClientController cria uma nova isnancia de ClientController.
func NewClientController() ClientController {
	return &clientController{
		clientService: services.NewClientService(),
	}
}
