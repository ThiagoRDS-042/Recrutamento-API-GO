package dtos

// AddressCreateDTO representa o modelo usado para cadastrar endereços.
type AddressCreateDTO struct {
	Logradouro string `json:"logradouro" form:"logradouro" binding:"required,min=3,max=128"`
	Bairro     string `json:"bairro" form:"bairro" binding:"required,min=3,max=128"`
	Numero     int    `json:"numero" gorm:"type:smallint" binding:"required"`
}

// AddressUpdateDTO representa o modelo usado para atualizar endereços.
type AddressUpdateDTO struct {
	Base
	Logradouro string `json:"logradouro" form:"logradouro"`
	Bairro     string `json:"bairro" form:"bairro"`
	Numero     int    `json:"numero" gorm:"type:smallint"`
}
