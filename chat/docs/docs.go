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
        "license": {
            "name": "BSD-3-Clause",
            "url": "https://opensource.org/license/bsd-3-clause"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "Авторизоваться в системе, используя логин и пароль",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Авторизация",
                "parameters": [
                    {
                        "description": "Запрос на авторизацию",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.LoginReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешная авторизация",
                        "schema": {
                            "$ref": "#/definitions/model.LoginResp"
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Неверный данные для входа",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Ошибка валидации данных",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/favorite": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Обновляет избранный предикт.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "favorite"
                ],
                "summary": "Обновляет избранный предикт.",
                "parameters": [
                    {
                        "description": "Параметры запроса",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.FavoriteUpdateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Предикт успешно обновлен",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Добавляет новый предикт в избранное.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "favorite"
                ],
                "summary": "Добавляет новый предикт в избранное.",
                "parameters": [
                    {
                        "description": "Параметры запроса",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.FavoriteCreateReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Предикт успешно добавлен в избранное",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/favorite/list": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Возвращает список избранных предиктов.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "favorite"
                ],
                "summary": "Возвращает список избранных предиктов.",
                "responses": {
                    "200": {
                        "description": "Список избранных предиктов",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.FavoriteDto"
                            }
                        }
                    }
                }
            }
        },
        "/favorite/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Возвращает избранный предикт по ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "favorite"
                ],
                "summary": "Возвращает избранный предикт по ID.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID предикта",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Избранный предикт",
                        "schema": {
                            "$ref": "#/definitions/model.FavoriteDto"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Предикт не найден",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Удаляет избранный предикт по ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "favorite"
                ],
                "summary": "Удаляет избранный предикт по ID.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID предикта",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Предикт успешно удален",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Предикт не найден",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/session/list": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Получить список сессий в порядке убывания момента создания от последней к первой для авторизованного пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "session"
                ],
                "summary": "Список сессий",
                "responses": {
                    "200": {
                        "description": "Список сессий",
                        "schema": {
                            "$ref": "#/definitions/model.ListResp"
                        }
                    }
                }
            }
        },
        "/session/new": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Создать новую сессию в чате для авторизованного пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "session"
                ],
                "summary": "Новая сессия",
                "responses": {
                    "201": {
                        "description": "ID новой сессии",
                        "schema": {
                            "$ref": "#/definitions/model.NewSessionResp"
                        }
                    },
                    "400": {
                        "description": "Сессия с таким ID уже существует",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/session/rename": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Обновить заголовок сессии, который отображается в интерфейсе",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "session"
                ],
                "summary": "Обновить заголовок",
                "parameters": [
                    {
                        "description": "ID и новое название сессии",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.RenameReq"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Сессия переименована",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Сессия с таким ID уже существует",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Сессия с таким ID не найдена",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/session/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Получить все запросы и ответы для сессии по ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "session"
                ],
                "summary": "Получить данные о сессии",
                "parameters": [
                    {
                        "type": "string",
                        "description": "UUID сессии",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Информация о сессии",
                        "schema": {
                            "$ref": "#/definitions/model.FindOneResp"
                        }
                    },
                    "400": {
                        "description": "Неверное значение ID",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Удалить сессию по ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "session"
                ],
                "summary": "Удалить сессию",
                "parameters": [
                    {
                        "type": "string",
                        "description": "UUID сессии",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Сессия удалена",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неверное значение ID",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Сессия с таким ID не найдена",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/stock/unique_codes/{organization_id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Получить набор уникальных записей с разделением на регулярные и нерегулярные товары",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "stock"
                ],
                "summary": "Регулярные товары",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID организации",
                        "name": "organization_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список с товарами",
                        "schema": {
                            "$ref": "#/definitions/model.UniqueCodesResp"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.FavoriteCreateReq": {
            "type": "object",
            "properties": {
                "response": {
                    "type": "object",
                    "additionalProperties": {}
                }
            }
        },
        "model.FavoriteDto": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "response": {
                    "type": "object",
                    "additionalProperties": {}
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "model.FavoriteUpdateReq": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "response": {
                    "type": "object",
                    "additionalProperties": {}
                }
            }
        },
        "model.FindOneResp": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.SessionContentDto"
                    }
                },
                "editable": {
                    "type": "boolean"
                },
                "id": {
                    "type": "string"
                },
                "tg": {
                    "type": "boolean"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "model.ListResp": {
            "type": "object",
            "properties": {
                "sessions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.SessionDto"
                    }
                }
            }
        },
        "model.LoginReq": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "model.LoginResp": {
            "type": "object",
            "properties": {
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "model.NewSessionResp": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "model.QueryDto": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "period": {
                    "type": "string"
                },
                "product": {
                    "type": "string"
                },
                "prompt": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "model.RenameReq": {
            "type": "object",
            "required": [
                "id",
                "title"
            ],
            "properties": {
                "id": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "model.ResponseDto": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "data": {
                    "type": "object",
                    "additionalProperties": {}
                },
                "data_type": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "model.SessionContentDto": {
            "type": "object",
            "properties": {
                "query": {
                    "$ref": "#/definitions/model.QueryDto"
                },
                "response": {
                    "$ref": "#/definitions/model.ResponseDto"
                }
            }
        },
        "model.SessionDto": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "tg": {
                    "type": "boolean"
                },
                "tg_id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "model.UniqueCodeDto": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "regular": {
                    "type": "boolean"
                },
                "segment": {
                    "type": "string"
                }
            }
        },
        "model.UniqueCodesResp": {
            "type": "object",
            "properties": {
                "codes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.UniqueCodeDto"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "api.misis.larek.tech",
	BasePath:         "/chat",
	Schemes:          []string{},
	Title:            "Chat service API",
	Description:      "Документация API чат-сервиса команды misis.tech",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
