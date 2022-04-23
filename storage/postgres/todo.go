package postgres

import (
	"database/sql"

	"github.com/perfectogo/template-service-with-mq/models"
)

func (r *todoRepo) InsertTodo(newTodo models.CrUpTodo) (todo models.RespTodo, err error) {
	var id string
	if err = r.db.QueryRow(
		`INSERT INTO 
			todos 
				todo_id, 
				title, 
				note,
				priority
		RETURNING todo_id`,
		newTodo.TodoId,
		newTodo.Title,
		newTodo.Notes,
		newTodo.Priority,
	).Scan(&id); err != nil {
		return models.RespTodo{}, err
	}

	todo, err = r.SelectTodo(id)
	if err != nil {
		return models.RespTodo{}, err
	}

	return todo, nil
}

func (r *todoRepo) SelectTodo(id string) (todo models.RespTodo, err error) {
	var updatedAt sql.NullTime
	if err = r.db.QueryRow(
		`SELECT
			todo_id,
			title,
			note,
			priority,
			created_at,
			updated_at
		FROM
			todos
		WHERE
			todo_id=$1 AND
			deleted_at IS NULL`,
		id,
	).Scan(
		&todo.TodoId,
		&todo.Title,
		&todo.Notes,
		&todo.Priority,
		&todo.CreatedAt,
		&updatedAt,
	); err != nil {
		return models.RespTodo{}, err
	}
	todo.UpdatedAt = updatedAt.Time
	return todo, nil
}
