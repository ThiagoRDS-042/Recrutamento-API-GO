package entities

// ContratoEvento representa a tabela t_contrato_evento no banco de dados.
type ContratoEvento struct {
	Base
	EstadoAnterior  ContractState `json:"estado_anterior" gorm:"not null"`
	EstadoPosterior ContractState `json:"estado_posterior" gorm:"not null"`
	ContratoID      string        `json:"contrato_id" gorm:"type:uuid;not null"`
	Contrato        Contrato      `json:"-" gorm:"foreignKey:ContratoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
