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
        "/organization": {
            "get": {
                "description": "Получить организацию для пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organization"
                ],
                "summary": "Получить организацию",
                "parameters": [
                    {
                        "type": "string",
                        "description": "access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Информация об организации",
                        "schema": {
                            "$ref": "#/definitions/model.OrganizationDto"
                        }
                    },
                    "404": {
                        "description": "Организация не найдена",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "description": "Обновить организацию",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organization"
                ],
                "summary": "Обновить организацию",
                "parameters": [
                    {
                        "type": "string",
                        "description": "access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Параметры обновления организации",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.OrganizationUpdateReq"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Организация обновлена"
                    },
                    "400": {
                        "description": "Неверные параметры запроса",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Создает новую организацию.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organization"
                ],
                "summary": "Создает новую организацию.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Параметры создания организации",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.OrganizationCreateReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "ID созданной организации",
                        "schema": {
                            "$ref": "#/definitions/model.OrganizationCreateResp"
                        }
                    },
                    "400": {
                        "description": "Неверные параметры запроса",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user": {
            "post": {
                "description": "Создает нового пользователя.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Создает нового пользователя.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Параметры пользователя",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.UserCreateReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Пользователь создан",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Ошибка в запросе",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Неверный формат данных",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/organization": {
            "post": {
                "description": "Добавляет пользователя в организацию.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Добавляет пользователя в организацию.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Параметры пользователя",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.UserUpdateOrganizationReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пользователь добавлен в организацию",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Ошибка в запросе",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Неверный формат данных",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/{organization}": {
            "get": {
                "description": "Возвращает список пользователей по организации.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Возвращает список пользователей по организации.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Название организации",
                        "name": "organization",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список пользователей",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/user/{username}": {
            "delete": {
                "description": "Удаляет организацию.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Удаляет организацию.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Имя пользователя",
                        "name": "username",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Организация удалена",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Ошибка в запросе",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
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
        "model.OrganizationCreateReq": {
            "type": "object",
            "properties": {
                "title": {
                    "type": "string"
                }
            }
        },
        "model.OrganizationCreateResp": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "model.OrganizationDto": {
            "type": "object",
            "properties": {
                "created_at": {
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
        "model.OrganizationUpdateReq": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "model.UserCreateReq": {
            "type": "object",
            "required": [
                "email",
                "first_name",
                "last_name",
                "organization",
                "password",
                "roles",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "organization": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "model.UserUpdateOrganizationReq": {
            "type": "object",
            "required": [
                "organization",
                "username"
            ],
            "properties": {
                "organization": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "admin.misis.larek.tech",
	BasePath:         "/chat",
	Schemes:          []string{},
	Title:            "Admin service API",
	Description:      "Документация API админ-сервиса команды misis.tech",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
