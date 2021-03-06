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
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "ovr.shevtsov@gmail.com"
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
        "/api/{channel}": {
            "get": {
                "description": "get messages for channel",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "get messages for channel",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name",
                        "name": "channel",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "pagination offset",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "pagination limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.HooksByChannel"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "post message",
                "consumes": [
                    "application/json"
                ],
                "summary": "post message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name",
                        "name": "channel",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "message",
                        "name": "message",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    {
                        "type": "integer",
                        "description": "fail requests until timestamp",
                        "name": "failUntil",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "do not handle message by service and just reply with status code",
                        "name": "justreply",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            },
            "delete": {
                "description": "delete channel",
                "consumes": [
                    "application/json"
                ],
                "summary": "delete channel",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name",
                        "name": "channel",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "entities.Hook": {
            "type": "object",
            "properties": {
                "channel": {
                    "type": "string"
                },
                "created_at": {
                    "type": "integer"
                },
                "headers": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "payload": {
                    "type": "string"
                },
                "statusOk": {
                    "type": "boolean"
                }
            }
        },
        "entities.HooksByChannel": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Hook"
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
	Host:        "webhook.monster",
	BasePath:    "/",
	Schemes:     []string{"https"},
	Title:       "webhook listener API",
	Description: "webhook listener api",
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
	swag.Register(swag.Name, &s{})
}
