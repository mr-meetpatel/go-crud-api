// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "email@swagger.io"
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
        "/": {
            "get": {
                "description": "Welcome Message",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Home"
                ],
                "summary": "Welcome Message",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/articles": {
            "get": {
                "description": "Return All Articles aviable in Database",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Article"
                ],
                "summary": "Return All Articles",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Article"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new Article with the input paylod",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Article"
                ],
                "summary": "Create a new Article",
                "parameters": [
                    {
                        "description": "Create article",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.Article"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/main.Article"
                        }
                    }
                }
            }
        },
        "/articles/{id}": {
            "get": {
                "description": "Return single Article by articleId",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Article"
                ],
                "summary": "Return single Article",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Article ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Article"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/main.CustomError"
                        }
                    }
                }
            },
            "put": {
                "description": "Update an Article by article id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Article"
                ],
                "summary": "Update an Article",
                "parameters": [
                    {
                        "description": "Update article",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.Article"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Article ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Article"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/main.CustomError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete an Article by article id",
                "tags": [
                    "Article"
                ],
                "summary": "Delete an Article",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Article ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Article": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "main.CustomError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Articles API",
	Description:      "This is a sample API for managing Articles",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
