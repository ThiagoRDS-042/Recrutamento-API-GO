package dtos

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
)

// ContractCreateDTO representa o modelo usado para cadastrar contratos.
type ContractCreateDTO struct {
	Estado  entities.ContractState `json:"estado" form:"estado"`
	PontoID string                 `json:"ponto_id" form:"ponto_id" binding:"required"`
}

// ContractUpdateDTO representa o modelo usado para atualizar contratos.
type ContractUpdateDTO struct {
	Base
	Estado entities.ContractState `json:"estado" form:"estado" binding:"required,eq=Em vigor|eq=Desativado Temporario|eq=Cancelado"`
}

// ContractResponse representa o modelo usado para retornar a resposta da pesquisa dos contratos.
type ContractResponse struct {
	ID                 string              `json:"id"`
	ClienteID          string              `json:"cliente_id"`
	ClienteNome        string              `json:"cliente_nome"`
	ClienteTipo        entities.ClientType `json:"cliente_tipo"`
	EnderecoID         string              `json:"endereco_id"`
	EnderecoLogradouro string              `json:"endereco_logradouro"`
	EnderecoBairro     string              `json:"endereco_bairro"`
	EnderecoNumero     int                 `json:"endereco_numero"`
}

// IsAuthorized verifica se a alteração de estado do contrato é valida.
func IsAuthorized(oldState entities.ContractState, newState entities.ContractState) bool {
	switch {
	case oldState == newState:
		return false
	case oldState == entities.CANCELADO:
		return false
	case oldState == entities.VIGOR && newState != entities.DESATIVADO:
		return false
	default:
		return true
	}
}

// CreateContractResponse cria a responsta modelada para a pesquisa de contratos.
func CreateContractResponse(contrat entities.Contrato) ContractResponse {
	contractResponse := ContractResponse{
		ID:                 contrat.ID,
		ClienteID:          contrat.Ponto.ClienteID,
		ClienteNome:        contrat.Ponto.Cliente.Nome,
		ClienteTipo:        contrat.Ponto.Cliente.Tipo,
		EnderecoID:         contrat.Ponto.EnderecoID,
		EnderecoLogradouro: contrat.Ponto.Endereco.Logradouro,
		EnderecoBairro:     contrat.Ponto.Endereco.Bairro,
		EnderecoNumero:     contrat.Ponto.Endereco.Numero,
	}

	return contractResponse
}
