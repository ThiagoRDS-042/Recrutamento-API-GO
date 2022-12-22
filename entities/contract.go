package entities

import "gorm.io/gorm"

// ContractState representa o type ContractState.
type ContractState string

// Constantes que representam os estados do contrato.
const (
	VIGOR      ContractState = "Em vigor"
	DESATIVADO ContractState = "Desativado Temporario"
	CANCELADO  ContractState = "Cancelado"
)

// Contrato representa a tabela t_contrato no banco de dados.
type Contrato struct {
	Base
	Estado      ContractState  `json:"-" gorm:"not null"`
	PontoID     string         `json:"ponto_id" gorm:"type:uuid;not null"`
	Ponto       Ponto          `json:"-" gorm:"foreignKey:PontoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DataRemocao gorm.DeletedAt `json:"-" gorm:"index"`
}
