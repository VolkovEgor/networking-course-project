package services

import (
	"errors"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"testing"

	mock_repositories "github.com/architectv/networking-course-project/backend/pkg/repositories/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type UrlIdsBuilder struct {
	UrlIds *models.UrlIds
}

func NewUrlIdsBuilder() *UrlIdsBuilder {
	urlIds := &models.UrlIds{
		ProjectId: 0,
		BoardId:   0,
		ListId:    0,
		TaskId:    0,
	}
	return &UrlIdsBuilder{UrlIds: urlIds}
}

func (u *UrlIdsBuilder) build() *models.UrlIds {
	return u.UrlIds
}

func (u *UrlIdsBuilder) withProject(id int) *UrlIdsBuilder {
	u.UrlIds.ProjectId = id
	return u
}

func (u *UrlIdsBuilder) withBoard(id int) *UrlIdsBuilder {
	u.UrlIds.BoardId = id
	return u
}

func (u *UrlIdsBuilder) withList(id int) *UrlIdsBuilder {
	u.UrlIds.ListId = id
	return u
}

func (u *UrlIdsBuilder) withTask(id int) *UrlIdsBuilder {
	u.UrlIds.TaskId = id
	return u
}
func TestUrlValidatorService_Create(t *testing.T) {
	type args struct {
		urlIds *models.UrlIds
	}
	type boardMockBehavior func(r *mock_repositories.MockBoard, boardId int)
	type listMockBehavior func(r *mock_repositories.MockTaskList, listId int)
	type taskBehavior func(r *mock_repositories.MockTask, taskId int)

	tests := []struct {
		name                string
		input               args
		boardMock           boardMockBehavior
		listMock            listMockBehavior
		taskMock            taskBehavior
		expectedApiResponse *models.ApiResponse
	}{
		{
			name: "Ok",
			input: args{
				urlIds: NewUrlIdsBuilder().withProject(1).withBoard(1).withList(1).withTask(1).build(),
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {
				r.EXPECT().GetById(boardId).Return(&models.Board{1, 1, 1, &models.Permission{true, true, false},
					&models.Datetimes{1, 1, 1}, "title"}, nil)
			},
			listMock: func(r *mock_repositories.MockTaskList, listId int) {
				r.EXPECT().GetById(listId).Return(&models.TaskList{1, 1, "title", 1}, nil)
			},
			taskMock: func(r *mock_repositories.MockTask, taskId int) {
				r.EXPECT().GetById(taskId).Return(&models.Task{1, 1, "title", "description", &models.Datetimes{1, 1, 1}, 1}, nil)
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusOK,
				Data: Map{},
			},
		},
		{
			name: "Board is not defined",
			input: args{
				urlIds: NewUrlIdsBuilder().withProject(1).withBoard(-1).build(),
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {
				r.EXPECT().GetById(boardId).Return(nil, errors.New(DbResultNotFound))
			},
			listMock: func(r *mock_repositories.MockTaskList, listId int) {},
			taskMock: func(r *mock_repositories.MockTask, taskId int) {},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusNotFound,
			},
		},
		{
			name: "Repo error for Get in Board",
			input: args{
				urlIds: NewUrlIdsBuilder().withProject(1).withBoard(1).build(),
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {
				r.EXPECT().GetById(boardId).Return(nil, errors.New("Some error"))
			},
			listMock: func(r *mock_repositories.MockTaskList, listId int) {},
			taskMock: func(r *mock_repositories.MockTask, taskId int) {},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},
		{
			name: "There is no requested board inside the project",
			input: args{
				urlIds: NewUrlIdsBuilder().withProject(1).withBoard(1).build(),
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {
				r.EXPECT().GetById(boardId).Return(&models.Board{1, 2, 1, &models.Permission{true, true, false},
					&models.Datetimes{1, 1, 1}, "title"}, nil)
			},
			listMock: func(r *mock_repositories.MockTaskList, listId int) {},
			taskMock: func(r *mock_repositories.MockTask, taskId int) {},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusNotFound,
			},
		},
		{
			name: "List is not defined",
			input: args{
				urlIds: NewUrlIdsBuilder().withProject(1).withBoard(1).withList(-1).build(),
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {
				r.EXPECT().GetById(boardId).Return(&models.Board{1, 1, 1, &models.Permission{true, true, false},
					&models.Datetimes{1, 1, 1}, "title"}, nil)
			},
			listMock: func(r *mock_repositories.MockTaskList, listId int) {
				r.EXPECT().GetById(listId).Return(nil, errors.New(DbResultNotFound))
			},
			taskMock: func(r *mock_repositories.MockTask, taskId int) {},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusNotFound,
			},
		},
		{
			name: "Repo error for Get in List",
			input: args{
				urlIds: NewUrlIdsBuilder().withProject(1).withBoard(1).withList(1).build(),
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {
				r.EXPECT().GetById(boardId).Return(&models.Board{1, 1, 1, &models.Permission{true, true, false},
					&models.Datetimes{1, 1, 1}, "title"}, nil)
			},
			listMock: func(r *mock_repositories.MockTaskList, listId int) {
				r.EXPECT().GetById(listId).Return(nil, errors.New("Some error"))
			},
			taskMock: func(r *mock_repositories.MockTask, taskId int) {},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},
		{
			name: "There is no requested list inside the board",
			input: args{
				urlIds: NewUrlIdsBuilder().withProject(1).withBoard(1).withList(1).build(),
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {
				r.EXPECT().GetById(boardId).Return(&models.Board{1, 1, 1, &models.Permission{true, true, false},
					&models.Datetimes{1, 1, 1}, "title"}, nil)
			},
			listMock: func(r *mock_repositories.MockTaskList, listId int) {
				r.EXPECT().GetById(listId).Return(&models.TaskList{1, 2, "title", 1}, nil)
			},
			taskMock: func(r *mock_repositories.MockTask, taskId int) {},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusNotFound,
			},
		},
		{
			name: "Task is not defined",
			input: args{
				urlIds: NewUrlIdsBuilder().withProject(1).withBoard(1).withList(1).withTask(-1).build(),
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {
				r.EXPECT().GetById(boardId).Return(&models.Board{1, 1, 1, &models.Permission{true, true, false},
					&models.Datetimes{1, 1, 1}, "title"}, nil)
			},
			listMock: func(r *mock_repositories.MockTaskList, listId int) {
				r.EXPECT().GetById(listId).Return(&models.TaskList{1, 1, "title", 1}, nil)
			},
			taskMock: func(r *mock_repositories.MockTask, taskId int) {
				r.EXPECT().GetById(taskId).Return(nil, errors.New(DbResultNotFound))
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusNotFound,
			},
		},
		{
			name: "Repo error for Get in Task",
			input: args{
				urlIds: NewUrlIdsBuilder().withProject(1).withBoard(1).withList(1).withTask(1).build(),
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {
				r.EXPECT().GetById(boardId).Return(&models.Board{1, 1, 1, &models.Permission{true, true, false},
					&models.Datetimes{1, 1, 1}, "title"}, nil)
			},
			listMock: func(r *mock_repositories.MockTaskList, listId int) {
				r.EXPECT().GetById(listId).Return(&models.TaskList{1, 1, "title", 1}, nil)
			},
			taskMock: func(r *mock_repositories.MockTask, taskId int) {
				r.EXPECT().GetById(taskId).Return(nil, errors.New("Some error"))
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},
		{
			name: "There is no requested task inside the list",
			input: args{
				urlIds: NewUrlIdsBuilder().withProject(1).withBoard(1).withList(1).withTask(1).build(),
			},
			boardMock: func(r *mock_repositories.MockBoard, boardId int) {
				r.EXPECT().GetById(boardId).Return(&models.Board{1, 1, 1, &models.Permission{true, true, false},
					&models.Datetimes{1, 1, 1}, "title"}, nil)
			},
			listMock: func(r *mock_repositories.MockTaskList, listId int) {
				r.EXPECT().GetById(listId).Return(&models.TaskList{1, 1, "title", 1}, nil)
			},
			taskMock: func(r *mock_repositories.MockTask, taskId int) {
				r.EXPECT().GetById(taskId).Return(&models.Task{1, 2, "title", "description", &models.Datetimes{1, 1, 1}, 1}, nil)
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusNotFound,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			boardRepo := mock_repositories.NewMockBoard(c)
			listRepo := mock_repositories.NewMockTaskList(c)
			taskRepo := mock_repositories.NewMockTask(c)

			test.boardMock(boardRepo, test.input.urlIds.BoardId)
			test.listMock(listRepo, test.input.urlIds.ListId)
			test.taskMock(taskRepo, test.input.urlIds.TaskId)
			s := &UrlValidatorService{boardRepo: boardRepo, listRepo: listRepo, taskRepo: taskRepo}

			got := s.Validation(test.input.urlIds)
			assert.Equal(t, test.expectedApiResponse.Code, got.Code)
			if test.expectedApiResponse.Code == StatusOK {
				assert.Equal(t, test.expectedApiResponse.Data, got.Data)
			}
		})
	}
}
