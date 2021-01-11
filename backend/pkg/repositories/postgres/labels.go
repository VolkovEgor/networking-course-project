package postgres

import (
	"fmt"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

type LabelPg struct {
	db *sqlx.DB
}

func NewLabelPg(db *sqlx.DB) *LabelPg {
	return &LabelPg{db: db}
}

func (r *LabelPg) GetAllInTask(taskId int) ([]*models.Label, error) {
	var labels []*models.Label
	query := fmt.Sprintf(
		`SELECT l.id, l.board_id, l.name, l.color
		FROM %s AS l
			INNER JOIN %s AS tl ON l.id = tl.label_id
		WHERE tl.task_id = $1`,
		labelsTable, taskLabelsTable)

	rows, err := r.db.Query(query, taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		label := &models.Label{}
		err := rows.Scan(&label.Id, &label.BoardId, &label.Name, &label.Color)
		if err != nil {
			return nil, err
		}
		labels = append(labels, label)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return labels, nil
}

func (r *LabelPg) GetAll(boardId int) ([]*models.Label, error) {
	var labels []*models.Label
	query := fmt.Sprintf(`SELECT * FROM %s AS l WHERE l.board_id = $1`, labelsTable)

	rows, err := r.db.Query(query, boardId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		label := &models.Label{}
		err := rows.Scan(&label.Id, &label.BoardId, &label.Name, &label.Color)
		if err != nil {
			return nil, err
		}
		labels = append(labels, label)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return labels, nil
}

func (r *LabelPg) GetById(labelId int) (*models.Label, error) {
	label := &models.Label{}

	query := fmt.Sprintf(`SELECT * FROM %s AS l WHERE l.id = $1`, labelsTable)
	row := r.db.QueryRow(query, labelId)
	err := row.Scan(&label.Id, &label.BoardId, &label.Name, &label.Color)
	if err != nil {
		return nil, err
	}

	return label, nil
}

func (r *LabelPg) Create(label *models.Label) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf(
		`INSERT INTO %s (board_id, name, color)
		VALUES ($1, $2, $3) RETURNING id`, labelsTable)

	var id int
	row := tx.QueryRow(query, label.BoardId, label.Name, label.Color)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return id, nil
}

func (r *LabelPg) CreateInTask(taskId, labelId int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	query := fmt.Sprintf(
		`INSERT INTO %s (task_id, label_id)
		VALUES ($1, $2) RETURNING id`, taskLabelsTable)
	row := tx.QueryRow(query, taskId, labelId)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return id, nil
}

func (r *LabelPg) Update(labelId int, input *models.UpdateLabel) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Color != nil {
		setValues = append(setValues, fmt.Sprintf("color=$%d", argId))
		args = append(args, *input.Color)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s SET %s where id=$%d`,
		labelsTable, setQuery, argId)
	args = append(args, labelId)
	fmt.Println(query)
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return err
}

func (r *LabelPg) DeleteInTask(taskId, labelId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`DELETE FROM %s AS tl WHERE tl.label_id = $1 AND tl.task_id = $2`, taskLabelsTable)
	_, err = tx.Exec(query, taskId, labelId)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return err
}

func (r *LabelPg) Delete(labelId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`DELETE FROM %s AS l WHERE l.id = $1`, labelsTable)
	_, err = tx.Exec(query, labelId)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return err
}

// func (r *TaskPg) getTaskForeignKey(boardId int) (int, error) {
// 	var datetimesId int
// 	query := fmt.Sprintf(
// 		`SELECT p.datetimes_id
// 		FROM %s AS p WHERE p.id = $1`, tasksTable)

// 	row := r.db.QueryRow(query, boardId)
// 	err := row.Scan(&datetimesId)
// 	return datetimesId, err
// }

// func (r *TaskPg) updateTaskPosition(tx *sql.Tx, listId, start int, operation string) error {
// 	query := fmt.Sprintf(
// 		`UPDATE %s SET position = position %s 1
// 		WHERE list_id = $1 AND position >= $2`,
// 		tasksTable, operation)
// 	_, err := tx.Exec(query, listId, start)
// 	return err
// }

// func getTaskMaxPosition(tx *sql.Tx, listId int) (int, error) {
// 	var position int
// 	query := fmt.Sprintf(
// 		`SELECT MAX(t.position)
// 		FROM %s AS t
// 			INNER JOIN %s AS tl ON tl.id = t.list_id
// 		WHERE tl.id = $1;`, tasksTable, taskListsTable)

// 	row := tx.QueryRow(query, listId)
// 	err := row.Scan(&position)
// 	return position, err
// }

// func checkTaskOutOfBounds(tx *sql.Tx, newPos, newListId int, is_insert bool) error {
// 	maxPos, err := getTaskMaxPosition(tx, newListId)
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}
// 	if is_insert == true {
// 		maxPos++
// 	}
// 	if newPos > maxPos {
// 		tx.Rollback()
// 		return errors.New("Task position out of bounds")
// 	}
// 	return err
// }
