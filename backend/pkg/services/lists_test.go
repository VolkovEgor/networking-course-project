package services

import (
	"errors"
	"github.com/architectv/networking-course-project/backend/pkg/builders"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"testing"

	mock_repositories "github.com/architectv/networking-course-project/backend/pkg/repositories/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTaskListService_Create(t *testing.T) {
	type args struct {
		userId    int
		projectId int
		boardId   int
		list      *models.TaskList
	}
	type mockBehavior func(r *mock_repositories.MockTaskList, list *models.TaskList)
	type projectMockBehavior func(r *mock_repositories.MockProject, userId, projectId int)
	type boardMockBehavior func(r *mock_repositories.MockBoard, userId, boardId int)

	tests := []struct {
		name                string
		input               args
		projectMock         projectMockBehavior
		boardMock           boardMockBehavior
		mock                mockBehavior
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
			mock: func(r *mock_repositories.MockTaskList, list *models.TaskList) {
				r.EXPECT().Create(list).Return(1, nil)
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusOK,
				Data: Map{"listId": 1},
			},
		},
		{
			name: "Project Perm Failed",
			input: args{
				userId:    1,
				projectId: 1,
				boardId:   1,
				list:      builders.NewListBuilder().WithTitle("List Builder").Build(),
			},
			projectMock: func(r *mock_repositories.MockProject, userId, projectId int) {
				r.EXPECT().GetPermissions(userId, projectId).Return(nil, errors.New("Forbidden"))
			},
			boardMock: func(r *mock_repositories.MockBoard, userId, boardId int) {},
			mock:      func(r *mock_repositories.MockTaskList, list *models.TaskList) {},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusForbidden,
			},
		},
		{
			name: "Board Perm Failed",
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
				r.EXPECT().GetPermissions(userId, boardId).Return(nil, errors.New("Forbidden"))
			},
			mock: func(r *mock_repositories.MockTaskList, list *models.TaskList) {},
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
				list:      builders.NewListBuilder().WithTitle("List Builder").Build(),
			},
			projectMock: func(r *mock_repositories.MockProject, userId, projectId int) {
				r.EXPECT().GetPermissions(userId, projectId).Return(&models.Permission{true, true, true}, nil)
			},
			boardMock: func(r *mock_repositories.MockBoard, userId, boardId int) {
				r.EXPECT().GetPermissions(userId, boardId).Return(&models.Permission{true, true, true}, nil)
			},
			mock: func(r *mock_repositories.MockTaskList, list *models.TaskList) {
				r.EXPECT().Create(list).Return(0, errors.New("repo error"))
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

			repo := mock_repositories.NewMockTaskList(c)
			projectRepo := mock_repositories.NewMockProject(c)
			boardRepo := mock_repositories.NewMockBoard(c)
			test.projectMock(projectRepo, test.input.userId, test.input.projectId)
			test.boardMock(boardRepo, test.input.userId, test.input.boardId)
			test.mock(repo, test.input.list)
			s := &TaskListService{repo: repo, projectRepo: projectRepo, boardRepo: boardRepo}

			got := s.Create(test.input.userId, test.input.projectId, test.input.boardId, test.input.list)
			assert.Equal(t, test.expectedApiResponse.Code, got.Code)
			if test.expectedApiResponse.Code == StatusOK {
				assert.Equal(t, test.expectedApiResponse.Data, got.Data)
			}
		})
	}
}
