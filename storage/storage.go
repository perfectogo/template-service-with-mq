package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/perfectogo/template-service-with-mq/storage/postgres"
	"github.com/perfectogo/template-service-with-mq/storage/repo"
)

type StorageInterface interface {
	Todo() repo.TodoStorageInterface
}

type storage struct {
	db       *sqlx.DB
	todoRepo repo.TodoStorageInterface
}

func NewStorage(db *sqlx.DB) *storage {
	return &storage{
		db:       db,
		todoRepo: postgres.NewTodoRepo(db),
	}
}

func (s *storage) Todo() repo.TodoStorageInterface {
	return s.todoRepo
}
