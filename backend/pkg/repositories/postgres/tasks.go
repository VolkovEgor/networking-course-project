package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

type TaskPg struct {
	db *sqlx.DB
}

func NewTaskPg(db *sqlx.DB) *TaskPg {
	return &TaskPg{db: db}
}

func (r *TaskPg) GetAll(listId int) ([]*models.Task, error) {
	var tasks []*models.Task
	query := fmt.Sprintf(
		`SELECT t.id, t.list_id, t.title, t.description, d.created, d.updated, d.accessed, t.position
		FROM %s AS t
			INNER JOIN %s AS tl ON t.list_id = tl.id
			INNER JOIN %s AS d ON t.datetimes_id = d.id
		WHERE tl.id = $1
		ORDER BY t.position`,

		tasksTable, taskListsTable, datetimesTable)

	rows, err := r.db.Query(query, listId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := &models.Task{}
		datetimes := &models.Datetimes{}

		err := rows.Scan(&task.Id, &task.ListId, &task.Title, &task.Description, &datetimes.Created,
			&datetimes.Updated, &datetimes.Accessed, &task.Position)

		if err != nil {
			return nil, err
		}

		task.Datetimes = datetimes
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskPg) GetById(taskId int) (*models.Task, error) {
	task := &models.Task{}
	datetimes := &models.Datetimes{}

	query := fmt.Sprintf(
		`SELECT t.id, t.list_id, t.title, t.description, d.created, d.updated, d.accessed, t.position
		FROM %s AS t
		INNER JOIN %s AS d ON t.datetimes_id = d.id
		WHERE t.id = $1`,
		tasksTable, datetimesTable)

	row := r.db.QueryRow(query, taskId)
	err := row.Scan(&task.Id, &task.ListId, &task.Title, &task.Description, &datetimes.Created,
		&datetimes.Updated, &datetimes.Accessed, &task.Position)
	if err != nil {
		return nil, err
	}

	task.Datetimes = datetimes
	return task, nil
}

func (r *TaskPg) Create(task *models.Task) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	datetimesId, err := createDatetimes(tx, task.Datetimes)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var id int
	position, err := getTaskMaxPosition(tx, task.ListId)
	if err != nil {
		// TODO: to use int pointer?
		position = -1
		// return 0, err
	}
	position++

	query := fmt.Sprintf(
		`INSERT INTO %s (list_id, title, description, datetimes_id, position)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`, tasksTable)

	row := tx.QueryRow(query, task.ListId, task.Title, task.Description, datetimesId, position)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	return id, nil
}

func (r *TaskPg) Update(taskId int, input *models.UpdateTask) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Position != nil {
		newPos := *input.Position

		var listId, oldPos int
		query := fmt.Sprintf(`SELECT list_id, position FROM %s WHERE id = $1`, tasksTable)
		row := tx.QueryRow(query, taskId)
		err := row.Scan(&listId, &oldPos)
		if err != nil {
			tx.Rollback()
			return err
		}

		if input.ListId != nil && *input.ListId != listId {
			newListId := *input.ListId

			// TODO: check func
			err = checkListIsExists(tx, newListId)
			if err != nil {
				tx.Rollback()
				return err
			}

			err = checkTaskOutOfBounds(tx, newPos, newListId, true)
			if err != nil {
				tx.Rollback()
				return err
			}

			err = r.updateTaskPosition(tx, listId, oldPos+1, "-")
			if err != nil {
				tx.Rollback()
				return err
			}
			err = r.updateTaskPosition(tx, newListId, newPos, "+")
			if err != nil {
				tx.Rollback()
				return err
			}
			setValues = append(setValues, fmt.Sprintf("position=$%d", argId))
			args = append(args, newPos)
			argId++

			setValues = append(setValues, fmt.Sprintf("list_id=$%d", argId))
			args = append(args, newListId)
			argId++

		} else {
			var operation string
			var start, end int

			err = checkTaskOutOfBounds(tx, newPos, listId, false)
			if err != nil {
				tx.Rollback()
				return err
			}

			if oldPos < newPos {
				operation = "-"
				start, end = oldPos+1, newPos
			} else if oldPos > newPos {
				operation = "+"
				start, end = newPos, oldPos-1
			}
			fmt.Println(oldPos, newPos)

			if operation != "" {
				setValues = append(setValues, fmt.Sprintf("position=$%d", argId))
				args = append(args, newPos)
				argId++

				query := fmt.Sprintf(
					`UPDATE %s SET position = position %s 1
					WHERE list_id = $1 AND position >= $2 AND position <= $3`,
					tasksTable, operation)
				_, err := tx.Exec(query, listId, start, end)
				if err != nil {
					tx.Rollback()
					return err
				}
			} else {
				tx.Commit()
				return nil
			}
		}
	}
	// TODO обновление Datetimes на всех уровнях
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s SET %s where id=$%d`,
		tasksTable, setQuery, argId)
	fmt.Println(query)
	args = append(args, taskId)
	fmt.Println(query)
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return err
}

func (r *TaskPg) Delete(taskId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var listId, position int
	query := fmt.Sprintf(`SELECT list_id, position FROM %s WHERE id = $1`, tasksTable)
	row := tx.QueryRow(query, taskId)
	err = row.Scan(&listId, &position)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.updateTaskPosition(tx, listId, position, "-")
	if err != nil {
		tx.Rollback()
		return err
	}

	datetimesId, err := r.getTaskForeignKey(taskId)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, datetimesTable)
	_, err = tx.Exec(query, datetimesId)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return err
}

func (r *TaskPg) getTaskForeignKey(boardId int) (int, error) {
	var datetimesId int
	query := fmt.Sprintf(
		`SELECT p.datetimes_id
		FROM %s AS p WHERE p.id = $1`, tasksTable)

	row := r.db.QueryRow(query, boardId)
	err := row.Scan(&datetimesId)
	return datetimesId, err
}

func (r *TaskPg) updateTaskPosition(tx *sql.Tx, listId, start int, operation string) error {
	query := fmt.Sprintf(
		`UPDATE %s SET position = position %s 1
		WHERE list_id = $1 AND position >= $2`,
		tasksTable, operation)
	_, err := tx.Exec(query, listId, start)
	return err
}

func getTaskMaxPosition(tx *sql.Tx, listId int) (int, error) {
	var position int
	query := fmt.Sprintf(
		`SELECT MAX(t.position)
		FROM %s AS t
			INNER JOIN %s AS tl ON tl.id = t.list_id
		WHERE tl.id = $1;`, tasksTable, taskListsTable)

	row := tx.QueryRow(query, listId)
	err := row.Scan(&position)
	return position, err
}

func checkTaskOutOfBounds(tx *sql.Tx, newPos, newListId int, is_insert bool) error {
	maxPos, err := getTaskMaxPosition(tx, newListId)
	if err != nil {
		// tx.Rollback()
		// return err
		maxPos = -1
	}
	if is_insert == true {
		maxPos++
	}
	if newPos > maxPos {
		tx.Rollback()
		return errors.New("Task position out of bounds")
	}
	return nil
}

func checkListIsExists(tx *sql.Tx, listId int) error {
	var id int
	query := fmt.Sprintf(
		`SELECT id FROM %s WHERE id = $1`, taskListsTable)
	row := tx.QueryRow(query, listId)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		if err.Error() == DbResultNotFound {
			return errors.New("List is not exists")
		} else {
			return err
		}
	}
	return nil
}
