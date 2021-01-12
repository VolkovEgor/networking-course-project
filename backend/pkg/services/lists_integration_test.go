// +build integration

package services

import (
	"errors"
	"github.com/architectv/networking-course-project/backend/pkg/builders"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/repositories/postgres"
	"testing"

	mock_repositories "github.com/architectv/networking-course-project/backend/pkg/repositories/mocks"

	"github.com/golang/mock/gomock"
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

	type projectMockBehavior func(r *mock_repositories.MockProject, userId, projectId int)
	type boardMockBehavior func(r *mock_repositories.MockBoard, userId, boardId int)

	tests := []struct {
		name                string
		input               args
		projectMock         projectMockBehavior
		boardMock           boardMockBehavior
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
			projectMock: func(r *mock_repositories.MockProject, userId, projectId int) {
				r.EXPECT().GetPermissions(userId, projectId).Return(&models.Permission{true, true, true}, nil)
			},
			boardMock: func(r *mock_repositories.MockBoard, userId, boardId int) {
				r.EXPECT().GetPermissions(userId, boardId).Return(&models.Permission{true, true, true}, nil)
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
			projectMock: func(r *mock_repositories.MockProject, userId, projectId int) {
				r.EXPECT().GetPermissions(userId, projectId).Return(nil, errors.New("Forbidden"))
			},
			boardMock: func(r *mock_repositories.MockBoard, userId, boardId int) {},
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
			projectMock: func(r *mock_repositories.MockProject, userId, projectId int) {
				r.EXPECT().GetPermissions(userId, projectId).Return(&models.Permission{true, true, true}, nil)
			},
			boardMock: func(r *mock_repositories.MockBoard, userId, boardId int) {
				r.EXPECT().GetPermissions(userId, boardId).Return(nil, errors.New("Forbidden"))
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
			projectMock: func(r *mock_repositories.MockProject, userId, projectId int) {
				r.EXPECT().GetPermissions(userId, projectId).Return(&models.Permission{true, true, true}, nil)
			},
			boardMock: func(r *mock_repositories.MockBoard, userId, boardId int) {
				r.EXPECT().GetPermissions(userId, boardId).Return(&models.Permission{true, true, true}, nil)
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := postgres.NewTaskListPg(db)
			projectRepo := mock_repositories.NewMockProject(c)
			boardRepo := mock_repositories.NewMockBoard(c)
			test.projectMock(projectRepo, test.input.userId, test.input.projectId)
			test.boardMock(boardRepo, test.input.userId, test.input.boardId)
			s := NewTaskListService(repo, boardRepo, projectRepo)

			got := s.Create(test.input.userId, test.input.projectId, test.input.boardId, test.input.list)
			assert.Equal(t, test.expectedApiResponse.Code, got.Code)
			if test.expectedApiResponse.Code == StatusOK {
				assert.Equal(t, test.expectedApiResponse.Data, got.Data)
			}
		})
	}
}
