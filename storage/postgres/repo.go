package postgres

import "github.com/jmoiron/sqlx"

type todoRepo struct {
	db *sqlx.DB
}

func NewTodoRepo(db *sqlx.DB) *todoRepo {
	return &todoRepo{db: db}
}
