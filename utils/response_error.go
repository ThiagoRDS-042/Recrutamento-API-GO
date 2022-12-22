package utils

import "strings"

// Response usada como corpo estatico para a resposta de error em json, que contém mensagem.
type Response struct {
	Message []string `json:"message"`
}

// NewResponse cria uma nova instancia de Response com a mensgaem passada.
func NewResponse(message string) Response {
	splitedMessage := strings.Split(message, "\n")

	return Response{
		Message: splitedMessage,
	}
}

// ResponseError usada como corpo estatico para a resposta de error em json, que contém mensagem e codigo de estado.
type ResponseError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

// NewResponseError cria uma nova instancia de ResponseError com a mensgaem e codigo de estado passado.
func NewResponseError(message string, statusCode int) ResponseError {
	return ResponseError{
		Message:    message,
		StatusCode: statusCode,
	}
}
