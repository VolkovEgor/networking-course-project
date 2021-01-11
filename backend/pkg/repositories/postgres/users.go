package postgres

import (
	"fmt"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

type UserPg struct {
	db *sqlx.DB
}

func NewUserPg(db *sqlx.DB) *UserPg {
	return &UserPg{db: db}
}

func (r *UserPg) GetById(id int) (*models.User, error) {
	user := &models.User{}

	query := fmt.Sprintf(
		`SELECT *
		FROM %s AS u
		WHERE u.id = $1`,
		usersTable)

	if err := r.db.Get(user, query, id); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserPg) Update(id int, profile *models.UpdateUser) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if profile.Nickname != nil {
		setValues = append(setValues, fmt.Sprintf("nickname=$%d", argId))
		args = append(args, *profile.Nickname)
		argId++
	}

	if profile.Firstname != nil {
		setValues = append(setValues, fmt.Sprintf("firstname=$%d", argId))
		args = append(args, *profile.Firstname)
		argId++
	}

	if profile.Lastname != nil {
		setValues = append(setValues, fmt.Sprintf("lastname=$%d", argId))
		args = append(args, *profile.Lastname)
		argId++
	}

	if profile.Email != nil {
		setValues = append(setValues, fmt.Sprintf("email=$%d", argId))
		args = append(args, *profile.Email)
		argId++
	}

	if profile.Phone != nil {
		setValues = append(setValues, fmt.Sprintf("phone=$%d", argId))
		args = append(args, *profile.Phone)
		argId++
	}

	if profile.Avatar != nil {
		setValues = append(setValues, fmt.Sprintf("avatar=$%d", argId))
		args = append(args, *profile.Avatar)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s SET %s where id=$%d`,
		usersTable, setQuery, argId)
	args = append(args, id)
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return err
}

func (r *UserPg) GetAll() ([]*models.User, error) {
	var users []*models.User
	query := fmt.Sprintf(`SELECT * FROM %s`, usersTable)
	if err := r.db.Select(&users, query); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserPg) Get(nickname, password string) (*models.User, error) {
	user := &models.User{}
	query := fmt.Sprintf(
		`SELECT * FROM %s WHERE nickname = $1 AND password = $2`, usersTable)
	err := r.db.Get(user, query, nickname, password)

	return user, err
}

func (r *UserPg) Create(user *models.User) (int, error) {
	var id int
	query := fmt.Sprintf(
		`INSERT INTO %s (nickname, email, password, avatar)
		VALUES ($1, $2, $3, $4) RETURNING id`, usersTable)

	row := r.db.QueryRow(query, user.Nickname, user.Email, user.Password, user.Avatar)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserPg) GetByNickname(nickname string) (*models.User, error) {
	user := &models.User{}
	query := fmt.Sprintf(`SELECT * FROM %s WHERE nickname = $1`, usersTable)
	err := r.db.Get(user, query, nickname)

	return user, err
}

func (r *UserPg) SignOut(token string) (int, error) {
	var id int
	query := fmt.Sprintf(
		`INSERT INTO %s (jwt)
		VALUES ($1) RETURNING id`, tokensTable)

	row := r.db.QueryRow(query, token)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserPg) FindToken(input string) error {
	type tokenDB struct {
		Id  int    `json:"id"`
		Jwt string `json:"jwt"`
	}
	token := &tokenDB{}
	query := fmt.Sprintf(`SELECT * FROM %s WHERE jwt = $1`, tokensTable)
	err := r.db.Get(token, query, input)

	return err
}
