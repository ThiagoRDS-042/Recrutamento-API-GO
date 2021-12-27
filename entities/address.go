package entities

// Endereco representa a tabela t_endereco no banco de dados.
type Endereco struct {
	Base
	Logradouro string `json:"logradouro" gorm:"type:text;size:128;not null"`
	Bairro     string `json:"bairro" gorm:"type:text;size:128;not null"`
	Numero     int    `json:"numero" gorm:"type:smallint"`
}
