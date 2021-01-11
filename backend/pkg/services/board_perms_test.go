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

func TestBoardPermsService_Create(t *testing.T) {
	type args struct {
		userId         int
		projectId      int
		boardId        int
		memberNickname string
		boardType      int
		projectType    int
		perms          *models.Permission
		defPerms       *models.Permission
	}

	type boardMockBehavior func(r *mock_repositories.MockBoard, boardId int)
	type getCallerProjectPermMockBehavior func(r *mock_repositories.MockObjectPerms, projectId, userId, projectType int)
	type getCallerBoardPermMockBehavior func(r *mock_repositories.MockObjectPerms, boardId, userId, boardType int)
	type getMemberProjectPermMockBehavior func(r *mock_repositories.MockObjectPerms, projectId, projectType int, memberNickname string)
	type mockBehavior func(r *mock_repositories.MockObjectPerms, projectId, objectType int, memberNickname string, permissions *models.Permission)

	tests := []struct {
		name                 string
		input                args
		boardMock            boardMockBehavior
		getCallerProjectPerm getCallerProjectPermMockBehavior
		getCallerBoardPerm   getCallerBoardPermMockBehavior
		getMemberProjectPerm getMemberProjectPermMockBehavior
		mock                 mockBehavior
		expectedApiResponse  *models.ApiResponse
	}{
		{
			name: "Ok",
			input: args{
				userId:         1,
				projectId:      1,
				boardId:        1,
				memberNickname: "first_user",
				boardType:      IsBoard,
				projectType:    IsProject,
				perms:          builders.NewPermsBuilder().Build(),
				defPerms:       builders.NewPermsBuilder().Build(),
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {},
			getCallerProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, userId, projectType int) {
				r.EXPECT().GetById(projectId, userId, projectType).Return(&models.Permission{true, true, true}, nil)
			},
			getCallerBoardPerm: func(r *mock_repositories.MockObjectPerms, boardId, userId, boardType int) {
				r.EXPECT().GetById(boardId, userId, boardType).Return(&models.Permission{true, true, true}, nil)
			},
			getMemberProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, projectType int, memberNickname string) {
				r.EXPECT().GetByNickname(projectId, projectType, memberNickname).Return(&models.Permission{true, true, false}, nil)
			},
			mock: func(r *mock_repositories.MockObjectPerms, boardId, objectType int, memberNickname string, permissions *models.Permission) {
				r.EXPECT().Create(boardId, objectType, memberNickname, permissions).Return(1, nil)
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusOK,
				Data: Map{"Board permissions id": 1},
			},
		},
		{
			name: "OK empty permissions",
			input: args{
				userId:         1,
				projectId:      1,
				boardId:        1,
				memberNickname: "second_user",
				boardType:      IsBoard,
				projectType:    IsProject,
				perms:          builders.NewPermsBuilder().WithPerm(false, false, false).Build(),
				defPerms:       builders.NewPermsBuilder().WithPerm(true, true, false).Build(),
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {
				r.EXPECT().GetById(boardId).Return(&models.Board{1, 1, 1, &models.Permission{true, true, false},
					&models.Datetimes{1, 1, 1}, "title"}, nil)
			},
			getCallerProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, userId, projectType int) {
				r.EXPECT().GetById(projectId, userId, projectType).Return(&models.Permission{true, true, true}, nil)
			},
			getCallerBoardPerm: func(r *mock_repositories.MockObjectPerms, boardId, userId, boardType int) {
				r.EXPECT().GetById(boardId, userId, boardType).Return(&models.Permission{true, true, true}, nil)
			},
			getMemberProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, projectType int, memberNickname string) {
				r.EXPECT().GetByNickname(projectId, projectType, memberNickname).Return(&models.Permission{true, true, false}, nil)
			},
			mock: func(r *mock_repositories.MockObjectPerms, boardId, objectType int, memberNickname string, permissions *models.Permission) {
				r.EXPECT().Create(boardId, objectType, memberNickname, permissions).Return(1, nil)
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusOK,
				Data: Map{"Board permissions id": 1},
			},
		},
		{
			name: "Repo error for Get in Board",
			input: args{
				userId:         1,
				projectId:      1,
				boardId:        1,
				memberNickname: "second_user",
				boardType:      IsBoard,
				projectType:    IsProject,
				perms:          builders.NewPermsBuilder().WithPerm(false, false, false).Build(),
				defPerms:       nil,
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {
				r.EXPECT().GetById(boardId).Return(nil, errors.New("some error"))
			},
			getCallerProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, userId, projectType int) {},
			getCallerBoardPerm:   func(r *mock_repositories.MockObjectPerms, boardId, userId, boardType int) {},
			getMemberProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, projectType int, memberNickname string) {},
			mock: func(r *mock_repositories.MockObjectPerms, boardId, objectType int, memberNickname string, permissions *models.Permission) {
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
				boardId:        1,
				memberNickname: "second_user",
				boardType:      IsBoard,
				projectType:    IsProject,
				perms:          builders.NewPermsBuilder().WithPerm(false, false, false).Build(),
				defPerms:       nil,
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {
				r.EXPECT().GetById(boardId).Return(&models.Board{1, 1, 1, nil,
					&models.Datetimes{1, 1, 1}, "title"}, nil)
			},
			getCallerProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, userId, projectType int) {},
			getCallerBoardPerm:   func(r *mock_repositories.MockObjectPerms, boardId, userId, boardType int) {},
			getMemberProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, projectType int, memberNickname string) {},
			mock: func(r *mock_repositories.MockObjectPerms, boardId, objectType int, memberNickname string, permissions *models.Permission) {
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
				boardId:        1,
				memberNickname: "second_user",
				boardType:      IsBoard,
				projectType:    IsProject,
				perms:          builders.NewPermsBuilder().WithPerm(true, false, true).Build(),
				defPerms:       nil,
			},
			boardMock:            func(r *mock_repositories.MockBoard, boardId int) {},
			getCallerProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, userId, projectType int) {},
			getCallerBoardPerm:   func(r *mock_repositories.MockObjectPerms, boardId, userId, boardType int) {},
			getMemberProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, projectType int, memberNickname string) {},
			mock: func(r *mock_repositories.MockObjectPerms, boardId, objectType int, memberNickname string, permissions *models.Permission) {
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
				boardId:        1,
				memberNickname: "second_user",
				boardType:      IsBoard,
				projectType:    IsProject,
				perms:          builders.NewPermsBuilder().Build(),
				defPerms:       nil,
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {},
			getCallerProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, userId, projectType int) {
				r.EXPECT().GetById(projectId, userId, projectType).Return(nil, errors.New(DbResultNotFound))
			},
			getCallerBoardPerm:   func(r *mock_repositories.MockObjectPerms, boardId, userId, boardType int) {},
			getMemberProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, projectType int, memberNickname string) {},
			mock: func(r *mock_repositories.MockObjectPerms, boardId, objectType int, memberNickname string, permissions *models.Permission) {
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusNotFound,
			},
		},
		{
			name: "Repo error for Get in caller ProjectPerms",
			input: args{
				userId:         1,
				projectId:      1,
				boardId:        1,
				memberNickname: "second_user",
				boardType:      IsBoard,
				projectType:    IsProject,
				perms:          builders.NewPermsBuilder().Build(),
				defPerms:       nil,
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {},
			getCallerProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, userId, projectType int) {
				r.EXPECT().GetById(projectId, userId, projectType).Return(nil, errors.New("Some error"))
			},
			getCallerBoardPerm:   func(r *mock_repositories.MockObjectPerms, boardId, userId, boardType int) {},
			getMemberProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, projectType int, memberNickname string) {},
			mock: func(r *mock_repositories.MockObjectPerms, boardId, objectType int, memberNickname string, permissions *models.Permission) {
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},
		{
			name: "Request author is not board member",
			input: args{
				userId:         1,
				projectId:      1,
				boardId:        1,
				memberNickname: "second_user",
				boardType:      IsBoard,
				projectType:    IsProject,
				perms:          builders.NewPermsBuilder().Build(),
				defPerms:       nil,
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {},
			getCallerProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, userId, projectType int) {
				r.EXPECT().GetById(projectId, userId, projectType).Return(&models.Permission{true, true, true}, nil)
			},
			getCallerBoardPerm: func(r *mock_repositories.MockObjectPerms, boardId, userId, boardType int) {
				r.EXPECT().GetById(boardId, userId, boardType).Return(nil, errors.New(DbResultNotFound))
			},
			getMemberProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, projectType int, memberNickname string) {},
			mock: func(r *mock_repositories.MockObjectPerms, boardId, objectType int, memberNickname string, permissions *models.Permission) {
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusNotFound,
			},
		},
		{
			name: "Repo error for Get in caller BoardPerms",
			input: args{
				userId:         1,
				projectId:      1,
				boardId:        1,
				memberNickname: "second_user",
				boardType:      IsBoard,
				projectType:    IsProject,
				perms:          builders.NewPermsBuilder().Build(),
				defPerms:       nil,
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {},
			getCallerProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, userId, projectType int) {
				r.EXPECT().GetById(projectId, userId, projectType).Return(&models.Permission{true, true, true}, nil)
			},
			getCallerBoardPerm: func(r *mock_repositories.MockObjectPerms, boardId, userId, boardType int) {
				r.EXPECT().GetById(boardId, userId, boardType).Return(nil, errors.New("Some error"))
			},
			getMemberProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, projectType int, memberNickname string) {},
			mock: func(r *mock_repositories.MockObjectPerms, boardId, objectType int, memberNickname string, permissions *models.Permission) {
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},
		{
			name: "Request author is not board admin",
			input: args{
				userId:         1,
				projectId:      1,
				boardId:        1,
				memberNickname: "second_user",
				boardType:      IsBoard,
				projectType:    IsProject,
				perms:          builders.NewPermsBuilder().Build(),
				defPerms:       nil,
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {},
			getCallerProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, userId, projectType int) {
				r.EXPECT().GetById(projectId, userId, projectType).Return(&models.Permission{true, true, true}, nil)
			},
			getCallerBoardPerm: func(r *mock_repositories.MockObjectPerms, boardId, userId, boardType int) {
				r.EXPECT().GetById(boardId, userId, boardType).Return(&models.Permission{true, true, false}, nil)
			},
			getMemberProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, projectType int, memberNickname string) {},
			mock: func(r *mock_repositories.MockObjectPerms, boardId, objectType int, memberNickname string, permissions *models.Permission) {
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusForbidden,
			},
		},
		{
			name: "New board member is not project member",
			input: args{
				userId:         1,
				projectId:      1,
				boardId:        1,
				memberNickname: "second_user",
				boardType:      IsBoard,
				projectType:    IsProject,
				perms:          builders.NewPermsBuilder().Build(),
				defPerms:       nil,
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {},
			getCallerProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, userId, projectType int) {
				r.EXPECT().GetById(projectId, userId, projectType).Return(&models.Permission{true, true, true}, nil)
			},
			getCallerBoardPerm: func(r *mock_repositories.MockObjectPerms, boardId, userId, boardType int) {
				r.EXPECT().GetById(boardId, userId, boardType).Return(&models.Permission{true, true, true}, nil)
			},
			getMemberProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, projectType int, memberNickname string) {
				r.EXPECT().GetByNickname(projectId, projectType, memberNickname).Return(nil, errors.New(DbResultNotFound))
			},
			mock: func(r *mock_repositories.MockObjectPerms, boardId, objectType int, memberNickname string, permissions *models.Permission) {
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusNotFound,
			},
		},
		{
			name: "Repo error for Get in member BoardPerms",
			input: args{
				userId:         1,
				projectId:      1,
				boardId:        1,
				memberNickname: "second_user",
				boardType:      IsBoard,
				projectType:    IsProject,
				perms:          builders.NewPermsBuilder().Build(),
				defPerms:       nil,
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {},
			getCallerProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, userId, projectType int) {
				r.EXPECT().GetById(projectId, userId, projectType).Return(&models.Permission{true, true, true}, nil)
			},
			getCallerBoardPerm: func(r *mock_repositories.MockObjectPerms, boardId, userId, boardType int) {
				r.EXPECT().GetById(boardId, userId, boardType).Return(&models.Permission{true, true, true}, nil)
			},
			getMemberProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, projectType int, memberNickname string) {
				r.EXPECT().GetByNickname(projectId, projectType, memberNickname).Return(nil, errors.New("Some error"))
			},
			mock: func(r *mock_repositories.MockObjectPerms, boardId, objectType int, memberNickname string, permissions *models.Permission) {
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},

		{
			name: "Repo error for Create in Board",
			input: args{
				userId:         1,
				projectId:      1,
				boardId:        1,
				memberNickname: "second_user",
				boardType:      IsBoard,
				projectType:    IsProject,
				perms:          builders.NewPermsBuilder().Build(),
				defPerms:       builders.NewPermsBuilder().Build(),
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {},
			getCallerProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, userId, projectType int) {
				r.EXPECT().GetById(projectId, userId, projectType).Return(&models.Permission{true, true, true}, nil)
			},
			getCallerBoardPerm: func(r *mock_repositories.MockObjectPerms, boardId, userId, boardType int) {
				r.EXPECT().GetById(boardId, userId, boardType).Return(&models.Permission{true, true, true}, nil)
			},
			getMemberProjectPerm: func(r *mock_repositories.MockObjectPerms, projectId, projectType int, memberNickname string) {
				r.EXPECT().GetByNickname(projectId, projectType, memberNickname).Return(&models.Permission{true, true, false}, nil)
			},
			mock: func(r *mock_repositories.MockObjectPerms, boardId, objectType int, memberNickname string, permissions *models.Permission) {
				r.EXPECT().Create(boardId, objectType, memberNickname, permissions).Return(0, errors.New("Some error"))
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

			test.boardMock(boardRepo, test.input.boardId)
			test.getCallerProjectPerm(repo, test.input.projectId, test.input.userId, test.input.projectType)
			test.getCallerBoardPerm(repo, test.input.boardId, test.input.userId, test.input.boardType)
			test.getMemberProjectPerm(repo, test.input.projectId, test.input.projectType, test.input.memberNickname)
			test.mock(repo, test.input.boardId, test.input.boardType, test.input.memberNickname,
				test.input.defPerms)
			s := &BoardPermsService{repo: repo, projectRepo: projectRepo, boardRepo: boardRepo}

			got := s.Create(test.input.userId, test.input.projectId, test.input.boardId,
				test.input.memberNickname, test.input.perms)
			assert.Equal(t, test.expectedApiResponse.Code, got.Code)
			if test.expectedApiResponse.Code == StatusOK {
				assert.Equal(t, test.expectedApiResponse.Data, got.Data)
			}
		})
	}
}
