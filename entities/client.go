package entities

import (
	"time"

	"gorm.io/gorm"
)

// ClientType representa o type ClientType.
type ClientType string

// Constantes que representam os tipos de clientes validos.
const (
	JURIDICO ClientType = "juridico"
	ESPECIAL ClientType = "especial"
	FISICO   ClientType = "fisico"
)

// Cliente representa a tabela t_cliente no banco de dados.
type Cliente struct {
	ID              string         `json:"-" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();not null"`
	Nome            string         `json:"nome" gorm:"type:text;size:128;not null;unique"`
	Tipo            ClientType     `json:"tipo" gorm:"not null"`
	DataCriacao     time.Time      `json:"-" gorm:"not null"`
	DataAtualizacao time.Time      `json:"-" gorm:"not null"`
	DataRemocao     gorm.DeletedAt `json:"-" gorm:"index"`
}

// IsValidClientName verifica se o nome do cliente é valido.
func IsValidClientName(clientName string) bool {
	if len(clientName) < 3 || len(clientName) > 128 {
		return false
	}

	return true
}

// IsValidClientType verifica se o tipo de cliente e valido.
func IsValidClientType(clientType ClientType) bool {
	if clientType != FISICO && clientType != ESPECIAL && clientType != JURIDICO {
		return false
	}

	return true
}
