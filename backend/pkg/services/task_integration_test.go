package services

import (
	"github.com/architectv/networking-course-project/backend/pkg/builders"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/repositories/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Integration_TaskService_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	type args struct {
		userId    int
		projectId int
		boardId   int
		listId    int
		task      *models.Task
	}

	db, err := prepareTestDatabase()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer db.Close()

	rt := postgres.NewTaskPg(db)
	rb := postgres.NewBoardPg(db)
	rp := postgres.NewProjectPg(db)
	s := NewTaskService(rt, rb, rp)

	tests := []struct {
		name                string
		input               args
		expectedApiResponse *models.ApiResponse
	}{
		{
			name: "Ok",
			input: args{
				userId:    1,
				projectId: 1,
				boardId:   1,
				listId:    1,
				task:      builders.NewTaskBuilder().WithTitle("Task Builder").Build(),
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusOK,
				Data: Map{"taskId": 10001},
			},
		},
		{
			name: "Project Perm Failed",
			input: args{
				userId:    1,
				projectId: 3,
				boardId:   3,
				listId:    1,
				task:      builders.NewTaskBuilder().WithTitle("Task Builder").Build(),
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusForbidden,
			},
		},
		{
			name: "Board Perm Failed",
			input: args{
				userId:    2,
				projectId: 1,
				boardId:   1,
				listId:    1,
				task:      builders.NewTaskBuilder().WithTitle("Task Builder").Build(),
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusForbidden,
			},
		},
		{
			name: "Repo Error",
			input: args{
				userId:    1,
				projectId: 1,
				boardId:   1,
				listId:    1,
				task:      builders.NewTaskBuilder().WithTitle("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA").Build(),
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := s.Create(test.input.userId, test.input.projectId, test.input.boardId,
				test.input.listId, test.input.task)
			assert.Equal(t, test.expectedApiResponse.Code, got.Code)
			if test.expectedApiResponse.Code == StatusOK {
				assert.Equal(t, test.expectedApiResponse.Data, got.Data)
			}
		})
	}
}
