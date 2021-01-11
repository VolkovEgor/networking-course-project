package postgres

import (
	"errors"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTaskpermsPg_GetById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTaskPg(db)

	type args struct {
		listId int
	}
	type mockBehavior func(args args)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    *models.Task
		wantErr bool
	}{
		{
			name: "Ok",
			input: args{
				listId: 1,
			},
			want: &models.Task{
				Id:          1,
				ListId:      1,
				Title:       "title",
				Description: "description",
				Datetimes:   &models.Datetimes{1, 1, 1},
				Position:    1,
			},
			mock: func(args args) {
				rows := sqlmock.NewRows([]string{"t.id", "t.list_id", "t.title", "t.description", "d.created",
					"d.updated", "d.accessed", "t.position"}).AddRow(1, 1, "title", "description", 1, 1, 1, 1)
				mock.ExpectQuery("SELECT (.+) FROM tasks").WithArgs(args.listId).WillReturnRows(rows)
			},
		},
		{
			name: "error repo",
			input: args{
				listId: 1,
			},
			want: nil,
			mock: func(args args) {
				rows := sqlmock.NewRows([]string{"t.id", "t.list_id", "t.title", "t.description", "d.created",
					"d.updated", "d.accessed", "t.position"}).RowError(0, errors.New("Some error"))
				mock.ExpectQuery("SELECT (.+) FROM tasks").WithArgs(args.listId).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.GetById(tt.input.listId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
