package entities

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
	Nome string     `json:"nome" gorm:"type:text;size:128;not null;unique"`
	Tipo ClientType `json:"tipo" gorm:"not null"`
}

// IsValidTextLenght verifica se o tamanho do texto é valido.
func IsValidTextLenght(text string) bool {
	if len(text) < 3 || len(text) > 128 {
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
