package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strings"
)

func createPermissions(tx *sql.Tx, permissions *models.Permission) (int, error) {
	var permissionId int
	query := fmt.Sprintf(
		`INSERT INTO %s (read, write, admin)
		VALUES ($1, $2, $3) RETURNING id`, permissionsTable)

	row := tx.QueryRow(query, permissions.Read, permissions.Write,
		permissions.Admin)
	err := row.Scan(&permissionId)
	return permissionId, err
}

func updatePermissions(tx *sql.Tx, permissionsId int, input *models.UpdatePermission) error {
	if input == nil {
		return errors.New("Permissions is not defined")
	}

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Read != nil {
		setValues = append(setValues, fmt.Sprintf("read=$%d", argId))
		args = append(args, *input.Read)
		argId++
	}

	if input.Write != nil {
		setValues = append(setValues, fmt.Sprintf("write=$%d", argId))
		args = append(args, *input.Write)
		argId++
	}

	if input.Admin != nil {
		setValues = append(setValues, fmt.Sprintf("admin=$%d", argId))
		args = append(args, *input.Admin)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s SET %s where id=$%d`,
		permissionsTable, setQuery, argId)
	args = append(args, permissionsId)
	_, err := tx.Exec(query, args...)
	return err
}

func deletePermissions(tx *sql.Tx, permissionsId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, permissionsTable)
	_, err := tx.Exec(query, permissionsId)
	return err
}
