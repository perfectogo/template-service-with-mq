package repo

import "github.com/perfectogo/template-service-with-mq/models"

type TodoStorageInterface interface {
	InsertTodo(models.CrUpTodo) (models.RespTodo, error)
	SelectAllTodo(queryParam models.TodoQueryParamModel) (models.RespTodoList, error)
	SelectTodo(id string) (models.RespTodo, error)
	UpdateTodo(editedTodo models.CrUpTodo) (models.RespTodo, error)
	DeleteTodo(req models.ReqById) error
}
