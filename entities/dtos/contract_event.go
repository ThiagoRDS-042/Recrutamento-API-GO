package dtos

import (
	"time"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
)

// ContratoEventCreateDTO representa a tabela t_contrato_evento no banco de dados.
type ContratoEventCreateDTO struct {
	EstadoAnterior  entities.ContractState `json:"estado_anterior" form:"estado_anterior" binding:"required"`
	EstadoPosterior entities.ContractState `json:"estado_posterior" form:"estado_posterior" binding:"required"`
	ContratoID      string                 `json:"contrato_id" form:"contrato_id" binding:"required"`
}

// ContractEventResponse representa o modelo usado para retornar a resposta do histórico de alteração de do contrato.
type ContractEventResponse struct {
	ID           string                 `json:"id"`
	DataEvento   time.Time              `json:"data_evento"`
	EstadoAntigo entities.ContractState `json:"estado_antigo"`
	EstadoNovo   entities.ContractState `json:"estado_novo"`
}

// CreateContractEventResponse cria a responsta modelada para a pesquisa do histórico de alteração de do contrato.
func CreateContractEventResponse(contractEvent entities.ContratoEvento) ContractEventResponse {
	contractEventResponse := ContractEventResponse{
		ID:           contractEvent.ID,
		DataEvento:   contractEvent.DataCriacao,
		EstadoAntigo: contractEvent.EstadoAnterior,
		EstadoNovo:   contractEvent.EstadoPosterior,
	}

	return contractEventResponse
}
