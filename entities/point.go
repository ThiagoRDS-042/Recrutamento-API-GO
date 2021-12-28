package entities

// Ponto representa a tabela t_ponto no banco de dados.
type Ponto struct {
	Base
	ClienteID  string   `json:"cliente_id" gorm:"not null"`
	EnderecoID string   `json:"endereco_id" gorm:"not null"`
	Cliente    Cliente  `json:"-" gorm:"foreignKey:ClienteID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Endereco   Endereco `json:"-" gorm:"foreignKey:EnderecoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
