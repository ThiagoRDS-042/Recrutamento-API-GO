package controllers

import (
	"net/http"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities/dtos"
	services "github.com/ThiagoRDS-042/Recrutamento-API-GO/services/address_service"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/utils"
	"github.com/gin-gonic/gin"
)

// AddressController representa o contracto de AddressController.
type AddressController interface {
	CreateAddress(ctx *gin.Context)
	UpdateAddress(ctx *gin.Context)
	FindAddressByID(ctx *gin.Context)
	DeleteAddress(ctx *gin.Context)
	FindAddress(ctx *gin.Context)
}

type addressController struct {
	addressService services.AddressService
}

// CreateAddress godoc
// @Summary cria um novo endereço
// @Description rota para o cadastro de novos endereços
// @Tags address
// @Accept json
// @Produce json
// @Param address body entities.Endereco true "Criar Novo Endereço"
// @Success 201 {object} entities.Endereco
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /enderecos [post]
func (controller *addressController) CreateAddress(ctx *gin.Context) {
	addressDTO := dtos.AddressCreateDTO{}

	if err := ctx.ShouldBindJSON(&addressDTO); err != nil {
		response := utils.NewResponse(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	address, responseError := controller.addressService.CreateAddress(addressDTO)
	if responseError != (utils.ResponseError{}) {
		response := utils.NewResponse(responseError.Message)
		ctx.AbortWithStatusJSON(responseError.StatusCode, response)
		return
	}

	ctx.JSON(http.StatusCreated, address)
}

// UpdateAddress godoc
// @Summary atualiza o endereço
// @Description rota para a atualização dos dados do endereço a partir do id
// @Tags address
// @Accept json
// @Produce json
// @Param address body entities.Endereco true "atualizar endereço"
// @Param id path string true "id do endereço"
// @Success 200 {object} entities.Endereco
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /endereco/{id} [put]
func (controller *addressController) UpdateAddress(ctx *gin.Context) {
	addressDTO := dtos.AddressUpdateDTO{}

	if err := ctx.ShouldBindJSON(&addressDTO); err != nil {
		response := utils.NewResponse(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	addressID := ctx.Param("id")

	addressDTO.ID = addressID

	address, responseError := controller.addressService.UpdateAddress(addressDTO)
	if responseError != (utils.ResponseError{}) {
		response := utils.NewResponse(responseError.Message)
		ctx.AbortWithStatusJSON(responseError.StatusCode, response)
		return
	}

	ctx.JSON(http.StatusOK, address)
}

// FindAddressByID godoc
// @Summary pesquisa o endereço
// @Description rota para a pesquisa do endereço pelo id
// @Tags address
// @Accept json
// @Produce json
// @Param id path string true "id do endereço"
// @Success 200 {object} entities.Endereco
// @Failure 404 {object} utils.Response
// @Router /endereco/{id} [get]
func (controller *addressController) FindAddressByID(ctx *gin.Context) {
	addressID := ctx.Param("id")

	addressFound := controller.addressService.FindAddressByID(addressID)

	if addressFound == (entities.Endereco{}) {
		response := utils.NewResponse(utils.AddressNotFound)
		ctx.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, addressFound)
}

// DeleteAddress godoc
// @Summary deleta o endereço
// @Description rota para a exclusão do endereço pelo id
// @Tags address
// @Accept json
// @Produce json
// @Param id path string true "id do endereço"
// @Success 204 "No Content"
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /endereco/{id} [delete]
func (controller *addressController) DeleteAddress(ctx *gin.Context) {
	addressID := ctx.Param("id")

	responseError := controller.addressService.DeleteAddressByID(addressID)
	if responseError != (utils.ResponseError{}) {
		response := utils.NewResponse(responseError.Message)
		ctx.AbortWithStatusJSON(responseError.StatusCode, response)
		return
	}

	ctx.JSON(http.StatusNoContent, entities.Endereco{})
}

// FindAddress godoc
// @Summary lista os endreços existentes
// @Description rota para a listagem de todos os endereços existentes no banco de dados
// @Tags address
// @Accept json
// @Produce json
// @Param logradouro query string false "logradouro"
// @Param bairro query string false "bairro"
// @Param numero query string false "numero da casa"
// @Success 200 {object} []entities.Endereco
// @Failure 404 {object} utils.Response
// @Router /enderecos [get]
func (controller *addressController) FindAddress(ctx *gin.Context) {
	addressNeighborhood := ctx.Query("bairro")
	addressStreet := ctx.Query("logradouro")
	addressNumber := ctx.Query("numero")

	addresses := controller.addressService.FindAddresses(
		addressStreet, addressNeighborhood, addressNumber)

	if len(addresses) == 0 {
		response := utils.NewResponse(utils.AddressNotFound)
		ctx.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	response := map[string][]entities.Endereco{
		"dados": addresses,
	}

	ctx.JSON(http.StatusOK, response)
}

// NewAddressController cria uma nova isnancia de AddressController.
func NewAddressController(addressService services.AddressService) AddressController {
	return &addressController{
		addressService: addressService,
	}
}
