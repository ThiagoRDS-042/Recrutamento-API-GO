// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Thiago Rodrigues",
            "url": "http://thiagords042/support",
            "email": "thiagords042@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/cliente/{id}": {
            "get": {
                "description": "rota para a pesquisa do cliente pelo id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "client"
                ],
                "summary": "pesquisa o cliente",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id do cliente",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.Cliente"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            },
            "put": {
                "description": "rota para a atualização dos dados do cliente a partir do id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "client"
                ],
                "summary": "atualiza o cliente",
                "parameters": [
                    {
                        "description": "atualizar cliente",
                        "name": "client",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.Cliente"
                        }
                    },
                    {
                        "type": "string",
                        "description": "id do cliente",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.Cliente"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "rota para a exclusão do cliente pelo id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "client"
                ],
                "summary": "deleta o cliente",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id do cliente",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/clientes": {
            "get": {
                "description": "rota para a listagem de todos os clientes existentes no banco de dados",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "client"
                ],
                "summary": "lista os clientes existentes",
                "parameters": [
                    {
                        "type": "string",
                        "description": "tipo de cliente",
                        "name": "tipo",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "nome do cliente",
                        "name": "nome",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.Cliente"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "rota para o cadastro de novos clientes",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "client"
                ],
                "summary": "cria um novo cliente",
                "parameters": [
                    {
                        "description": "Criar Novo Cliente",
                        "name": "client",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.Cliente"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/entities.Cliente"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/contrato/{id}": {
            "get": {
                "description": "rota para a pesquisa do contrato pelo id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "contract"
                ],
                "summary": "pesquisa o contrato",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id do contrato",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.ContractResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            },
            "put": {
                "description": "rota para a atualização dos dados do contrato a partir do id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "contract"
                ],
                "summary": "atualiza o contrato",
                "parameters": [
                    {
                        "enum": [
                            "Em vigor",
                            "Desativado Temporario",
                            "Cancelado"
                        ],
                        "description": "atualizar contrato",
                        "name": "estado",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "type": "string",
                        "description": "id do contrato",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.Contrato"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "rota para a exclusão do contrato pelo id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "contract"
                ],
                "summary": "deleta o contrato",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id do contrato",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/contrato/{id}/historico": {
            "get": {
                "description": "rota para a pesquisa do hitorico de evento de contrato pelo id do contrato",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "contractEvent"
                ],
                "summary": "pesquisa de evento de contrato",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id do contrato",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dtos.ContractEventResponse"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/contratos": {
            "get": {
                "description": "rota para a listagem de todos os contratos existentes no banco de dados",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "contract"
                ],
                "summary": "lista os contratos existentes",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id do cliente",
                        "name": "cliente_id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "id do endereço",
                        "name": "endereco_id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dtos.ContractResponse"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "rota para o cadastro de novos contratos a partir do id do ponto",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "contract"
                ],
                "summary": "cria um novo contrato",
                "parameters": [
                    {
                        "description": "Criar Novo Contrato",
                        "name": "contract",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.Contrato"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/entities.Contrato"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/endereco/{id}": {
            "get": {
                "description": "rota para a pesquisa do endereço pelo id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "address"
                ],
                "summary": "pesquisa o endereço",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id do endereço",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.Endereco"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            },
            "put": {
                "description": "rota para a atualização dos dados do endereço a partir do id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "address"
                ],
                "summary": "atualiza o endereço",
                "parameters": [
                    {
                        "description": "atualizar endereço",
                        "name": "address",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.Endereco"
                        }
                    },
                    {
                        "type": "string",
                        "description": "id do endereço",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.Endereco"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "rota para a exclusão do endereço pelo id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "address"
                ],
                "summary": "deleta o endereço",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id do endereço",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/enderecos": {
            "get": {
                "description": "rota para a listagem de todos os endereços existentes no banco de dados",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "address"
                ],
                "summary": "lista os endreços existentes",
                "parameters": [
                    {
                        "type": "string",
                        "description": "logradouro",
                        "name": "logradouro",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "bairro",
                        "name": "bairro",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "numero da casa",
                        "name": "numero",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.Endereco"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "rota para o cadastro de novos endereços",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "address"
                ],
                "summary": "cria um novo endereço",
                "parameters": [
                    {
                        "description": "Criar Novo Endereço",
                        "name": "address",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.Endereco"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/entities.Endereco"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/ponto/{id}": {
            "delete": {
                "description": "rota para a exclusão do ponto pelo id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "point"
                ],
                "summary": "deleta o ponto",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id do ponto",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/pontos": {
            "get": {
                "description": "rota para a listagem de todos os pontos existentes no banco de dados",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "point"
                ],
                "summary": "lista os pontos existentes",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id do cliente",
                        "name": "cliente_id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "id do endereço",
                        "name": "endereco_id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dtos.PointResponse"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "rota para o cadastro de novos pontos",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "point"
                ],
                "summary": "cria um novo ponto",
                "parameters": [
                    {
                        "description": "Criar Novo Ponto",
                        "name": "point",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.Ponto"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/entities.Ponto"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dtos.ContractEventResponse": {
            "type": "object",
            "properties": {
                "data_evento": {
                    "type": "string"
                },
                "estado_antigo": {
                    "type": "string"
                },
                "estado_novo": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "dtos.ContractResponse": {
            "type": "object",
            "properties": {
                "cliente_id": {
                    "type": "string"
                },
                "cliente_nome": {
                    "type": "string"
                },
                "cliente_tipo": {
                    "type": "string"
                },
                "endereco_bairro": {
                    "type": "string"
                },
                "endereco_id": {
                    "type": "string"
                },
                "endereco_logradouro": {
                    "type": "string"
                },
                "endereco_numero": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "dtos.PointResponse": {
            "type": "object",
            "properties": {
                "cliente_id": {
                    "type": "string"
                },
                "cliente_nome": {
                    "type": "string"
                },
                "cliente_tipo": {
                    "type": "string"
                },
                "endereco_bairro": {
                    "type": "string"
                },
                "endereco_id": {
                    "type": "string"
                },
                "endereco_logradouro": {
                    "type": "string"
                },
                "endereco_numero": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "entities.Cliente": {
            "type": "object",
            "properties": {
                "nome": {
                    "type": "string"
                },
                "tipo": {
                    "type": "string"
                }
            }
        },
        "entities.Contrato": {
            "type": "object",
            "properties": {
                "ponto_id": {
                    "type": "string"
                }
            }
        },
        "entities.Endereco": {
            "type": "object",
            "properties": {
                "bairro": {
                    "type": "string"
                },
                "logradouro": {
                    "type": "string"
                },
                "numero": {
                    "type": "integer"
                }
            }
        },
        "entities.Ponto": {
            "type": "object",
            "properties": {
                "cliente_id": {
                    "type": "string"
                },
                "endereco_id": {
                    "type": "string"
                }
            }
        },
        "utils.Response": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:2222",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "API Recrutamento",
	Description: "Uma API básica para cadastro de clientes, endereços e contratos.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
