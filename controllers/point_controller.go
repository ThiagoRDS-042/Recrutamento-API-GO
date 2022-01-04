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

// PointController representa o contracto de pointController.
type PointController interface {
	CreatePoint(ctx *gin.Context)
	DeletePoint(ctx *gin.Context)
	FindPoints(ctx *gin.Context)
}

type pointController struct {
	pointService    services.PointService
	clientService   services.ClientService
	addressService  services.AddressService
	contractService services.ContractService
}

// CreatePoint godoc
// @Summary cria um novo ponto
// @Description rota para o cadastro de novos pontos
// @Tags point
// @Accept json
// @Produce json
// @Param point body entities.Ponto true "Criar Novo Ponto"
// @Success 201 {object} entities.Ponto
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /pontos [post]
func (controller *pointController) CreatePoint(ctx *gin.Context) {
	pointDTO := dtos.PointCreateDTO{}

	if err := ctx.ShouldBindJSON(&pointDTO); err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	clientExists := controller.clientService.FindClientByID(pointDTO.ClienteID)
	if clientExists == (entities.Cliente{}) {
		response := utils.BuildErrorResponse(utils.ClientNotFound)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	addressExists := controller.addressService.FindAddressByID(pointDTO.EnderecoID)
	if addressExists == (entities.Endereco{}) {
		response := utils.BuildErrorResponse(utils.AddressNotFound)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	pointAlreadyExists := controller.pointService.FindPointByClientIDAndAddressID(
		pointDTO.ClienteID, pointDTO.EnderecoID)

	switch {
	case pointAlreadyExists.DataRemocao.Valid:
		pointUpdateDTO := dtos.PointUpdateDTO{}
		err := smapping.FillStruct(&pointUpdateDTO, smapping.MapFields(&pointDTO))
		if err != nil {
			log.Fatalf("failed to map: %v", err)
		}

		pointUpdateDTO.ID = pointAlreadyExists.ID
		point, err := controller.pointService.UpdatePoint(pointUpdateDTO)
		if err != nil {
			response := utils.BuildErrorResponse(err.Error())
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		ctx.JSON(http.StatusCreated, point)

	case (pointAlreadyExists != entities.Ponto{}):
		response := utils.BuildErrorResponse(utils.PointAlreadyExists)
		ctx.JSON(http.StatusConflict, response)

	default:
		point, err := controller.pointService.CreatePoint(pointDTO)
		if err != nil {
			response := utils.BuildErrorResponse(err.Error())
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		ctx.JSON(http.StatusCreated, point)
	}
}

// DeletePoint godoc
// @Summary deleta o ponto
// @Description rota para a exclusão do ponto pelo id
// @Tags point
// @Accept json
// @Produce json
// @Param id path string true "id do ponto"
// @Success 204 "No Content"
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /ponto/{id} [delete]
func (controller *pointController) DeletePoint(ctx *gin.Context) {
	pointID := ctx.Param("id")

	pointFound := controller.pointService.FindPointByID(pointID)

	if pointFound == (entities.Ponto{}) {
		response := utils.BuildErrorResponse(utils.PointNotFound)
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	err := controller.pointService.DeletePoint(pointFound)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err = controller.contractService.DeleteContractByPontoID(pointFound.ID)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	ctx.JSON(http.StatusNoContent, entities.Cliente{})
}

// FindPoints godoc
// @Summary lista os pontos existentes
// @Description rota para a listagem de todos os pontos existentes no banco de dados
// @Tags point
// @Accept json
// @Produce json
// @Param cliente_id query string false "id do cliente"
// @Param endereco_id query string false "id do endereço"
// @Success 200 {object} []dtos.PointResponse
// @Failure 404 {object} utils.Response
// @Router /pontos [get]
func (controller *pointController) FindPoints(ctx *gin.Context) {
	clientID := ctx.Query("cliente_id")
	addressID := ctx.Query("endereco_id")

	points := controller.pointService.FindPoints(clientID, addressID)

	if len(points) == 0 {
		response := utils.BuildErrorResponse(utils.PointNotFound)
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	pointsResponse := []dtos.PointResponse{}

	for _, point := range points {
		pointsResponse = append(pointsResponse, dtos.CreatePointResponse(point))
	}

	response := map[string][]dtos.PointResponse{
		"dados": pointsResponse,
	}

	ctx.JSON(http.StatusOK, response)
}

// NewPointController cria uma nova isnancia de PointController.
func NewPointController(pointService services.PointService, clientService services.ClientService, addressService services.AddressService, contractService services.ContractService) PointController {

	return &pointController{
		pointService:    pointService,
		clientService:   clientService,
		addressService:  addressService,
		contractService: contractService,
	}
}
