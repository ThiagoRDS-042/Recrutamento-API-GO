package dtos

// Base utilizada para representar aos capos genericos de todas os dtos da API.
type Base struct {
	ID string `json:"id" form:"id"`
}

// IsValidTextLenght verifica se o tamanho do texto é valido.
func IsValidTextLenght(text string) bool {
	if len(text) < 3 || len(text) > 128 {
		return false
	}

	return true
}
