package controllers

import (
	"net/http"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/services"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/utils"
	"github.com/gin-gonic/gin"
)

// ContractEventController representa o contracto de ContractEventController.
type ContractEventController interface {
	CreateContractEvent(ctx *gin.Context)
	FindContractEventsByContractID(ctx *gin.Context)
}

type contractEventController struct {
	contractEventService services.ContractEventService
	contractService      services.ContractService
}

func (controller *contractEventController) CreateContractEvent(ctx *gin.Context) {
	contractEventDTO := dtos.ContratoEventCreateDTO{}

	if err := ctx.ShouldBindJSON(&contractEventDTO); err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	contractExists := controller.contractService.FindContractByID(contractEventDTO.ContratoID)
	if contractExists == (entities.Contrato{}) {
		response := utils.BuildErrorResponse(utils.ContractNotFound)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	contractEvent, err := controller.contractEventService.CreateContractEvent(contractEventDTO)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	ctx.JSON(http.StatusCreated, contractEvent)
}

func (controller *contractEventController) FindContractEventsByContractID(ctx *gin.Context) {
	contractID := ctx.Param("id")

	contractEvents := controller.contractEventService.FindContractEventsByContractID(contractID)

	if len(contractEvents) == 0 {
		response := utils.BuildErrorResponse(utils.HitoryOfContractNotFound)
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	contractEventsResponse := []dtos.ContractEventResponse{}

	for _, contractEvent := range contractEvents {
		contractEventsResponse = append(contractEventsResponse, dtos.CreateContractEventResponse(contractEvent))
	}

	response := map[string][]dtos.ContractEventResponse{
		"dados": contractEventsResponse,
	}

	ctx.JSON(http.StatusOK, response)
}

// NewContractEventController cria uma nova isnancia de ContractEventController.
func NewContractEventController() ContractEventController {
	return &contractEventController{
		contractEventService: services.NewContractEventService(),
		contractService:      services.NewContractService(),
	}
}
