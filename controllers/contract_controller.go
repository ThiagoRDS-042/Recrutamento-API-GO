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
	contractService services.ContractService
	pontoService    services.PointService
}

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

		ctx.JSON(http.StatusCreated, contract)
	}
}

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

	ctx.JSON(http.StatusOK, contract)
}

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
		contractService: services.NewContractService(),
		pontoService:    services.NewPointService(),
	}
}
