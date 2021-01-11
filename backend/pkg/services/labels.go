package services

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/repositories"
)

type LabelService struct {
	repo        repositories.Label
	boardRepo   repositories.Board
	projectRepo repositories.Project
}

func NewLabelService(repo repositories.Label, boardRepo repositories.Board, projectRepo repositories.Project) *LabelService {
	return &LabelService{repo: repo, boardRepo: boardRepo, projectRepo: projectRepo}
}

func (s *LabelService) GetAllInTask(userId, projectId, boardId, taskId int) *models.ApiResponse {
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

	labels, err := s.repo.GetAllInTask(taskId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"labels": labels})
	return r
}

func (s *LabelService) GetAll(userId, projectId, boardId int) *models.ApiResponse {
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

	labels, err := s.repo.GetAll(boardId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"labels": labels})
	return r
}

func (s *LabelService) GetById(userId, projectId, boardId, labelId int) *models.ApiResponse {
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

	label, err := s.repo.GetById(labelId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"label": label})
	return r
}

func (s *LabelService) Create(userId, projectId, boardId int, label *models.Label) *models.ApiResponse {
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

	label.BoardId = boardId
	labelId, err := s.repo.Create(label)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"labelId": labelId})
	return r
}

func (s *LabelService) CreateInTask(userId, projectId, boardId, taskId, labelId int) *models.ApiResponse {
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

	taskLabelId, err := s.repo.CreateInTask(taskId, labelId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"taskLabelId": taskLabelId})
	return r
}

func (s *LabelService) Update(userId, projectId, boardId, labelId int, label *models.UpdateLabel) *models.ApiResponse {
	r := &models.ApiResponse{}
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

	if err = s.repo.Update(labelId, label); err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{})
	return r
}

func (s *LabelService) DeleteInTask(userId, projectId, boardId, taskId, labelId int) *models.ApiResponse {
	r := &models.ApiResponse{}
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

	err = s.repo.DeleteInTask(taskId, labelId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{})
	return r
}

func (s *LabelService) Delete(userId, projectId, boardId, labelId int) *models.ApiResponse {
	r := &models.ApiResponse{}
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

	err = s.repo.Delete(labelId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{})
	return r
}
