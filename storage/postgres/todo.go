package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/perfectogo/template-service-with-mq/models"
)

func (r *todoRepo) InsertTodo(newTodo models.CrUpTodo) (todo models.RespTodo, err error) {
	var id string
	if err = r.db.QueryRow(
		`INSERT INTO 
			todos 
				(todo_id, 
				title, 
				note,
				priority)
		VALUES
				($1, $2, $3, $4)
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

//SELECT
func (r *todoRepo) SelectAllTodo(queryParam models.TodoQueryParamModel) (models.RespTodoList, error) {
	//
	params := make(map[string]interface{})
	selct_query :=
		`SELECT
			todo_id,
			title,
			note,
			priority,
			created_at,
			updated_at
		FROM
			todos
		`

	//filters
	filter := " WHERE deleted_at IS NULL"
	order := " ORDER BY created_at"
	arrangement := " DESC"
	offset := " OFFSET 0"
	limit := " LIMIT 10"

	if len(queryParam.Search) > 0 {
		params["search"] = queryParam.Search
		filter += " AND (body ILIKE '%' || :search || '%')"
	}

	if len(queryParam.Order) > 0 {
		valid := regexp.MustCompile("^[A-Za-z0-9_]+$")
		if valid.MatchString(queryParam.Order) {
			order = fmt.Sprintf(" ORDER BY %s", queryParam.Order)
		} else {
			return models.RespTodoList{}, errors.New("wrong order query param")
		}
	}

	switch strings.ToUpper(queryParam.Arrangement) {
	case "DESC":
		arrangement = " DESC"
	case "ASC":
		arrangement = " ASC"
	}

	if queryParam.Offset > 0 {
		params["offset"] = queryParam.Offset
		offset = " OFFSET :offset"
	}

	if queryParam.Limit > 0 {
		params["limit"] = queryParam.Limit
		limit = " LIMIT :limit"
	}

	//Maked query
	selct_query += filter + order + arrangement + offset + limit
	//
	rows, err := r.db.NamedQuery(selct_query, params)
	if err != nil {
		return models.RespTodoList{}, err
	}
	defer rows.Close()
	var (
		todos []models.RespTodo
		count int64
	)

	for rows.Next() {

		var (
			updatedAt sql.NullTime
			todo      models.RespTodo
		)

		if err := rows.Scan(
			&todo.TodoId,
			&todo.Title,
			&todo.Notes,
			&todo.Priority,
			&todo.CreatedAt,
			&updatedAt,
		); err != nil {
			return models.RespTodoList{}, nil
		}

		todo.UpdatedAt = updatedAt.Time

		todos = append(todos, todo)
	}
	if err := r.db.QueryRow(
		`SELECT COUNT(todo_id) FROM
			todos` + filter).Scan(&count); err != nil {
		return models.RespTodoList{}, err
	}

	return models.RespTodoList{
		Todos: &todos,
		Count: count,
	}, nil
}

//
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

//
func (r *todoRepo) UpdateTodo(editedTodo models.CrUpTodo) (models.RespTodo, error) {
	var id string
	if err := r.db.QueryRow(
		`UPDATE
			todos
		SET
			title=$2,
			note=$3,
			priority=$4,
			updated_at=NOW()
		WHERE
			todo_id=$1
		RETURNING
			todo_id
			`,
		editedTodo.TodoId,
		editedTodo.Title,
		editedTodo.Notes,
		editedTodo.Priority,
	).Scan(&id); err != nil {
		return models.RespTodo{}, err
	}
	todo, err := r.SelectTodo(id)
	if err != nil {
		return models.RespTodo{}, err
	}
	return todo, nil
}

//
func (r *todoRepo) DeleteTodo(req models.ReqById) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := `UPDATE 
				todos
			SET
				deleted_at=NOW() 
			WHERE 
				todo_id = $1`

	result, err := tx.Exec(query, req.TodoId)
	if err != nil {
		_err := tx.Rollback()
		if _err != nil {
			return _err
		}
		return err
	}

	if i, err := result.RowsAffected(); i == 0 {
		if err != nil {
			_err := tx.Rollback()
			if _err != nil {
				return _err
			}
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		_err := tx.Rollback()
		if _err != nil {
			return _err
		}
		return err
	}

	return nil
}
