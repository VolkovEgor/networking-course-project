package postgres

import (
	"errors"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestBoardPg_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewBoardPg(db)

	type args struct {
		userId int
		board  *models.Board
	}
	type mockBehavior func(args args, id int)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			input: args{
				board: &models.Board{
					ProjectId: 1,
					DefaultPermissions: &models.Permission{
						Read:  true,
						Write: true,
						Admin: false,
					},
					Datetimes: &models.Datetimes{
						Created:  1,
						Updated:  1,
						Accessed: 1,
					},
					Title: "New Test Board",
				},
			},
			want: 1,
			mock: func(args args, id int) {
				mock.ExpectBegin()

				defPermissionId := 1
				datetimesId := 1

				defPerm := args.board.DefaultPermissions
				defPermRows := sqlmock.NewRows([]string{"id"}).AddRow(defPermissionId)
				mock.ExpectQuery("INSERT INTO permissions").
					WithArgs(defPerm.Read, defPerm.Write, defPerm.Admin).
					WillReturnRows(defPermRows)

				datetimes := args.board.Datetimes
				dateRows := sqlmock.NewRows([]string{"id"}).AddRow(datetimesId)
				mock.ExpectQuery("INSERT INTO datetimes").
					WithArgs(datetimes.Created, datetimes.Updated, datetimes.Accessed).
					WillReturnRows(dateRows)

				board := args.board
				boardRows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO boards").WithArgs(board.ProjectId, board.OwnerId,
					defPermissionId, datetimesId, board.Title).
					WillReturnRows(boardRows)

				perm := &models.Permission{
					Read:  true,
					Write: true,
					Admin: true,
				}
				permissionId := defPermissionId + 1
				permRows := sqlmock.NewRows([]string{"id"}).AddRow(permissionId)
				mock.ExpectQuery("INSERT INTO permissions").
					WithArgs(perm.Read, perm.Write, perm.Admin).
					WillReturnRows(permRows)

				mock.ExpectExec("INSERT INTO board_users").WithArgs(args.userId, id, permissionId).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			input: args{
				board: &models.Board{
					ProjectId: 1,
					DefaultPermissions: &models.Permission{
						Read:  true,
						Write: true,
						Admin: false,
					},
					Datetimes: &models.Datetimes{
						Created:  1,
						Updated:  1,
						Accessed: 1,
					},
					Title: "New Test Board",
				},
			},
			mock: func(args args, id int) {
				mock.ExpectBegin()

				defPermissionId := 1
				datetimesId := 1

				defPerm := args.board.DefaultPermissions
				defPermRows := sqlmock.NewRows([]string{"id"}).AddRow(defPermissionId)
				mock.ExpectQuery("INSERT INTO permissions").
					WithArgs(defPerm.Read, defPerm.Write, defPerm.Admin).
					WillReturnRows(defPermRows)

				datetimes := args.board.Datetimes
				dateRows := sqlmock.NewRows([]string{"id"}).AddRow(datetimesId)
				mock.ExpectQuery("INSERT INTO datetimes").
					WithArgs(datetimes.Created, datetimes.Updated, datetimes.Accessed).
					WillReturnRows(dateRows)

				board := args.board
				boardRows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO boards").WithArgs(board.ProjectId, board.OwnerId,
					defPermissionId, datetimesId, board.Title).
					WillReturnRows(boardRows)

				perm := &models.Permission{
					Read:  true,
					Write: true,
					Admin: true,
				}
				permissionId := defPermissionId + 1
				permRows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(0, errors.New("insert error"))
				mock.ExpectQuery("INSERT INTO permissions").
					WithArgs(perm.Read, perm.Write, perm.Admin).
					WillReturnRows(permRows)

				mock.ExpectExec("INSERT INTO board_users").WithArgs(args.userId, id, permissionId).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Failed Insert",
			input: args{
				board: &models.Board{
					ProjectId: 1,
					DefaultPermissions: &models.Permission{
						Read:  true,
						Write: true,
						Admin: false,
					},
					Datetimes: &models.Datetimes{
						Created:  1,
						Updated:  1,
						Accessed: 1,
					},
					Title: "New Test Board",
				},
			},
			mock: func(args args, id int) {
				mock.ExpectBegin()

				defPermissionId := 1

				defPerm := args.board.DefaultPermissions
				defPermRows := sqlmock.NewRows([]string{"id"}).AddRow(defPermissionId).RowError(0, errors.New("insert error"))
				mock.ExpectQuery("INSERT INTO permissions").
					WithArgs(defPerm.Read, defPerm.Write, defPerm.Admin).
					WillReturnRows(defPermRows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input, tt.want)

			got, err := r.Create(tt.input.userId, tt.input.board)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
