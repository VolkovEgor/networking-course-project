package postgres

import (
	"errors"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestProjectPg_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewProjectPg(db)

	type args struct {
		project *models.Project
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
				project: &models.Project{
					OwnerId: 1,
					DefaultPermissions: &models.Permission{
						Read:  true,
						Write: true,
						Admin: true,
					},
					Datetimes: &models.Datetimes{
						Created:  1,
						Updated:  1,
						Accessed: 1,
					},
					Title:       "New Test Project",
					Description: "Some Description",
				},
			},
			want: 1,
			mock: func(args args, id int) {
				mock.ExpectBegin()

				defPermissionId := 1
				datetimesId := 1

				defPerm := args.project.DefaultPermissions
				defPermRows := sqlmock.NewRows([]string{"id"}).AddRow(defPermissionId)
				mock.ExpectQuery("INSERT INTO permissions").
					WithArgs(defPerm.Read, defPerm.Write, defPerm.Admin).
					WillReturnRows(defPermRows)

				datetimes := args.project.Datetimes
				dateRows := sqlmock.NewRows([]string{"id"}).AddRow(datetimesId)
				mock.ExpectQuery("INSERT INTO datetimes").
					WithArgs(datetimes.Created, datetimes.Updated, datetimes.Accessed).
					WillReturnRows(dateRows)

				project := args.project
				projectRows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO projects").WithArgs(project.OwnerId, defPermissionId,
					datetimesId, project.Title, project.Description).
					WillReturnRows(projectRows)

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

				mock.ExpectExec("INSERT INTO project_users").WithArgs(project.OwnerId, id, permissionId).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			input: args{
				project: &models.Project{
					OwnerId: 1,
					DefaultPermissions: &models.Permission{
						Read:  true,
						Write: true,
						Admin: true,
					},
					Datetimes: &models.Datetimes{
						Created:  1,
						Updated:  1,
						Accessed: 1,
					},
					Title:       "",
					Description: "Some Description",
				},
			},
			mock: func(args args, id int) {
				mock.ExpectBegin()

				defPermissionId := 1
				datetimesId := 1

				defPerm := args.project.DefaultPermissions
				defPermRows := sqlmock.NewRows([]string{"id"}).AddRow(defPermissionId)
				mock.ExpectQuery("INSERT INTO permissions").
					WithArgs(defPerm.Read, defPerm.Write, defPerm.Admin).
					WillReturnRows(defPermRows)

				datetimes := args.project.Datetimes
				dateRows := sqlmock.NewRows([]string{"id"}).AddRow(datetimesId)
				mock.ExpectQuery("INSERT INTO datetimes").
					WithArgs(datetimes.Created, datetimes.Updated, datetimes.Accessed).
					WillReturnRows(dateRows)

				project := args.project
				projectRows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(0, errors.New("insert error"))
				mock.ExpectQuery("INSERT INTO projects").WithArgs(project.OwnerId, defPermissionId,
					datetimesId, project.Title, project.Description).
					WillReturnRows(projectRows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Failed Insert",
			input: args{
				project: &models.Project{
					OwnerId: 1,
					DefaultPermissions: &models.Permission{
						Read:  true,
						Write: true,
						Admin: true,
					},
					Datetimes: &models.Datetimes{
						Created:  1,
						Updated:  1,
						Accessed: 1,
					},
					Title:       "New Test Project",
					Description: "Some Description",
				},
			},
			mock: func(args args, id int) {
				mock.ExpectBegin()

				defPermissionId := 1

				defPerm := args.project.DefaultPermissions
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

			got, err := r.Create(tt.input.project)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
