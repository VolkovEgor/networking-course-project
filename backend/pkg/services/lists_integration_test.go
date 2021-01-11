package services

import (
	// "errors"
	"github.com/architectv/networking-course-project/backend/pkg/builders"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/repositories/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Integration_TaskListService_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	type args struct {
		userId    int
		projectId int
		boardId   int
		list      *models.TaskList
	}
	db, err := prepareTestDatabase()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer db.Close()

	rl := postgres.NewTaskListPg(db)
	rb := postgres.NewBoardPg(db)
	rp := postgres.NewProjectPg(db)
	s := NewTaskListService(rl, rb, rp)

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
				list:      builders.NewListBuilder().WithTitle("List Builder").Build(),
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusOK,
				Data: Map{"listId": 10001},
			},
		},
		{
			name: "Project Perm Failed",
			input: args{
				userId:    1,
				projectId: 3,
				boardId:   3,
				list:      builders.NewListBuilder().WithTitle("List Builder").Build(),
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
				list:      builders.NewListBuilder().WithTitle("List Builder").Build(),
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
				list:      builders.NewListBuilder().WithTitle("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA").Build(),
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := s.Create(test.input.userId, test.input.projectId, test.input.boardId, test.input.list)
			assert.Equal(t, test.expectedApiResponse.Code, got.Code)
			if test.expectedApiResponse.Code == StatusOK {
				assert.Equal(t, test.expectedApiResponse.Data, got.Data)
			}
		})
	}
}
