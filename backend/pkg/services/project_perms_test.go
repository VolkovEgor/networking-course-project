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

func TestProjectPermsService_Create(t *testing.T) {
	type args struct {
		userId         int
		projectId      int
		memberNickname string
		objectType     int
		perms          *models.Permission
		defPerms       *models.Permission
	}

	type projectMockBehavior func(r *mock_repositories.MockProject, projectId int)
	type getMockBehavior func(r *mock_repositories.MockObjectPerms, projectId, userId, objectType int)
	type mockBehavior func(r *mock_repositories.MockObjectPerms, projectId, objectType int, memberNickname string, permissions *models.Permission)

	tests := []struct {
		name                string
		input               args
		projectMock         projectMockBehavior
		getMock             getMockBehavior
		mock                mockBehavior
		expectedApiResponse *models.ApiResponse
	}{
		{
			name: "Ok",
			input: args{
				userId:         1,
				projectId:      1,
				memberNickname: "first_user",
				objectType:     IsProject,
				perms:          builders.NewPermsBuilder().Build(),
				defPerms:       builders.NewPermsBuilder().Build(),
			},
			projectMock: func(r *mock_repositories.MockProject, projectId int) {},
			getMock: func(r *mock_repositories.MockObjectPerms, projectId, userId, objectType int) {
				r.EXPECT().GetById(projectId, userId, objectType).Return(&models.Permission{true, true, true}, nil)
			},
			mock: func(r *mock_repositories.MockObjectPerms, projectId, objectType int, memberNickname string, permissions *models.Permission) {
				r.EXPECT().Create(projectId, objectType, memberNickname, permissions).Return(1, nil)
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusOK,
				Data: Map{"Project permissions id": 1},
			},
		},
		{
			name: "OK empty permissions",
			input: args{
				userId:         1,
				projectId:      1,
				memberNickname: "second_user",
				objectType:     IsProject,
				perms:          builders.NewPermsBuilder().WithPerm(false, false, false).Build(),
				defPerms:       builders.NewPermsBuilder().WithPerm(true, true, false).Build(),
			},
			projectMock: func(r *mock_repositories.MockProject, projectId int) {
				r.EXPECT().GetById(projectId).Return(&models.Project{1, 1, &models.Permission{true, true, false},
					&models.Datetimes{1, 1, 1}, "title", "description"}, nil)
			},
			getMock: func(r *mock_repositories.MockObjectPerms, projectId, userId, objectType int) {
				r.EXPECT().GetById(projectId, userId, objectType).Return(&models.Permission{true, true, true}, nil)
			},
			mock: func(r *mock_repositories.MockObjectPerms, projectId, objectType int, memberNickname string, permissions *models.Permission) {
				r.EXPECT().Create(projectId, objectType, memberNickname, permissions).Return(1, nil)
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusOK,
				Data: Map{"Project permissions id": 1},
			},
		},
		{
			name: "Repo error for Get in Project",
			input: args{
				userId:         1,
				projectId:      1,
				memberNickname: "second_user",
				objectType:     IsProject,
				perms:          builders.NewPermsBuilder().WithPerm(false, false, false).Build(),
				defPerms:       builders.NewPermsBuilder().Build(),
			},
			projectMock: func(r *mock_repositories.MockProject, projectId int) {
				r.EXPECT().GetById(projectId).Return(nil, errors.New("some error"))
			},
			getMock: func(r *mock_repositories.MockObjectPerms, projectId, userId, objectType int) {},
			mock: func(r *mock_repositories.MockObjectPerms, projectId, objectType int, memberNickname string, permissions *models.Permission) {
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},
		{
			name: "Default permissions is not defined",
			input: args{
				userId:         1,
				projectId:      1,
				memberNickname: "second_user",
				objectType:     IsProject,
				perms:          builders.NewPermsBuilder().WithPerm(false, false, false).Build(),
				defPerms:       builders.NewPermsBuilder().Build(),
			},
			projectMock: func(r *mock_repositories.MockProject, projectId int) {
				r.EXPECT().GetById(projectId).Return(&models.Project{1, 1, nil,
					&models.Datetimes{1, 1, 1}, "title", "description"}, nil)
			},
			getMock: func(r *mock_repositories.MockObjectPerms, projectId, userId, objectType int) {},
			mock: func(r *mock_repositories.MockObjectPerms, projectId, objectType int, memberNickname string, permissions *models.Permission) {
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},
		{
			name: "Permissions are set incorrectly",
			input: args{
				userId:         1,
				projectId:      1,
				memberNickname: "second_user",
				objectType:     IsProject,
				perms:          builders.NewPermsBuilder().WithPerm(true, false, true).Build(),
				defPerms:       builders.NewPermsBuilder().Build(),
			},
			projectMock: func(r *mock_repositories.MockProject, projectId int) {},
			getMock:     func(r *mock_repositories.MockObjectPerms, projectId, userId, objectType int) {},
			mock: func(r *mock_repositories.MockObjectPerms, projectId, objectType int, memberNickname string, permissions *models.Permission) {
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusBadRequest,
			},
		},
		{
			name: "Request author is not project member",
			input: args{
				userId:         1,
				projectId:      1,
				memberNickname: "second_user",
				objectType:     IsProject,
				perms:          builders.NewPermsBuilder().Build(),
				defPerms:       builders.NewPermsBuilder().Build(),
			},
			projectMock: func(r *mock_repositories.MockProject, projectId int) {},
			getMock: func(r *mock_repositories.MockObjectPerms, projectId, userId, objectType int) {
				r.EXPECT().GetById(projectId, userId, objectType).Return(nil, errors.New(DbResultNotFound))
			},
			mock: func(r *mock_repositories.MockObjectPerms, projectId, objectType int, memberNickname string, permissions *models.Permission) {
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusNotFound,
			},
		},
		{
			name: "Repo error for Get in ProjectPerms",
			input: args{
				userId:         1,
				projectId:      1,
				memberNickname: "second_user",
				objectType:     IsProject,
				perms:          builders.NewPermsBuilder().Build(),
				defPerms:       builders.NewPermsBuilder().Build(),
			},
			projectMock: func(r *mock_repositories.MockProject, projectId int) {},
			getMock: func(r *mock_repositories.MockObjectPerms, projectId, userId, objectType int) {
				r.EXPECT().GetById(projectId, userId, objectType).Return(nil, errors.New("some error"))
			},
			mock: func(r *mock_repositories.MockObjectPerms, projectId, objectType int, memberNickname string, permissions *models.Permission) {
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},
		{
			name: "Request author is not project admin",
			input: args{
				userId:         1,
				projectId:      1,
				memberNickname: "second_user",
				objectType:     IsProject,
				perms:          builders.NewPermsBuilder().Build(),
				defPerms:       builders.NewPermsBuilder().Build(),
			},
			projectMock: func(r *mock_repositories.MockProject, projectId int) {},
			getMock: func(r *mock_repositories.MockObjectPerms, projectId, userId, objectType int) {
				r.EXPECT().GetById(projectId, userId, objectType).Return(&models.Permission{true, true, false}, nil)
			},
			mock: func(r *mock_repositories.MockObjectPerms, projectId, objectType int, memberNickname string, permissions *models.Permission) {
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusForbidden,
			},
		},
		{
			name: "Repo error for Create in Project",
			input: args{
				userId:         1,
				projectId:      1,
				memberNickname: "second_user",
				objectType:     IsProject,
				perms:          builders.NewPermsBuilder().Build(),
				defPerms:       builders.NewPermsBuilder().Build(),
			},
			projectMock: func(r *mock_repositories.MockProject, projectId int) {},
			getMock: func(r *mock_repositories.MockObjectPerms, projectId, userId, objectType int) {
				r.EXPECT().GetById(projectId, userId, objectType).Return(&models.Permission{true, true, true}, nil)
			},
			mock: func(r *mock_repositories.MockObjectPerms, projectId, objectType int, memberNickname string, permissions *models.Permission) {
				r.EXPECT().Create(projectId, objectType, memberNickname, permissions).Return(0, errors.New("Some error"))
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

			repo := mock_repositories.NewMockObjectPerms(c)
			projectRepo := mock_repositories.NewMockProject(c)
			boardRepo := mock_repositories.NewMockBoard(c)

			test.projectMock(projectRepo, test.input.projectId)
			test.getMock(repo, test.input.projectId, test.input.userId, test.input.objectType)
			test.mock(repo, test.input.projectId, test.input.objectType, test.input.memberNickname,
				test.input.defPerms)
			s := &ProjectPermsService{repo: repo, projectRepo: projectRepo, boardRepo: boardRepo}

			got := s.Create(test.input.userId, test.input.projectId, test.input.memberNickname, test.input.perms)
			assert.Equal(t, test.expectedApiResponse.Code, got.Code)
			if test.expectedApiResponse.Code == StatusOK {
				assert.Equal(t, test.expectedApiResponse.Data, got.Data)
			}
		})
	}
}
