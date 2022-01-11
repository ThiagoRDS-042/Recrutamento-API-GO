package controllers

import (
	"net/http"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	services "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/client_service"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/utils"
	"github.com/gin-gonic/gin"
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

// CreateClient godoc
// @Summary cria um novo cliente
// @Description rota para o cadastro de novos clientes
// @Tags client
// @Accept json
// @Produce json
// @Param client body entities.Cliente true "Criar Novo Cliente"
// @Success 201 {object} entities.Cliente
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /clientes [post]
func (controller *clientController) CreateClient(ctx *gin.Context) {
	clientDTO := dtos.ClientCreateDTO{}

	if err := ctx.ShouldBindJSON(&clientDTO); err != nil {
		response := utils.NewResponse(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	client, responseError := controller.clientService.CreateClient(clientDTO)
	if responseError != (utils.ResponseError{}) {
		response := utils.NewResponse(responseError.Message)
		ctx.AbortWithStatusJSON(responseError.StatusCode, response)
		return
	}

	ctx.JSON(http.StatusCreated, client)
}

// UpdateClient godoc
// @Summary atualiza o cliente
// @Description rota para a atualização dos dados do cliente a partir do id
// @Tags client
// @Accept json
// @Produce json
// @Param client body entities.Cliente true "atualizar cliente"
// @Param id path string true "id do cliente"
// @Success 200 {object} entities.Cliente
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /cliente/{id} [put]
func (controller *clientController) UpdateClient(ctx *gin.Context) {
	clientDTO := dtos.ClientUpdateDTO{}

	if err := ctx.ShouldBindJSON(&clientDTO); err != nil {
		response := utils.NewResponse(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	clientID := ctx.Param("id")

	clientDTO.ID = clientID

	client, responseError := controller.clientService.UpdateClient(clientDTO)
	if responseError != (utils.ResponseError{}) {
		response := utils.NewResponse(responseError.Message)
		ctx.AbortWithStatusJSON(responseError.StatusCode, response)
		return
	}

	ctx.JSON(http.StatusOK, client)
}

// FindClientByID godoc
// @Summary pesquisa o cliente
// @Description rota para a pesquisa do cliente pelo id
// @Tags client
// @Accept json
// @Produce json
// @Param id path string true "id do cliente"
// @Success 200 {object} entities.Cliente
// @Failure 404 {object} utils.Response
// @Router /cliente/{id} [get]
func (controller *clientController) FindClientByID(ctx *gin.Context) {
	clientID := ctx.Param("id")

	clientFound := controller.clientService.FindClientByID(clientID)

	if clientFound == (entities.Cliente{}) {
		response := utils.NewResponse(utils.ClientNotFound)
		ctx.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, clientFound)
}

// DeleteClient godoc
// @Summary deleta o cliente
// @Description rota para a exclusão do cliente pelo id
// @Tags client
// @Accept json
// @Produce json
// @Param id path string true "id do cliente"
// @Success 204 "No Content"
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /cliente/{id} [delete]
func (controller *clientController) DeleteClient(ctx *gin.Context) {
	clientID := ctx.Param("id")

	responseError := controller.clientService.DeleteClientByID(clientID)
	if responseError != (utils.ResponseError{}) {
		response := utils.NewResponse(responseError.Message)
		ctx.AbortWithStatusJSON(responseError.StatusCode, response)
		return
	}

	ctx.JSON(http.StatusNoContent, entities.Cliente{})
}

// FindClients godoc
// @Summary lista os clientes existentes
// @Description rota para a listagem de todos os clientes existentes no banco de dados
// @Tags client
// @Accept json
// @Produce json
// @Param tipo query string false "tipo de cliente"
// @Param nome query string false "nome do cliente"
// @Success 200 {object} []entities.Cliente
// @Failure 404 {object} utils.Response
// @Router /clientes [get]
func (controller *clientController) FindClients(ctx *gin.Context) {
	clientType := ctx.Query("tipo")
	clientName := ctx.Query("nome")

	clients := controller.clientService.FindClients(clientName, entities.ClientType(clientType))

	if len(clients) == 0 {
		response := utils.NewResponse(utils.ClientNotFound)
		ctx.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	response := map[string][]entities.Cliente{
		"dados": clients,
	}

	ctx.JSON(http.StatusOK, response)
}

// NewClientController cria uma nova isnancia de ClientController.
func NewClientController(clientService services.ClientService) ClientController {
	return &clientController{
		clientService: clientService,
	}
}
