package postgres

import (
	"errors"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestUserPg_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserPg(db)

	tests := []struct {
		name    string
		mock    func()
		input   *models.User
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("TestName", "test@test.com", "password", "avatar").WillReturnRows(rows)
			},
			input: &models.User{
				Nickname: "TestName",
				Email:    "test@test.com",
				Password: "password",
				Avatar:   "avatar",
			},
			want: 1,
		},
		{
			name: "Empty Fields",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("TestName", "test@test.com", "", "avatar").WillReturnRows(rows)
			},
			input: &models.User{
				Nickname: "TestName",
				Email:    "test@test.com",
				Password: "",
				Avatar:   "avatar",
			},
			wantErr: true,
		},
		{
			name: "Not Unique Name",
			mock: func() {
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("NotUniqueName", "test@test.com", "password", "avatar").
					WillReturnError(errors.New("some error"))
			},
			input: &models.User{
				Nickname: "NotUniqueName",
				Email:    "test@test.com",
				Password: "password",
				Avatar:   "avatar",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.Create(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUserPg_Get(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserPg(db)

	type args struct {
		nickname string
		password string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    *models.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "nickname", "email", "password", "avatar"}).
					AddRow(1, "TestName", "test@test.com", "password", "avatar")
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("TestName", "password").WillReturnRows(rows)
			},
			input: args{"TestName", "password"},
			want: &models.User{
				Id:       1,
				Nickname: "TestName",
				Email:    "test@test.com",
				Password: "password",
				Avatar:   "avatar",
			},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "nickname", "email", "password", "avatar"})
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("not", "found").WillReturnRows(rows)
			},
			input:   args{"not", "found"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.Get(tt.input.nickname, tt.input.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUserPg_GetByNickname(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserPg(db)

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    *models.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "nickname", "email", "password", "avatar"}).
					AddRow(1, "TestName", "test@test.com", "password", "avatar")
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("TestName").WillReturnRows(rows)
			},
			input: "TestName",
			want: &models.User{
				Id:       1,
				Nickname: "TestName",
				Email:    "test@test.com",
				Password: "password",
				Avatar:   "avatar",
			},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "nickname", "email", "password", "avatar"})
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("NotFound").WillReturnRows(rows)
			},
			input:   "NotFound",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetByNickname(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
