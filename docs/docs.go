// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/songs": {
            "get": {
                "description": "Get songs with optional filters (group, artist, album, song) and pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Get list of songs",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Filter by group",
                        "name": "group",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by artist",
                        "name": "artist",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by album",
                        "name": "album",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by song title",
                        "name": "song",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Number of items per page",
                        "name": "pageSize",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of songs",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Song"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "$ref": "#/definitions/error_message.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/error_message.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Add a new song to the library",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Add a new song",
                "parameters": [
                    {
                        "description": "Song data",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Song"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Song added successfully with ID",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "$ref": "#/definitions/error_message.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/error_message.ErrorResponse"
                        }
                    },
                    "502": {
                        "description": "Failed to fetch data from external API",
                        "schema": {
                            "$ref": "#/definitions/error_message.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/songs/{id}": {
            "get": {
                "description": "Get song text by song ID, verse, and limit",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Get song text by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Verse number",
                        "name": "verse",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 4,
                        "description": "Number of lines per verse",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Song text",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "$ref": "#/definitions/error_message.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Song not found",
                        "schema": {
                            "$ref": "#/definitions/error_message.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/error_message.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update song details by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Update an existing song",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated song data",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Song"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Song updated successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "$ref": "#/definitions/error_message.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/error_message.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Remove a song from the library by its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Delete a song by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Song deleted successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Invalid song ID",
                        "schema": {
                            "$ref": "#/definitions/error_message.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/error_message.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "error_message.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "model.Song": {
            "type": "object",
            "properties": {
                "details": {
                    "description": "Дополнительные детали о песне",
                    "allOf": [
                        {
                            "$ref": "#/definitions/model.SongDetail"
                        }
                    ]
                },
                "group": {
                    "description": "Исполнитель или группа",
                    "type": "string"
                },
                "id": {
                    "description": "Уникальный идентификатор песни",
                    "type": "integer"
                },
                "song": {
                    "description": "Название песни",
                    "type": "string"
                }
            }
        },
        "model.SongDetail": {
            "type": "object",
            "properties": {
                "album": {
                    "type": "string"
                },
                "artist": {
                    "type": "string"
                },
                "duration": {
                    "description": "Продолжительность",
                    "type": "string"
                },
                "genre": {
                    "description": "Жанр",
                    "type": "string"
                },
                "key": {
                    "description": "Ключ",
                    "type": "string"
                },
                "link": {
                    "type": "string"
                },
                "release_date": {
                    "type": "string"
                },
                "tempo": {
                    "description": "Темп",
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9000",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Music Library API",
	Description:      "This is a simple API for managing a music library.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
