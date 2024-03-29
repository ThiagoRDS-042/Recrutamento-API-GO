package controllers

import (
	"net/http"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	services "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/contract_service"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/utils"
	"github.com/gin-gonic/gin"
)

// ContractController representa o contracto de contractController.
type ContractController interface {
	CreateContract(ctx *gin.Context)
	UpdateContract(ctx *gin.Context)
	FindContractByID(ctx *gin.Context)
	DeleteContract(ctx *gin.Context)
	FindContracts(ctx *gin.Context)
}

type contractController struct {
	contractService services.ContractService
}

// CreateContract godoc
// @Summary cria um novo contrato
// @Description rota para o cadastro de novos contratos a partir do id do ponto
// @Tags contract
// @Accept json
// @Produce json
// @Param contract body entities.Contrato true "Criar Novo Contrato"
// @Success 201 {object} entities.Contrato
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /contratos [post]
func (controller *contractController) CreateContract(ctx *gin.Context) {
	contractDTO := dtos.ContractCreateDTO{}

	if err := ctx.ShouldBindJSON(&contractDTO); err != nil {
		response := utils.NewResponse(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	contractDTO.Estado = entities.VIGOR

	contract, responseError := controller.contractService.CreateContract(contractDTO)
	if responseError != (utils.ResponseError{}) {
		response := utils.NewResponse(responseError.Message)
		ctx.AbortWithStatusJSON(responseError.StatusCode, response)
		return
	}

	ctx.JSON(http.StatusCreated, contract)
}

// UpdateContract godoc
// @Summary atualiza o contrato
// @Description rota para a atualização dos dados do contrato a partir do id
// @Tags contract
// @Accept json
// @Produce json
// @Param estado body string true "atualizar contrato" Enums(Em vigor, Desativado Temporario, Cancelado)
// @Param id path string true "id do contrato"
// @Success 200 {object} entities.Contrato
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /contrato/{id} [put]
func (controller *contractController) UpdateContract(ctx *gin.Context) {
	contractDTO := dtos.ContractUpdateDTO{}

	if err := ctx.ShouldBindJSON(&contractDTO); err != nil {
		response := utils.NewResponse(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	contractID := ctx.Param("id")

	contractDTO.ID = contractID

	contract, responseError := controller.contractService.UpdateContract(contractDTO)
	if responseError != (utils.ResponseError{}) {
		response := utils.NewResponse(responseError.Message)
		ctx.AbortWithStatusJSON(responseError.StatusCode, response)
		return
	}

	ctx.JSON(http.StatusOK, contract)
}

// FindContractByID godoc
// @Summary pesquisa o contrato
// @Description rota para a pesquisa do contrato pelo id
// @Tags contract
// @Accept json
// @Produce json
// @Param id path string true "id do contrato"
// @Success 200 {object} dtos.ContractResponse
// @Failure 404 {object} utils.Response
// @Router /contrato/{id} [get]
func (controller *contractController) FindContractByID(ctx *gin.Context) {
	contractID := ctx.Param("id")

	contractFound := controller.contractService.FindContractByID(contractID)

	if contractFound == (entities.Contrato{}) {
		response := utils.NewResponse(utils.ContractNotFound)
		ctx.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	contractResponse := dtos.CreateContractResponse(contractFound)

	ctx.JSON(http.StatusOK, contractResponse)
}

// DeleteContract godoc
// @Summary deleta o contrato
// @Description rota para a exclusão do contrato pelo id
// @Tags contract
// @Accept json
// @Produce json
// @Param id path string true "id do contrato"
// @Success 204 "No Content"
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /contrato/{id} [delete]
func (controller *contractController) DeleteContract(ctx *gin.Context) {
	contractID := ctx.Param("id")

	responseError := controller.contractService.DeleteContractByID(contractID)
	if responseError != (utils.ResponseError{}) {
		response := utils.NewResponse(responseError.Message)
		ctx.AbortWithStatusJSON(responseError.StatusCode, response)
		return
	}

	ctx.JSON(http.StatusNoContent, entities.Contrato{})
}

// FindContracts godoc
// @Summary lista os contratos existentes
// @Description rota para a listagem de todos os contratos existentes no banco de dados
// @Tags contract
// @Accept json
// @Produce json
// @Param cliente_id query string false "id do cliente"
// @Param endereco_id query string false "id do endereço"
// @Success 200 {object} []dtos.ContractResponse
// @Failure 404 {object} utils.Response
// @Router /contratos [get]
func (controller *contractController) FindContracts(ctx *gin.Context) {
	clientID := ctx.Query("cliente_id")
	addressID := ctx.Query("endereco_id")

	contracts := controller.contractService.FindContracts(clientID, addressID)

	if len(contracts) == 0 {
		response := utils.NewResponse(utils.ContractNotFound)
		ctx.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	contractsResponse := []dtos.ContractResponse{}

	for _, contract := range contracts {
		contractsResponse = append(contractsResponse, dtos.CreateContractResponse(contract))
	}

	response := map[string][]dtos.ContractResponse{
		"dados": contractsResponse,
	}

	ctx.JSON(http.StatusOK, response)
}

// NewContractController cria uma nova isnancia de ContractController.
func NewContractController(contractService services.ContractService) ContractController {
	return &contractController{
		contractService: contractService,
	}
}
