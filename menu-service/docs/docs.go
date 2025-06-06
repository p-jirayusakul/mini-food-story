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
                "description": "Search menu items with filters",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Menu"
                ],
                "summary": "Search menu items",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "X-Session-Id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "pageNumber",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page size",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Search by name",
                        "name": "search",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by category IDs e.g. 333,444,555",
                        "name": "categoryID",
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
                        "name": "orderType",
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
                                            "$ref": "#/definitions/food-story_menu-service_internal_domain.SearchProductResult"
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
        "/category": {
            "get": {
                "description": "Get list of all available product categories",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Category"
                ],
                "summary": "Get list of categories",
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
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/food-story_menu-service_internal_domain.Category"
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
            "get": {
                "description": "Get menu item details by product ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Menu"
                ],
                "summary": "Get menu item by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "X-Session-Id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
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
                                            "$ref": "#/definitions/food-story_menu-service_internal_domain.Product"
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
        }
    },
    "definitions": {
        "food-story_menu-service_internal_domain.Category": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "1921144250070732800"
                },
                "name": {
                    "type": "string",
                    "example": "ขนม"
                },
                "nameEN": {
                    "type": "string",
                    "example": "Dessert"
                }
            }
        },
        "food-story_menu-service_internal_domain.Product": {
            "type": "object",
            "properties": {
                "categoryID": {
                    "type": "string",
                    "example": "1921143886227443712"
                },
                "categoryName": {
                    "type": "string",
                    "example": "อาหาร"
                },
                "categoryNameEN": {
                    "type": "string",
                    "example": "Food"
                },
                "description": {
                    "type": "string",
                    "example": "lorem ipsum"
                },
                "id": {
                    "type": "string",
                    "example": "1921144250070732800"
                },
                "imageURL": {
                    "type": "string",
                    "example": "https://example.com/image.jpg"
                },
                "isAvailable": {
                    "type": "boolean",
                    "example": true
                },
                "name": {
                    "type": "string",
                    "example": "ข้าวมันไก่"
                },
                "nameEN": {
                    "type": "string",
                    "example": "Chicken rice"
                },
                "price": {
                    "type": "number",
                    "example": 100
                }
            }
        },
        "food-story_menu-service_internal_domain.SearchProductResult": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/food-story_menu-service_internal_domain.Product"
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
