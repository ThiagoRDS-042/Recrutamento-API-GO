package entities

import (
	"time"
)

// Base utilizada para representar aos capos genericos de todas as entidades do banco de dados.
type Base struct {
	ID              string    `json:"-" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();not null"`
	DataCriacao     time.Time `json:"-" gorm:"not null"`
	DataAtualizacao time.Time `json:"-" gorm:"not null"`
}
