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
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/companies": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Creates a company",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "company"
                ],
                "summary": "Create a company",
                "parameters": [
                    {
                        "description": "company creation parameter",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.CreateCompanyParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Company"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/company/{companyId}/projects/": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "create a  project",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "projects"
                ],
                "summary": "Create a project",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Company ID",
                        "name": "companyId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "project creation parameters",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.CreateProjectParam"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Project"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/company/{companyId}/projects/{projectId}/tasks/": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Create a tasks",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Create a Task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Project ID",
                        "name": "projectId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Company ID",
                        "name": "companyId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "task creation parameter",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.CreateTaskParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Task"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/company/{companyId}/projects/{projectId}/tasks/{taskId}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "get a task by its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Get a Task by its ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "taskId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Project ID",
                        "name": "projectId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Company ID",
                        "name": "companyId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Task"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Updates a tasks",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Update a Task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "taskId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Project ID",
                        "name": "projectId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Company ID",
                        "name": "companyId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "task update parameter",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.UpdateTaskParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Task"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "delete a task by its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Delete a Task by its ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "taskId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Project ID",
                        "name": "projectId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Company ID",
                        "name": "companyId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/payment/{paymentId}/projects/{projectId}/pay/{taskId}": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Adds a payment source and a payment amout for a task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Pay for a task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "taskId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Project ID",
                        "name": "projectId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Company ID",
                        "name": "companyId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "payment source",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.CreateTaskPaymentParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.TaskPayment"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/users/": {
            "post": {
                "description": "Create a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Create a User",
                "parameters": [
                    {
                        "description": "user creation parameters",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.createUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/users/login": {
            "post": {
                "description": "Logs an user in",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Logs an user in",
                "parameters": [
                    {
                        "description": "user creation parameters",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.loginUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.loginUserResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        }
    },
    "definitions": {
        "api.createUserRequest": {
            "type": "object",
            "required": [
                "company_id",
                "email",
                "full_name",
                "password",
                "user_name",
                "user_role"
            ],
            "properties": {
                "company_id": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "full_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 6
                },
                "user_name": {
                    "type": "string"
                },
                "user_role": {
                    "type": "string"
                }
            }
        },
        "api.loginUserRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "minLength": 6
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "api.loginUserResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "access_token_expires_at": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                },
                "refresh_token_expires_at": {
                    "type": "string"
                },
                "session_id": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/api.userResponse"
                }
            }
        },
        "api.userResponse": {
            "type": "object",
            "properties": {
                "company_id": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "full_name": {
                    "type": "string"
                },
                "password_changed_at": {
                    "type": "string"
                },
                "user_name": {
                    "type": "string"
                },
                "user_role": {
                    "type": "string"
                }
            }
        },
        "db.CreateCompanyParams": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "company_name": {
                    "type": "string"
                }
            }
        },
        "db.CreateProjectParam": {
            "type": "object",
            "required": [
                "client",
                "project_name",
                "responsible"
            ],
            "properties": {
                "address": {
                    "type": "string"
                },
                "client": {
                    "type": "string"
                },
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                },
                "project_name": {
                    "type": "string"
                },
                "responsible": {
                    "type": "string"
                }
            }
        },
        "db.CreateTaskParams": {
            "type": "object",
            "required": [
                "budget",
                "project_id",
                "task_name",
                "task_order"
            ],
            "properties": {
                "budget": {
                    "type": "number"
                },
                "project_id": {
                    "type": "integer"
                },
                "task_name": {
                    "type": "string"
                },
                "task_order": {
                    "type": "integer"
                }
            }
        },
        "db.CreateTaskPaymentParams": {
            "type": "object",
            "required": [
                "amount",
                "source",
                "taskIf"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                },
                "createdBy": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                },
                "taskIf": {
                    "type": "integer"
                }
            }
        },
        "db.UpdateTaskParams": {
            "type": "object",
            "required": [
                "budget",
                "done",
                "project_id",
                "rating",
                "task_name",
                "version"
            ],
            "properties": {
                "budget": {
                    "type": "number"
                },
                "done": {
                    "type": "boolean"
                },
                "id": {
                    "type": "integer"
                },
                "project_id": {
                    "type": "integer"
                },
                "rating": {
                    "type": "integer",
                    "maximum": 5
                },
                "task_name": {
                    "type": "string"
                },
                "task_order": {
                    "type": "integer"
                },
                "version": {
                    "type": "integer"
                }
            }
        },
        "domain.Company": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "company_name": {
                    "type": "string"
                },
                "createdBy": {
                    "type": "string"
                },
                "createdOn": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "updatedBy": {
                    "type": "string"
                },
                "updatedOn": {
                    "type": "string"
                }
            }
        },
        "domain.Project": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "budget": {
                    "type": "number"
                },
                "client": {
                    "type": "string"
                },
                "company_id": {
                    "type": "integer"
                },
                "completion_percentage": {
                    "type": "number"
                },
                "created_by": {
                    "type": "string"
                },
                "created_on": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                },
                "project_name": {
                    "type": "string"
                },
                "responsible": {
                    "type": "string"
                },
                "tasks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.Task"
                    }
                }
            }
        },
        "domain.Task": {
            "type": "object",
            "properties": {
                "budget": {
                    "type": "number"
                },
                "company_id": {
                    "type": "integer"
                },
                "createdBy": {
                    "type": "string"
                },
                "createdOn": {
                    "type": "string"
                },
                "done": {
                    "type": "boolean"
                },
                "id": {
                    "type": "integer"
                },
                "paid_amount": {
                    "type": "number"
                },
                "project_id": {
                    "type": "integer"
                },
                "rating": {
                    "type": "integer"
                },
                "task_name": {
                    "type": "string"
                },
                "task_order": {
                    "type": "integer"
                },
                "updatedBy": {
                    "type": "string"
                },
                "updatedOn": {
                    "type": "string"
                },
                "version": {
                    "type": "integer"
                }
            }
        },
        "domain.TaskPayment": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "createdBy": {
                    "type": "string"
                },
                "createdOn": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "source": {
                    "type": "string"
                },
                "task_id": {
                    "type": "integer"
                },
                "updatedBy": {
                    "type": "string"
                },
                "updatedOn": {
                    "type": "string"
                }
            }
        },
        "domain.User": {
            "type": "object",
            "properties": {
                "company_id": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "full_name": {
                    "type": "string"
                },
                "hashed_password": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "password_changed_at": {
                    "type": "string"
                },
                "user_name": {
                    "type": "string"
                },
                "user_role": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Swagger ChainTasks API",
	Description:      "This a server helping organizing tasks in a chain.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
