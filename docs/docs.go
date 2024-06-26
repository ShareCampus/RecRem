// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/all_users": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "获取所有用户",
                "responses": {
                    "100": {
                        "description": "Continue",
                        "schema": {
                            "$ref": "#/definitions/utils.Result"
                        }
                    },
                    "104": {
                        "schema": {
                            "$ref": "#/definitions/utils.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/users": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "更新用户信息",
                "parameters": [
                    {
                        "description": "用户信息表单",
                        "name": "userForm",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/forms.UserInfoForm"
                        }
                    }
                ],
                "responses": {
                    "100": {
                        "description": "Continue",
                        "schema": {
                            "$ref": "#/definitions/utils.Result"
                        }
                    },
                    "104": {
                        "schema": {
                            "$ref": "#/definitions/utils.Result"
                        }
                    }
                }
            }
        },
        "/api/v1/users/pwd": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "修改密码",
                "parameters": [
                    {
                        "description": "修改密码表单",
                        "name": "pwdForm",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/forms.PwdForm"
                        }
                    }
                ],
                "responses": {
                    "100": {
                        "description": "Continue",
                        "schema": {
                            "$ref": "#/definitions/utils.Result"
                        }
                    },
                    "104": {
                        "schema": {
                            "$ref": "#/definitions/utils.Result"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "forms.PwdForm": {
            "type": "object",
            "required": [
                "confirm_new_pwd",
                "new_pwd",
                "old_pwd",
                "username"
            ],
            "properties": {
                "confirm_new_pwd": {
                    "type": "string"
                },
                "new_pwd": {
                    "type": "string"
                },
                "old_pwd": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "forms.UserInfoForm": {
            "type": "object",
            "required": [
                "email",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "nickname": {
                    "type": "string"
                },
                "signature": {
                    "type": "string"
                },
                "user_img": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "utils.Result": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "状态码",
                    "type": "integer",
                    "example": 100
                },
                "data": {
                    "description": "数据",
                    "type": "object"
                },
                "msg": {
                    "description": "提示",
                    "type": "string",
                    "example": "密码错误"
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
	Host:        "localhost:8088",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Gin Swagger",
	Description: "recrem 开源搜索引擎 API 接口文档",
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
