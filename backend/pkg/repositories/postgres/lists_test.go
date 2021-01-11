package postgres

import (
	"errors"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTaskListPg_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTaskListPg(db)

	type args struct {
		list *models.TaskList
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
				list: &models.TaskList{
					BoardId:  1,
					Title:    "Default List Title",
					Position: 1,
				},
			},
			want: 1,
			mock: func(args args, id int) {
				mock.ExpectBegin()
				position := 1
				rows := sqlmock.NewRows([]string{"maxPos"}).AddRow(position)
				mock.ExpectQuery(`SELECT (.+) FROM task_lists AS tl INNER JOIN boards AS b ON (.+) WHERE (.+)`).
					WithArgs(args.list.BoardId).
					WillReturnRows(rows)
				position++

				list := args.list
				rows = sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO task_lists").
					WithArgs(list.BoardId, list.Title, position).WillReturnRows(rows)

				mock.ExpectCommit()
			},
		},
		{
			name: "Failed Select",
			input: args{
				list: &models.TaskList{
					BoardId:  1,
					Title:    "",
					Position: 1,
				},
			},
			mock: func(args args, id int) {
				mock.ExpectBegin()
				position := 1
				rows := sqlmock.NewRows([]string{"maxPos"}).AddRow(position).RowError(0, errors.New("insert error"))
				mock.ExpectQuery(`SELECT (.+) FROM task_lists AS tl INNER JOIN boards AS b ON (.+) WHERE (.+)`).
					WithArgs(args.list.BoardId).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Failed Insert",
			input: args{
				list: &models.TaskList{
					BoardId:  1,
					Title:    "",
					Position: 1,
				},
			},
			mock: func(args args, id int) {
				mock.ExpectBegin()
				position := 1
				rows := sqlmock.NewRows([]string{"maxPos"}).AddRow(position)
				mock.ExpectQuery(`SELECT (.+) FROM task_lists AS tl INNER JOIN boards AS b ON (.+) WHERE (.+)`).
					WithArgs(args.list.BoardId).
					WillReturnRows(rows)
				position++

				list := args.list
				rows = sqlmock.NewRows([]string{"id"}).AddRow(1).RowError(0, errors.New("insert error"))
				mock.ExpectQuery("INSERT INTO task_lists").
					WithArgs(list.BoardId, list.Title, position).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input, tt.want)

			got, err := r.Create(tt.input.list)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
