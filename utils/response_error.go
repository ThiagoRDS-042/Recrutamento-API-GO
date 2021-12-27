package utils

import "strings"

// Response usada como corpo estaico para a resposta de error em json.
type Response struct {
	Message interface{} `json:"message"`
}

// BuildErrorResponse injeta um valor dinamico na resposta de error em json.
func BuildErrorResponse(message string) Response {
	splitedSMessage := strings.Split(message, "\n")

	return Response{
		Message: splitedSMessage,
	}
}
