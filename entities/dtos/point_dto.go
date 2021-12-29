package dtos

import "github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"

// PointCreateDTO representa o modelo usado para cadastrar pontos.
type PointCreateDTO struct {
	ClienteID  string `json:"cliente_id" form:"cliente_id" binding:"required"`
	EnderecoID string `json:"endereco_id" form:"endereco_id" binding:"required"`
}

// PointUpdateDTO representa o modelo usado para atualizar pontos.
type PointUpdateDTO struct {
	Base
	ClienteID  string `json:"cliente_id" form:"cliente_id"`
	EnderecoID string `json:"endereco_id" form:"endereco_id"`
}

// PointResponse representa o modelo usado para retornar a resposta da pesquisa dos pontos.
type PointResponse struct {
	ID                 string              `json:"id"`
	ClienteID          string              `json:"cliente_id"`
	ClienteNome        string              `json:"cliente_nome"`
	ClienteTipo        entities.ClientType `json:"cliente_tipo"`
	EnderecoID         string              `json:"endereco_id"`
	EnderecoLogradouro string              `json:"endereco_logradouro"`
	EnderecoBairro     string              `json:"endereco_bairro"`
	EnderecoNumero     int                 `json:"endereco_numero"`
}

// CreatePointResponse cria a responsta modela para a pesquisa de pontos.
func CreatePointResponse(point entities.Ponto) PointResponse {
	pointResponse := PointResponse{
		ID:                 point.ID,
		ClienteID:          point.ClienteID,
		ClienteNome:        point.Cliente.Nome,
		ClienteTipo:        point.Cliente.Tipo,
		EnderecoID:         point.EnderecoID,
		EnderecoLogradouro: point.Endereco.Logradouro,
		EnderecoBairro:     point.Endereco.Bairro,
		EnderecoNumero:     point.Endereco.Numero}

	return pointResponse
}
