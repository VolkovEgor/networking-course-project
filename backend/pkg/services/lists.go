package services

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/repositories"
)

type TaskListService struct {
	repo        repositories.TaskList
	boardRepo   repositories.Board
	projectRepo repositories.Project
}

func NewTaskListService(repo repositories.TaskList, boardRepo repositories.Board, projectRepo repositories.Project) *TaskListService {
	return &TaskListService{repo: repo, boardRepo: boardRepo, projectRepo: projectRepo}
}

func (s *TaskListService) GetAll(userId, projectId, boardId int) *models.ApiResponse {
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

	lists, err := s.repo.GetAll(boardId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"lists": lists})
	return r
}

func (s *TaskListService) GetById(userId, projectId, boardId, listId int) *models.ApiResponse {
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

	list, err := s.repo.GetById(listId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"list": list})
	return r
}

func (s *TaskListService) Create(userId, projectId, boardId int, list *models.TaskList) *models.ApiResponse {
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

	list.BoardId = boardId
	listId, err := s.repo.Create(list)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"listId": listId})
	return r
}

func (s *TaskListService) Delete(userId, projectId, boardId, listId int) *models.ApiResponse {
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

	err = s.repo.Delete(listId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{})
	return r
}

func (s *TaskListService) Update(userId, projectId, boardId, listId int, list *models.UpdateTaskList) *models.ApiResponse {
	r := &models.ApiResponse{}

	if list.Position != nil && *list.Position < 0 {
		r.Error(StatusBadRequest, "List position out of bounds")
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

	if err = s.repo.Update(listId, list); err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{})
	return r
}
