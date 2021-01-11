package services

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/repositories"
	"time"
)

type TaskService struct {
	repo        repositories.Task
	boardRepo   repositories.Board
	projectRepo repositories.Project
}

func NewTaskService(repo repositories.Task, boardRepo repositories.Board, projectRepo repositories.Project) *TaskService {
	return &TaskService{repo: repo, boardRepo: boardRepo, projectRepo: projectRepo}
}

func (s *TaskService) GetAll(userId, projectId, boardId, listId int) *models.ApiResponse {
	r := &models.ApiResponse{}

	projectPermissions, err := s.projectRepo.GetPermissions(userId, projectId)
	if err != nil || projectPermissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	boardPermissions, err := s.boardRepo.GetPermissions(userId, boardId)
	if err != nil || boardPermissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	tasks, err := s.repo.GetAll(listId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"tasks": tasks})
	return r
}

func (s *TaskService) GetById(userId, projectId, boardId, listId, taksId int) *models.ApiResponse {
	r := &models.ApiResponse{}
	projectPermissions, err := s.projectRepo.GetPermissions(userId, projectId)
	if err != nil || projectPermissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	boardPermissions, err := s.boardRepo.GetPermissions(userId, boardId)
	if err != nil || boardPermissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	task, err := s.repo.GetById(taksId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"task": task})
	return r
}

func (s *TaskService) Create(userId, projectId, boardId, listId int, task *models.Task) *models.ApiResponse {
	r := &models.ApiResponse{}
	projectPermissions, err := s.projectRepo.GetPermissions(userId, projectId)
	if err != nil || projectPermissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	boardPermissions, err := s.boardRepo.GetPermissions(userId, boardId)
	if err != nil || boardPermissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	task.ListId = listId
	curTime := time.Now().Unix()
	datetimes := &models.Datetimes{
		Created:  curTime,
		Updated:  curTime,
		Accessed: curTime,
	}
	task.Datetimes = datetimes

	taskId, err := s.repo.Create(task)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"taskId": taskId})
	return r
}

func (s *TaskService) Update(userId, projectId, boardId, listId, taskId int, task *models.UpdateTask) *models.ApiResponse {
	r := &models.ApiResponse{}

	if task.Position != nil && *task.Position < 0 {
		r.Error(StatusBadRequest, "Task position out of bounds")
		return r
	}

	if task.ListId != nil && *task.ListId < 1 {
		r.Error(StatusBadRequest, "New list id out of bounds")
		return r
	}

	// TODO: права projectId = read ?
	projectPermissions, err := s.projectRepo.GetPermissions(userId, projectId)
	if err != nil || projectPermissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	boardPermissions, err := s.boardRepo.GetPermissions(userId, boardId)
	if err != nil || boardPermissions.Write == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	curTime := time.Now().Unix()
	task.Datetimes = &models.UpdateDatetimes{
		Updated:  &curTime,
		Accessed: &curTime,
	}

	if err = s.repo.Update(taskId, task); err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{})
	return r
}

func (s *TaskService) Delete(userId, projectId, boardId, listId, taskId int) *models.ApiResponse {
	r := &models.ApiResponse{}
	// TODO: права projectId = read ?
	projectPermissions, err := s.projectRepo.GetPermissions(userId, projectId)
	if err != nil || projectPermissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	boardPermissions, err := s.boardRepo.GetPermissions(userId, boardId)
	if err != nil || boardPermissions.Write == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	err = s.repo.Delete(taskId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{})
	return r
}
