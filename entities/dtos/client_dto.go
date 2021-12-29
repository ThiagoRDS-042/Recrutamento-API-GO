package dtos

import "github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"

// ClientCreateDTO representa o modelo usado para cadastrar clientes.
type ClientCreateDTO struct {
	Nome string              `json:"nome" form:"nome" binding:"required,min=3,max=128"`
	Tipo entities.ClientType `json:"tipo" form:"tipo" binding:"required,eq=juridico|eq=fisico|eq=especial"`
}

// ClientUpdateDTO representa o modelo usado para atualizar clientes.
type ClientUpdateDTO struct {
	Base `json:"base" form:"base"`
	Nome string              `json:"nome" form:"nome"`
	Tipo entities.ClientType `json:"tipo" form:"tipo"`
}

// IsValidClientType verifica se o tipo de cliente e valido.
func IsValidClientType(clientType entities.ClientType) bool {
	if clientType != entities.FISICO && clientType != entities.ESPECIAL && clientType != entities.JURIDICO {
		return false
	}

	return true
}
