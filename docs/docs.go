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
        "/music-library/list": {
            "get": {
                "description": "page and pageSize are required, every other field is a filter, if it's empty, it is treated as absence of filter",
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "music-library"
                ],
                "summary": "Fetches song data in pages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "group name",
                        "name": "group",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "song name",
                        "name": "song",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "dates before this will not show up",
                        "name": "releaseDateLower",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "dates after this will not show up",
                        "name": "releaseDateUpper",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "lyrics",
                        "name": "text",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "page number",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "number of songs displayed per page",
                        "name": "pageSize",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.FilteredListResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/music-library/lyrics": {
            "get": {
                "description": "If page/pageSize combination results in an empty page, you will get status code 400 :)",
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "music-library"
                ],
                "summary": "Fetches lyrics divided into verses",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "song id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "page number",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "number of verses per page",
                        "name": "pageSize",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.VersesResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/music-library/song": {
            "put": {
                "description": "You need to provide id and 3 other fields, on success returns provided id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "music-library"
                ],
                "summary": "Updates song info",
                "parameters": [
                    {
                        "description": "id and additional info",
                        "name": "UpdateRequestJSON",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.UpdateRequestJSON"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.IdResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Makes a request into externalAPIURL, if it fails, returns status 201,\nsaves basic song info and writes encountered errors into response body (field \"songDetails\")",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "music-library"
                ],
                "summary": "Adds song into library",
                "parameters": [
                    {
                        "description": "group and song names",
                        "name": "BasicSongInfoJSON",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.BasicSongInfoJSON"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.IdResponse"
                        }
                    },
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/handlers.AddSongFailedExternalAPIResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Returns provided id if deletion succeeded",
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "music-library"
                ],
                "summary": "Deletes song from library",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "song id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.IdResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.AddSongFailedExternalAPIResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "songDetails": {
                    "type": "string"
                }
            }
        },
        "handlers.BasicSongInfoJSON": {
            "type": "object",
            "properties": {
                "group": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                }
            }
        },
        "handlers.FilteredListResponse": {
            "type": "object",
            "properties": {
                "filteredRows": {}
            }
        },
        "handlers.UpdateRequestJSON": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "link": {
                    "type": "string"
                },
                "releaseDate": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "models.ErrorResponse": {
            "type": "object",
            "properties": {
                "errors": {}
            }
        },
        "models.IdResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "models.VersesResponse": {
            "type": "object",
            "properties": {
                "verses": {}
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
