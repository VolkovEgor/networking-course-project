package postgres

import (
	"database/sql"
	"fmt"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type BoardPg struct {
	db *sqlx.DB
}

func NewBoardPg(db *sqlx.DB) *BoardPg {
	return &BoardPg{db: db}
}

func (r *BoardPg) Create(userId int, board *models.Board) (int, error) {
	logrus.Info(board.Title, board.Datetimes, board.DefaultPermissions)
	var boardId int

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	defPermissionId, err := createPermissions(tx, board.DefaultPermissions)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	datetimesId, err := createDatetimes(tx, board.Datetimes)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	query := fmt.Sprintf(
		`INSERT INTO %s
		(project_id, owner_id, default_permissions_id, datetimes_id, title)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`, boardsTable)

	row := tx.QueryRow(query, board.ProjectId, board.OwnerId, defPermissionId,
		datetimesId, board.Title)
	if err := row.Scan(&boardId); err != nil {
		tx.Rollback()
		return 0, err
	}

	permission := &models.Permission{
		Read:  true,
		Write: true,
		Admin: true,
	}
	permissionId, err := createPermissions(tx, permission)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// TODO: права владельца проекта
	query = fmt.Sprintf(
		`INSERT INTO %s (user_id, board_id, permissions_id)
		VALUES ($1, $2, $3)`, boardUsersTable)

	_, err = tx.Exec(query, userId, boardId, permissionId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return boardId, nil
}

func (r *BoardPg) GetById(boardId int) (*models.Board, error) {
	board := &models.Board{}
	defaultPermissions := &models.Permission{}
	datetimes := &models.Datetimes{}

	query := fmt.Sprintf(
		`SELECT b.id, b.project_id, b.owner_id, bper.read, bper.write, bper.admin, 
		d.created, d.updated, d.accessed, b.title
		FROM %s AS b
			INNER JOIN %s AS bper ON b.default_permissions_id = bper.id
			INNER JOIN %s AS d ON b.datetimes_id = d.id
		WHERE b.id = $1`,
		boardsTable, permissionsTable, datetimesTable)

	row := r.db.QueryRow(query, boardId)
	err := row.Scan(&board.Id, &board.ProjectId, &board.OwnerId,
		&defaultPermissions.Read, &defaultPermissions.Write, &defaultPermissions.Admin,
		&datetimes.Created, &datetimes.Updated, &datetimes.Accessed,
		&board.Title)
	if err != nil {
		return nil, err
	}

	board.DefaultPermissions = defaultPermissions
	board.Datetimes = datetimes
	return board, nil
}

func (r *BoardPg) GetAll(userId, projectId int) ([]*models.Board, error) {
	var boards []*models.Board

	query := fmt.Sprintf(
		`SELECT b.id, b.project_id, b.owner_id, bper.read, bper.write, bper.admin, 
		d.created, d.updated, d.accessed, b.title
		FROM %s AS bu
			INNER JOIN %s AS per ON bu.permissions_id = per.id
			INNER JOIN %s AS b ON bu.board_id = b.id
			INNER JOIN %s AS bper ON b.default_permissions_id = bper.id
			INNER JOIN %s AS d ON b.datetimes_id = d.id
		WHERE bu.user_id = $1 AND b.project_id = $2 AND per.read = true`,
		boardUsersTable, permissionsTable, boardsTable, permissionsTable,
		datetimesTable)

	rows, err := r.db.Query(query, userId, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		board := &models.Board{}
		defaultPermissions := &models.Permission{}
		datetimes := &models.Datetimes{}

		err := rows.Scan(&board.Id, &board.ProjectId, &board.OwnerId,
			&defaultPermissions.Read, &defaultPermissions.Write, &defaultPermissions.Admin,
			&datetimes.Created, &datetimes.Updated, &datetimes.Accessed,
			&board.Title)

		if err != nil {
			return nil, err
		}

		board.DefaultPermissions = defaultPermissions
		board.Datetimes = datetimes
		boards = append(boards, board)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return boards, nil
}

func (r *BoardPg) Update(boardId int, input *models.UpdateBoard) error {
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

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s SET %s where id=$%d`,
		boardsTable, setQuery, argId)
	args = append(args, boardId)
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	defPermissionsId, datetimesId, err := r.getBoardForeignKeys(boardId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if input.DefaultPermissions != nil {
		if err = updatePermissions(tx, defPermissionsId, input.DefaultPermissions); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = updateDatetimes(tx, datetimesId, input.Datetimes); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return err
}

func (r *BoardPg) Delete(boardId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(
		`DELETE FROM %s AS per USING %s AS bu
		WHERE per.id = bu.permissions_id AND bu.board_id=$1`,
		permissionsTable, boardUsersTable)
	_, err = tx.Exec(query, boardId)
	if err != nil {
		tx.Rollback()
		return err
	}

	defPermissionsId, datetimesId, err := r.getBoardForeignKeys(boardId)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = deletePermissions(tx, defPermissionsId); err != nil {
		tx.Rollback()
		return err
	}
	if err = deleteDatetimes(tx, datetimesId); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return err
}

func (r *BoardPg) getBoardForeignKeys(boardId int) (int, int, error) {
	var defPermissionsId, datetimesId int
	query := fmt.Sprintf(
		`SELECT p.default_permissions_id, p.datetimes_id
		FROM %s AS p WHERE p.id = $1`, boardsTable)

	row := r.db.QueryRow(query, boardId)
	err := row.Scan(&defPermissionsId, &datetimesId)
	return defPermissionsId, datetimesId, err
}

func (r *BoardPg) GetPermissions(userId, boardId int) (*models.Permission, error) {
	permissions := &models.Permission{}

	query := fmt.Sprintf(
		`SELECT per.read, per.write, per.admin
		FROM %s AS bu
			INNER JOIN %s AS per ON bu.permissions_id = per.id
		WHERE bu.board_id = $1 AND bu.user_id = $2`,
		boardUsersTable, permissionsTable)

	row := r.db.QueryRow(query, boardId, userId)
	err := row.Scan(&permissions.Read, &permissions.Write, &permissions.Admin)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *BoardPg) GetBoardsCountByOwnerId(projectId, ownerId int) (int, error) {
	var count int

	query := fmt.Sprintf(
		`SELECT COUNT(*)
		FROM %s AS b
		WHERE b.project_id = $1 AND b.owner_id = $2`,
		boardsTable)
	err := r.db.QueryRow(query, projectId, ownerId).Scan(&count)
	return count, err
}

func updateOwnerIdByProjectId(tx *sql.Tx, projectId, oldOwnerId, newOwnerId int) error {
	query := fmt.Sprintf(`UPDATE %s SET owner_id=$1
		WHERE project_id = (
			SELECT project_id from %s WHERE id=$2
		) AND owner_id = $3`,
		boardsTable, boardsTable)
	fmt.Println(newOwnerId, projectId, oldOwnerId)
	_, err := tx.Exec(query, newOwnerId, projectId, oldOwnerId)
	return err
}

func updateOwnerIdByBoardId(tx *sql.Tx, boardId, oldOwnerId, newOwnerId int) error {
	query := fmt.Sprintf(`UPDATE %s SET owner_id=$1
		WHERE id = $2 AND owner_id = $3`,
		boardsTable)
	fmt.Println(newOwnerId, boardId, oldOwnerId)
	_, err := tx.Exec(query, newOwnerId, boardId, oldOwnerId)
	return err
}

func (r *BoardPg) GetMembers(boardId int) ([]*models.Member, error) {
	var members []*models.Member

	query := fmt.Sprintf(
		`SELECT u.id, u.nickname, u.avatar, per.read, per.write, per.admin,
		CASE b.owner_id
		WHEN user_id THEN true
		ELSE false
		END AS isOwner
		FROM %s AS pu
			INNER JOIN %s AS per ON pu.permissions_id = per.id
			INNER JOIN %s AS u ON pu.user_id = u.id
			INNER JOIN %s AS b ON pu.board_id = b.id
		WHERE pu.board_id = $1`,
		boardUsersTable, permissionsTable, usersTable, boardsTable)

	rows, err := r.db.Query(query, boardId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		member := &models.Member{}
		permissions := &models.Permission{}

		err := rows.Scan(&member.Id, &member.Nickname, &member.Avatar, &permissions.Read,
			&permissions.Write, &permissions.Admin, &member.IsOwner)
		if err != nil {
			return nil, err
		}

		member.Permissions = permissions
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, err
}
