package repo

import "github.com/perfectogo/template-service-with-mq/models"

type TodoStorageInterface interface {
	InsertTodo(models.CrUpTodo) (models.RespTodo, error)
	SelectTodo(id string) (models.RespTodo, error)
}
