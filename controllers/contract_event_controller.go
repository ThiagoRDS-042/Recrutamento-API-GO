package controllers

import (
	"net/http"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/services"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/utils"
	"github.com/gin-gonic/gin"
)

// ContractEventController representa o contracto de ContractEventController.
type ContractEventController interface {
	FindContractEventsByContractID(ctx *gin.Context)
}

type contractEventController struct {
	contractEventService services.ContractEventService
	contractService      services.ContractService
}

// FindContractEventsByContractID godoc
// @Summary pesquisa de evento de contrato
// @Description rota para a pesquisa do hitorico de evento de contrato pelo id do contrato
// @Tags contractEvent
// @Accept json
// @Produce json
// @Param id path string true "id do contrato"
// @Success 200 {object} []dtos.ContractEventResponse
// @Failure 404 {object} utils.Response
// @Router /contrato/{id}/historico [get]
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
func NewContractEventController(contractEventService services.ContractEventService, contractService services.ContractService) ContractEventController {
	return &contractEventController{
		contractEventService: contractEventService,
		contractService:      contractService,
	}
}
