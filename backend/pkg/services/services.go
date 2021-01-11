package services

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/repositories"
)

type User interface {
	GetAll() ([]*models.User, error)
	Get(id int) *models.ApiResponse
	Create(user *models.User) *models.ApiResponse
	GenerateToken(username, password string) *models.ApiResponse
	ParseToken(token string) (int, error)
	SignOut(token string) *models.ApiResponse
	Update(id int, profile *models.UpdateUser) *models.ApiResponse
}

type Project interface {
	Create(userId int, project *models.Project) *models.ApiResponse
	GetAll(userId int) *models.ApiResponse
	GetById(userId, projectId int) *models.ApiResponse
	Delete(userId, projectId int) *models.ApiResponse
	Update(userId, projectId int, project *models.UpdateProject) *models.ApiResponse
	GetMembers(userId, projectId int) *models.ApiResponse
}

type Board interface {
	Create(userId, projectId int, board *models.Board) *models.ApiResponse
	GetAll(userId, projectId int) *models.ApiResponse
	GetById(userId, projectId, boardId int) *models.ApiResponse
	Delete(userId, projectId, boardId int) *models.ApiResponse
	Update(userId, projectId, boardId int, board *models.UpdateBoard) *models.ApiResponse
	GetMembers(userId, projectId, boardId int) *models.ApiResponse
}

type TaskList interface {
	Create(userId, projectId, boardId int, list *models.TaskList) *models.ApiResponse
	GetAll(userId, projectId, boardId int) *models.ApiResponse
	GetById(userId, projectId, boardId, listId int) *models.ApiResponse
	Delete(userId, projectId, boardId, listId int) *models.ApiResponse
	Update(userId, projectId, boardId, listId int, list *models.UpdateTaskList) *models.ApiResponse
}

type Task interface {
	Create(userId, projectId, boardId, listId int, list *models.Task) *models.ApiResponse
	GetAll(userId, projectId, boardId, listId int) *models.ApiResponse
	GetById(userId, projectId, boardId, listId, taskId int) *models.ApiResponse
	Delete(userId, projectId, boardId, listId, taskId int) *models.ApiResponse
	Update(userId, projectId, boardId, listId, taskId int, list *models.UpdateTask) *models.ApiResponse
}

type Label interface {
	Create(userId, projectId, boardId int, label *models.Label) *models.ApiResponse
	CreateInTask(userId, projectId, boardId, taskId, labelId int) *models.ApiResponse
	GetAllInTask(userId, projectId, boardId, taskId int) *models.ApiResponse
	GetAll(userId, projectId, boardId int) *models.ApiResponse
	GetById(userId, projectId, boardId, labelId int) *models.ApiResponse
	DeleteInTask(userId, projectId, boardId, taskId, labelId int) *models.ApiResponse
	Delete(userId, projectId, boardId, labelId int) *models.ApiResponse
	Update(userId, projectId, boardId, labelId int, label *models.UpdateLabel) *models.ApiResponse
}

type UrlValidator interface {
	Validation(urlIds *models.UrlIds) *models.ApiResponse
}

type ProjectPerms interface {
	Create(userId, projectId int, memberNickname string, permissions *models.Permission) *models.ApiResponse
	Get(userId, projectId, memberId int) *models.ApiResponse
	Delete(userId, projectId, memberId int) *models.ApiResponse
	Update(userId, projectId, memberId int, list *models.UpdatePermission) *models.ApiResponse
}

type BoardPerms interface {
	Create(userId, projectId, boardId int, memberNickname string, permissions *models.Permission) *models.ApiResponse
	Get(userId, projectId, boardId, memberId int) *models.ApiResponse
	Delete(userId, projectId, boardId, memberId int) *models.ApiResponse
	Update(userId, projectId, boardId, memberId int, list *models.UpdatePermission) *models.ApiResponse
}

type Service struct {
	User
	Project
	Board
	TaskList
	Task
	Label
	UrlValidator
	ProjectPerms
	BoardPerms
}

func NewService(repos *repositories.Repository) *Service {
	return &Service{
		User:         NewUserService(repos.User),
		Project:      NewProjectService(repos.Project),
		Board:        NewBoardService(repos.Board, repos.Project),
		TaskList:     NewTaskListService(repos.TaskList, repos.Board, repos.Project),
		Task:         NewTaskService(repos.Task, repos.Board, repos.Project),
		Label:        NewLabelService(repos.Label, repos.Board, repos.Project),
		UrlValidator: NewUrlValidatorService(repos.Board, repos.TaskList, repos.Task),
		ProjectPerms: NewProjectPermsService(repos.ObjectPerms, repos.Project, repos.Board),
		BoardPerms:   NewBoardPermsService(repos.ObjectPerms, repos.Board, repos.Project),
	}
}
