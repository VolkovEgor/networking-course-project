package services

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/repositories"
	"time"
)

type BoardService struct {
	repo        repositories.Board
	projectRepo repositories.Project
}

func NewBoardService(repo repositories.Board, projectRepo repositories.Project) *BoardService {
	return &BoardService{repo: repo, projectRepo: projectRepo}
}

func (s *BoardService) GetAll(userId, projectId int) *models.ApiResponse {
	r := &models.ApiResponse{}
	permissions, err := s.projectRepo.GetPermissions(userId, projectId)
	if err != nil || permissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	boards, err := s.repo.GetAll(userId, projectId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"boards": boards})
	return r
}

func (s *BoardService) Create(userId, projectId int, board *models.Board) *models.ApiResponse {
	r := &models.ApiResponse{}
	permissions, err := s.projectRepo.GetPermissions(userId, projectId)
	if err != nil || permissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	board.OwnerId = userId
	board.ProjectId = projectId
	curTime := time.Now().Unix()
	datetimes := &models.Datetimes{
		Created:  curTime,
		Updated:  curTime,
		Accessed: curTime,
	}
	board.Datetimes = datetimes

	if board.DefaultPermissions == nil {
		board.DefaultPermissions = &models.Permission{
			Read:  true,
			Write: true,
			Admin: false,
		}
	}

	boardId, err := s.repo.Create(userId, board)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"boardId": boardId})
	return r
}

func (s *BoardService) GetById(userId, projectId, boardId int) *models.ApiResponse {
	r := &models.ApiResponse{}
	// TODO: права админов и автора проекта
	projectPermissions, err := s.projectRepo.GetPermissions(userId, projectId)
	if err != nil || projectPermissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	boardPermissions, err := s.repo.GetPermissions(userId, boardId)
	if err != nil || boardPermissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	board, err := s.repo.GetById(boardId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"board": board})
	return r
}

func (s *BoardService) Delete(userId, projectId, boardId int) *models.ApiResponse {
	r := &models.ApiResponse{}

	board, err := s.repo.GetById(boardId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}
	// TODO: права админов и автора проекта
	if board.OwnerId != userId {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	err = s.repo.Delete(boardId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{})
	return r
}

func (s *BoardService) Update(userId, projectId, boardId int, board *models.UpdateBoard) *models.ApiResponse {
	r := &models.ApiResponse{}
	projectPermissions, err := s.projectRepo.GetPermissions(userId, projectId)
	if err != nil || projectPermissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	boardPermissions, err := s.repo.GetPermissions(userId, boardId)
	if err != nil || boardPermissions.Write == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	curTime := time.Now().Unix()
	board.Datetimes = &models.UpdateDatetimes{
		Updated:  &curTime,
		Accessed: &curTime,
	}

	if err = s.repo.Update(boardId, board); err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{})
	return r
}

func (s *BoardService) GetMembers(userId, projectId, boardId int) *models.ApiResponse {
	r := &models.ApiResponse{}

	projectPermissions, err := s.projectRepo.GetPermissions(userId, projectId)
	if err != nil || projectPermissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	boardPermissions, err := s.repo.GetPermissions(userId, boardId)
	if err != nil || boardPermissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	members, err := s.repo.GetMembers(boardId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"members": members})
	return r
}
