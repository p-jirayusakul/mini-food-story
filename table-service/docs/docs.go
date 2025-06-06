// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Search tables by filters like number of people, table number, seats, and status",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Table"
                ],
                "summary": "Search table availability",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Number of people",
                        "name": "numberOfPeople",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Search by table number",
                        "name": "search",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by seats",
                        "name": "seats",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "Filter by status codes",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page number for pagination",
                        "name": "pageNumber",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page size for pagination",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Order by field (id, tableNumber, seats, status)",
                        "name": "orderBy",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Order direction (asc, desc)",
                        "name": "orderByType",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/food-story_table-service_internal_domain.SearchTablesResult"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Create a new table with specified number and seats",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Table"
                ],
                "summary": "Create new table",
                "parameters": [
                    {
                        "description": "Table details",
                        "name": "table",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_adapter_http.Table"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/internal_adapter_http.createResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/quick-search": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Quickly search for available tables based on number of people",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Table"
                ],
                "summary": "Quick search for available tables",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Number of people required",
                        "name": "numberOfPeople",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page number for pagination",
                        "name": "pageNumber",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page size for pagination",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Order by field (id, tableNumber, seats, status)",
                        "name": "orderBy",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Order direction (asc, desc)",
                        "name": "orderByType",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/food-story_table-service_internal_domain.SearchTablesResult"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Create a new session for a table with specified number of people",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Table"
                ],
                "summary": "Create new table session",
                "parameters": [
                    {
                        "description": "Table session details",
                        "name": "table",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_adapter_http.TableSession"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/internal_adapter_http.createSessionResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/session/current": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get details of the current active table session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Table"
                ],
                "summary": "Get current table session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "X-Session-Id",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.CurrentTableSession"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/status": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get list of all available table statuses",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Table"
                ],
                "summary": "Get list of table status",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.SuccessResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/food-story_table-service_internal_domain.Status"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/{id}": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Update table number and seats for existing table",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Table"
                ],
                "summary": "Update table details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Table ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Table details",
                        "name": "table",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_adapter_http.Table"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/middleware.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/{id}/status": {
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Update status for existing table",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Table"
                ],
                "summary": "Update table status",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Table ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Table status details",
                        "name": "status",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_adapter_http.updateTableStatus"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/middleware.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "food-story_table-service_internal_domain.SearchTablesResult": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/food-story_table-service_internal_domain.Table"
                    }
                },
                "totalItems": {
                    "type": "integer",
                    "example": 10
                },
                "totalPages": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "food-story_table-service_internal_domain.Status": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string",
                    "example": "ORDERED"
                },
                "id": {
                    "type": "string",
                    "example": "1921144250070732800"
                },
                "name": {
                    "type": "string",
                    "example": "สั่งอาหารแล้ว"
                },
                "nameEN": {
                    "type": "string",
                    "example": "Ordered"
                }
            }
        },
        "food-story_table-service_internal_domain.Table": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "1923564209627467776"
                },
                "seats": {
                    "type": "integer",
                    "example": 5
                },
                "status": {
                    "type": "string",
                    "example": "สั่งอาหารแล้ว"
                },
                "statusEN": {
                    "type": "string",
                    "example": "Ordered"
                },
                "tableNumber": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "internal_adapter_http.Table": {
            "type": "object",
            "required": [
                "seats",
                "tableNumber"
            ],
            "properties": {
                "seats": {
                    "type": "integer",
                    "minimum": 1,
                    "example": 5
                },
                "tableNumber": {
                    "type": "integer",
                    "minimum": 1,
                    "example": 1
                }
            }
        },
        "internal_adapter_http.TableSession": {
            "type": "object",
            "required": [
                "numberOfPeople",
                "tableID"
            ],
            "properties": {
                "numberOfPeople": {
                    "type": "integer",
                    "minimum": 1,
                    "example": 3
                },
                "tableID": {
                    "type": "string",
                    "example": "1923564209627467776"
                }
            }
        },
        "internal_adapter_http.createResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "1923564209627467776"
                }
            }
        },
        "internal_adapter_http.createSessionResponse": {
            "type": "object",
            "properties": {
                "url": {
                    "type": "string",
                    "example": "http://localhost:3000?s=tYlC7uRGdaIT0-wDvCngv9qEWwM4Tdg4Jzoywe00fK2WhDqDHKCzqtybFDVgCPIGv1_isM3aXdb16KNGD4E-q2kewBRaVXZ5N7vgdi46Tc5po6_ZCHpFAT-Ei3xkKT0dL_f3Ruoiz9IzsnBhxlexXTEhKN9myrECimbKaDI="
                }
            }
        },
        "internal_adapter_http.updateTableStatus": {
            "type": "object",
            "required": [
                "statusID"
            ],
            "properties": {
                "statusID": {
                    "type": "string",
                    "example": "1919968486671519744"
                }
            }
        },
        "middleware.ErrorResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "message": {
                    "type": "string",
                    "example": "something went wrong"
                },
                "status": {
                    "type": "string",
                    "example": "error"
                }
            }
        },
        "middleware.SuccessResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "message": {
                    "type": "string",
                    "example": "do something completed"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "model.CurrentTableSession": {
            "type": "object",
            "properties": {
                "orderID": {
                    "type": "string",
                    "example": "1922535048335069184"
                },
                "sessionID": {
                    "type": "string",
                    "example": "a9213539-b135-42cc-b714-60cfd1b099ec"
                },
                "startedAt": {
                    "type": "string",
                    "example": "2025-05-23T11:59:50.010316+07:00"
                },
                "status": {
                    "type": "string",
                    "example": "active"
                },
                "tableID": {
                    "type": "string",
                    "example": "1920153361642950656"
                },
                "tableNumber": {
                    "type": "integer",
                    "example": 1
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
