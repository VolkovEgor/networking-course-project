package services

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/repositories"
)

type UrlValidatorService struct {
	boardRepo repositories.Board
	listRepo  repositories.TaskList
	taskRepo  repositories.Task
}

func NewUrlValidatorService(boardRepo repositories.Board, listRepo repositories.TaskList,
	taskRepo repositories.Task) *UrlValidatorService {
	return &UrlValidatorService{
		boardRepo: boardRepo,
		listRepo:  listRepo,
		taskRepo:  taskRepo,
	}
}

func (s *UrlValidatorService) Validation(urlIds *models.UrlIds) *models.ApiResponse {
	r := &models.ApiResponse{}

	var projectId int
	if urlIds.ProjectId != 0 && urlIds.BoardId != 0 {
		board, err := s.boardRepo.GetById(urlIds.BoardId)
		if err != nil {
			if err.Error() == DbResultNotFound {
				r.Error(StatusNotFound, "Board is not defined")
				return r
			}
			r.Error(StatusInternalServerError, err.Error())
			return r
		} else {
			projectId = board.ProjectId
		}
		if urlIds.ProjectId != projectId {
			r.Error(StatusNotFound, "There is no requested board inside the project")
			return r
		}

		if urlIds.ListId != 0 {
			var boardId int
			list, err := s.listRepo.GetById(urlIds.ListId)
			if err != nil {
				if err.Error() == DbResultNotFound {
					r.Error(StatusNotFound, "List is not defined")
					return r
				}
				r.Error(StatusInternalServerError, err.Error())
				return r
			} else {
				boardId = list.BoardId
			}
			if urlIds.BoardId != boardId {
				r.Error(StatusNotFound, "There is no requested list inside the board")
				return r
			}

			if urlIds.TaskId != 0 {
				var listId int
				task, err := s.taskRepo.GetById(urlIds.TaskId)
				if err != nil {
					if err.Error() == DbResultNotFound {
						r.Error(StatusNotFound, "Task is not defined")
						return r
					}
					r.Error(StatusInternalServerError, err.Error())
					return r
				} else {
					listId = task.ListId
				}
				if urlIds.ListId != listId {
					r.Error(StatusNotFound, "There is no requested task inside the list")
					return r
				}
			}
		}
	}

	r.Set(StatusOK, "OK", Map{})
	return r
}
