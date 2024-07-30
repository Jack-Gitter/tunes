// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Jack Gitter",
            "email": "jack.a.gitter@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/users/current/followers/": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Gets the current users followers",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Gets the current users followers",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Pagination Key for follow up responses. This key is a spotify ID",
                        "name": "spotifyID",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.PaginationResponse-array_responses_User-string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users/current/unfollow/{otherUserSpotifyID}": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Unfollowers a user for the currently signed in user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Unfollowers a user for the currently signed in user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User spotify ID",
                        "name": "spotifyID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "User to unfollow spotify ID",
                        "name": "otherUserSpotifyID",
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
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users/{spotifyID}": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Gets a tunes user by their spotifyID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Gets a tunes user by their spotify ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User Spotify ID",
                        "name": "spotifyID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users/{spotifyID}/followers/": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Gets a users followers by their spotify ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Gets a users followers by their spotify ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User spotify ID",
                        "name": "spotifyID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Pagination Key for follow up responses. This key is a spotify ID",
                        "name": "spotifyID",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.PaginationResponse-array_responses_User-string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "responses.PaginationResponse-array_responses_User-string": {
            "type": "object",
            "properties": {
                "dataResponse": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/responses.User"
                    }
                },
                "paginationKey": {
                    "type": "string"
                }
            }
        },
        "responses.Role": {
            "type": "string",
            "enum": [
                "BASIC",
                "MODERATOR",
                "ADMIN"
            ],
            "x-enum-varnames": [
                "BASIC_USER",
                "MODERATOR",
                "ADMIN"
            ]
        },
        "responses.User": {
            "type": "object",
            "properties": {
                "bio": {
                    "type": "string"
                },
                "role": {
                    "$ref": "#/definitions/responses.Role"
                },
                "spotifyID": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "description": "\"Authorization header value\"",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Tunes backend API",
	Description:      "The backend REST API for Tunes",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
