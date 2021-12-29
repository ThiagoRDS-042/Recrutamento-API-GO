package entities

import "gorm.io/gorm"

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
	Base
	Nome        string         `json:"nome" gorm:"type:text;size:128;not null;unique"`
	Tipo        ClientType     `json:"tipo" gorm:"not null"`
	DataRemocao gorm.DeletedAt `json:"-" gorm:"index"`
}
