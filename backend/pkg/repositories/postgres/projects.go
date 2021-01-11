package postgres

import (
	"fmt"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type ProjectPg struct {
	db *sqlx.DB
}

func NewProjectPg(db *sqlx.DB) *ProjectPg {
	return &ProjectPg{db: db}
}

func (r *ProjectPg) Create(project *models.Project) (int, error) {
	var projectId int

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	defPermissionId, err := createPermissions(tx, project.DefaultPermissions)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	datetimesId, err := createDatetimes(tx, project.Datetimes)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	query := fmt.Sprintf(
		`INSERT INTO %s
		(owner_id, default_permissions_id, datetimes_id, title, description)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`, projectsTable)

	row := tx.QueryRow(query, project.OwnerId, defPermissionId,
		datetimesId, project.Title, project.Description)
	if err := row.Scan(&projectId); err != nil {
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

	query = fmt.Sprintf(
		`INSERT INTO %s (user_id, project_id, permissions_id)
		VALUES ($1, $2, $3)`, projectUsersTable)

	_, err = tx.Exec(query, project.OwnerId, projectId, permissionId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return projectId, nil
}

func (r *ProjectPg) GetById(projectId int) (*models.Project, error) {
	project := &models.Project{}
	defaultPermissions := &models.Permission{}
	datetimes := &models.Datetimes{}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	_, datetimesId, err := r.getProjectForeignKeys(projectId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	curTime := time.Now().Unix()
	upDatetimes := &models.UpdateDatetimes{
		Accessed: &curTime,
	}
	if err = updateDatetimes(tx, datetimesId, upDatetimes); err != nil {
		tx.Rollback()
		return nil, err
	}

	query := fmt.Sprintf(
		`SELECT p.id, p.owner_id, dper.read, dper.write, dper.admin, 
		d.created, d.updated, d.accessed, p.title, p.description
		FROM %s AS p
			INNER JOIN %s AS dper ON p.default_permissions_id = dper.id
			INNER JOIN %s AS d ON p.datetimes_id = d.id
		WHERE p.id = $1`,
		projectsTable, permissionsTable, datetimesTable)

	row := tx.QueryRow(query, projectId)
	err = row.Scan(&project.Id, &project.OwnerId, &defaultPermissions.Read,
		&defaultPermissions.Write, &defaultPermissions.Admin,
		&datetimes.Created, &datetimes.Updated, &datetimes.Accessed,
		&project.Title, &project.Description)

	if err != nil {
		tx.Rollback()
		return nil, err
	}
	project.DefaultPermissions = defaultPermissions
	project.Datetimes = datetimes
	tx.Commit()
	return project, err
}

func (r *ProjectPg) GetAll(userId int) ([]*models.Project, error) {
	var projects []*models.Project

	query := fmt.Sprintf(
		`SELECT p.id, p.owner_id, dper.read, dper.write, dper.admin, 
		d.created, d.updated, d.accessed, p.title, p.description
		FROM %s AS pu
			INNER JOIN %s AS per ON pu.permissions_id = per.id
			INNER JOIN %s AS p ON pu.project_id = p.id
			INNER JOIN %s AS dper ON p.default_permissions_id = dper.id
			INNER JOIN %s AS d ON p.datetimes_id = d.id
		WHERE pu.user_id = $1 AND per.read = true`,
		projectUsersTable, permissionsTable, projectsTable, permissionsTable,
		datetimesTable)

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		project := &models.Project{}
		defaultPermissions := &models.Permission{}
		datetimes := &models.Datetimes{}

		err := rows.Scan(&project.Id, &project.OwnerId, &defaultPermissions.Read,
			&defaultPermissions.Write, &defaultPermissions.Admin,
			&datetimes.Created, &datetimes.Updated, &datetimes.Accessed,
			&project.Title, &project.Description)

		if err != nil {
			return nil, err
		}

		project.DefaultPermissions = defaultPermissions
		project.Datetimes = datetimes
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, err
}

func (r *ProjectPg) Update(projectId int, input *models.UpdateProject) error {
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

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s SET %s where id=$%d`,
		projectsTable, setQuery, argId)
	args = append(args, projectId)
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	defPermissionsId, datetimesId, err := r.getProjectForeignKeys(projectId)
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

func (r *ProjectPg) Delete(projectId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(
		`DELETE FROM %s AS per USING %s AS pu
		WHERE per.id = pu.permissions_id AND pu.project_id=$1`,
		permissionsTable, projectUsersTable)
	_, err = tx.Exec(query, projectId)
	if err != nil {
		tx.Rollback()
		return err
	}

	defPermissionsId, datetimesId, err := r.getProjectForeignKeys(projectId)
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

func (r *ProjectPg) GetPermissions(userId, projectId int) (*models.Permission, error) {
	permissions := &models.Permission{}

	query := fmt.Sprintf(
		`SELECT per.read, per.write, per.admin
		FROM %s AS pu
			INNER JOIN %s AS per ON pu.permissions_id = per.id
		WHERE pu.project_id = $1 AND pu.user_id = $2`,
		projectUsersTable, permissionsTable)

	row := r.db.QueryRow(query, projectId, userId)
	err := row.Scan(&permissions.Read, &permissions.Write, &permissions.Admin)
	return permissions, err
}

func (r *ProjectPg) getProjectForeignKeys(projectId int) (int, int, error) {
	var defPermissionsId, datetimesId int
	query := fmt.Sprintf(
		`SELECT p.default_permissions_id, p.datetimes_id
		FROM %s AS p WHERE p.id = $1`, projectsTable)

	row := r.db.QueryRow(query, projectId)
	err := row.Scan(&defPermissionsId, &datetimesId)
	return defPermissionsId, datetimesId, err
}

func (r *ProjectPg) GetMembers(projectId int) ([]*models.Member, error) {
	var members []*models.Member

	query := fmt.Sprintf(
		`SELECT u.id, u.nickname, u.avatar, per.read, per.write, per.admin,
		CASE p.owner_id
		WHEN user_id THEN true
		ELSE false
		END AS isOwner
		FROM %s AS pu
			INNER JOIN %s AS per ON pu.permissions_id = per.id
			INNER JOIN %s AS u ON pu.user_id = u.id
			INNER JOIN %s AS p ON pu.project_id = p.id
		WHERE pu.project_id = $1`,
		projectUsersTable, permissionsTable, usersTable, projectsTable)

	rows, err := r.db.Query(query, projectId)
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
