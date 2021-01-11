package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strings"
)

func createDatetimes(tx *sql.Tx, permissions *models.Datetimes) (int, error) {
	var datetimesId int
	query := fmt.Sprintf(
		`INSERT INTO %s (created, updated, accessed)
		VALUES ($1, $2, $3) RETURNING id`, datetimesTable)

	row := tx.QueryRow(query, permissions.Created, permissions.Updated,
		permissions.Accessed)
	err := row.Scan(&datetimesId)
	return datetimesId, err
}

func updateDatetimes(tx *sql.Tx, datetimesId int, input *models.UpdateDatetimes) error {
	if input == nil {
		return errors.New("Datetimes is not defined")
	}

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Created != nil {
		setValues = append(setValues, fmt.Sprintf("created=$%d", argId))
		args = append(args, *input.Created)
		argId++
	}

	if input.Updated != nil {
		setValues = append(setValues, fmt.Sprintf("updated=$%d", argId))
		args = append(args, *input.Updated)
		argId++
	}

	if input.Accessed != nil {
		setValues = append(setValues, fmt.Sprintf("accessed=$%d", argId))
		args = append(args, *input.Accessed)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s SET %s where id=$%d`,
		datetimesTable, setQuery, argId)
	args = append(args, datetimesId)
	_, err := tx.Exec(query, args...)

	return err
}

func deleteDatetimes(tx *sql.Tx, datetimesId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, datetimesTable)
	_, err := tx.Exec(query, datetimesId)
	return err
}
