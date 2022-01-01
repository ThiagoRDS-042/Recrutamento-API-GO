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

// ContractController representa o contracto de contractController.
type ContractController interface {
	CreateContract(ctx *gin.Context)
	UpdateContract(ctx *gin.Context)
	FindContractByID(ctx *gin.Context)
	DeleteContract(ctx *gin.Context)
	FindContracts(ctx *gin.Context)
}

type contractController struct {
	contractService      services.ContractService
	contractEventService services.ContractEventService
	pontoService         services.PointService
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
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	pontoExists := controller.pontoService.FindPointByID(contractDTO.PontoID)
	if pontoExists == (entities.Ponto{}) {
		response := utils.BuildErrorResponse(utils.PointNotFound)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	contractDTO.Estado = entities.VIGOR

	contractAlreadyExists := controller.contractService.FindContractByPontoID(contractDTO.PontoID)

	switch {
	case contractAlreadyExists.DataRemocao.Valid:
		contractUpdateDTO := dtos.ContractUpdateDTO{}
		err := smapping.FillStruct(&contractUpdateDTO, smapping.MapFields(&contractDTO))
		if err != nil {
			log.Fatalf("failed to map: %v", err)
		}

		contractUpdateDTO.ID = contractAlreadyExists.ID
		contract, err := controller.contractService.UpdateContract(contractUpdateDTO)
		if err != nil {
			response := utils.BuildErrorResponse(err.Error())
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		contractEventDTO := dtos.ContratoEventCreateDTO{
			ContratoID:      contract.ID,
			EstadoAnterior:  contract.Estado,
			EstadoPosterior: contract.Estado,
		}

		_, err = controller.contractEventService.CreateContractEvent(contractEventDTO)
		if err != nil {
			response := utils.BuildErrorResponse(err.Error())
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		ctx.JSON(http.StatusCreated, contract)

	case (contractAlreadyExists != entities.Contrato{}):
		response := utils.BuildErrorResponse(utils.ContractAlreadyExists)
		ctx.JSON(http.StatusConflict, response)

	default:
		contract, err := controller.contractService.CreateContract(contractDTO)
		if err != nil {
			response := utils.BuildErrorResponse(err.Error())
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		contractEventDTO := dtos.ContratoEventCreateDTO{
			ContratoID:      contract.ID,
			EstadoAnterior:  contract.Estado,
			EstadoPosterior: contract.Estado,
		}

		_, err = controller.contractEventService.CreateContractEvent(contractEventDTO)
		if err != nil {
			response := utils.BuildErrorResponse(err.Error())
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		ctx.JSON(http.StatusCreated, contract)
	}
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
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	contractID := ctx.Param("id")

	contractFound := controller.contractService.FindContractByID(contractID)

	if contractFound == (entities.Contrato{}) {
		response := utils.BuildErrorResponse(utils.ContractNotFound)
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	if !dtos.IsAuthorized(contractFound.Estado, contractDTO.Estado) {
		response := utils.BuildErrorResponse(utils.Unathorized)
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	contractDTO.ID = contractID

	if contractDTO.PontoID == "" {
		contractDTO.PontoID = contractFound.PontoID
	}

	contract, err := controller.contractService.UpdateContract(contractDTO)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	contractEventDTO := dtos.ContratoEventCreateDTO{
		ContratoID:      contract.ID,
		EstadoAnterior:  contractFound.Estado,
		EstadoPosterior: contract.Estado,
	}

	_, err = controller.contractEventService.CreateContractEvent(contractEventDTO)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
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
		response := utils.BuildErrorResponse(utils.ContractNotFound)
		ctx.JSON(http.StatusNotFound, response)
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

	contractFound := controller.contractService.FindContractByID(contractID)

	if contractFound == (entities.Contrato{}) {
		response := utils.BuildErrorResponse(utils.ClientNotFound)
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	err := controller.contractService.DeleteContract(contractFound)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
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
		response := utils.BuildErrorResponse(utils.ContractNotFound)
		ctx.JSON(http.StatusNotFound, response)
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
func NewContractController() ContractController {
	return &contractController{
		contractService:      services.NewContractService(),
		contractEventService: services.NewContractEventService(),
		pontoService:         services.NewPointService(),
	}
}
