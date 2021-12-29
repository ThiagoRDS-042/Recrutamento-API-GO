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
	pointService   services.PointService
}

func (controller *addressController) CreateAddress(ctx *gin.Context) {
	addressDTO := dtos.AddressCreateDTO{}

	if err := ctx.ShouldBindJSON(&addressDTO); err != nil {
		response := utils.BuildErrorResponse("createError: " + err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	addressAlreadyExists := controller.addressService.FindAddressByFields(
		addressDTO.Logradouro, addressDTO.Bairro, addressDTO.Numero)

	switch {
	case addressAlreadyExists.DataRemocao.Valid:
		addressUpdateDTO := dtos.AddressUpdateDTO{}
		err := smapping.FillStruct(&addressUpdateDTO, smapping.MapFields(&addressDTO))
		if err != nil {
			log.Fatalf("failed to map: %v", err)
		}

		addressUpdateDTO.ID = addressAlreadyExists.ID

		client, err := controller.addressService.UpdateAddress(addressUpdateDTO)
		if err != nil {
			response := utils.BuildErrorResponse(err.Error())
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		ctx.JSON(http.StatusCreated, client)

	case (addressAlreadyExists != entities.Endereco{}):
		response := utils.BuildErrorResponse(utils.AddressAlreadyExists)
		ctx.JSON(http.StatusConflict, response)

	default:
		address, err := controller.addressService.CreateAddress(addressDTO)
		if err != nil {
			response := utils.BuildErrorResponse(err.Error())
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		ctx.JSON(http.StatusCreated, address)
	}
}

func (controller *addressController) UpdateAddress(ctx *gin.Context) {
	addressDTO := dtos.AddressUpdateDTO{}

	if err := ctx.ShouldBindJSON(&addressDTO); err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	addressID := ctx.Param("id")

	addressFound := controller.addressService.FindAddressByID(addressID)

	if addressFound == (entities.Endereco{}) {
		response := utils.BuildErrorResponse(utils.AddressNotFound)
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	addressDTO.ID = addressID

	if addressDTO.Logradouro == "" {
		addressDTO.Logradouro = addressFound.Logradouro
	} else {
		if !dtos.IsValidTextLenght(addressDTO.Logradouro) {
			response := utils.BuildErrorResponse("logradouro: " + utils.InvalidNumberOfCaracter)
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
	}

	if addressDTO.Bairro == "" {
		addressDTO.Bairro = addressFound.Bairro
	} else {
		if !dtos.IsValidTextLenght(addressDTO.Bairro) {
			response := utils.BuildErrorResponse("bairro: " + utils.InvalidNumberOfCaracter)
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
	}

	if addressDTO.Numero == 0 {
		addressDTO.Numero = addressFound.Numero
	}

	addressAlreadyExists := controller.addressService.FindAddressByFields(
		addressDTO.Logradouro, addressDTO.Bairro, addressDTO.Numero)

	if (addressAlreadyExists != entities.Endereco{}) && (addressFound.ID != addressAlreadyExists.ID) {
		response := utils.BuildErrorResponse(utils.AddressAlreadyExists)
		ctx.JSON(http.StatusConflict, response)
		return
	}

	address, err := controller.addressService.UpdateAddress(addressDTO)
	if err != nil {
		response := utils.BuildErrorResponse("updateError: " + err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	ctx.JSON(http.StatusOK, address)
}

func (controller *addressController) FindAddressByID(ctx *gin.Context) {
	addressID := ctx.Param("id")

	addressFound := controller.addressService.FindAddressByID(addressID)

	if addressFound == (entities.Endereco{}) {
		response := utils.BuildErrorResponse(utils.AddressNotFound)
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, addressFound)
}

func (controller *addressController) DeleteAddress(ctx *gin.Context) {
	addressID := ctx.Param("id")

	addressFound := controller.addressService.FindAddressByID(addressID)

	if addressFound == (entities.Endereco{}) {
		response := utils.BuildErrorResponse(utils.AddressNotFound)
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	err := controller.addressService.DeleteAddress(addressFound)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err = controller.pointService.DeletePointsByAddressID(addressFound.ID)
	if err != nil {
		response := utils.BuildErrorResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	ctx.JSON(http.StatusNoContent, entities.Endereco{})
}

func (controller *addressController) FindAddress(ctx *gin.Context) {
	addressNeighborhood := ctx.Query("bairro")
	addressStreet := ctx.Query("logradouro")
	addressNumber := ctx.Query("numero")

	addresses := controller.addressService.FindAddresses(
		addressStreet, addressNeighborhood, addressNumber)

	if len(addresses) == 0 {
		response := utils.BuildErrorResponse(utils.AddressNotFound)
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	response := map[string][]entities.Endereco{
		"dados": addresses,
	}

	ctx.JSON(http.StatusOK, response)
}

// NewAddressController cria uma nova isnancia de AddressController.
func NewAddressController() AddressController {
	return &addressController{
		addressService: services.NewAddressService(),
		pointService:   services.NewPointService(),
	}
}
