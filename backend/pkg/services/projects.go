package services

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/repositories"
	"time"
)

type ProjectService struct {
	repo repositories.Project
}

func NewProjectService(repo repositories.Project) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(userId int, project *models.Project) *models.ApiResponse {
	r := &models.ApiResponse{}

	project.OwnerId = userId
	curTime := time.Now().Unix()
	project.Datetimes = &models.Datetimes{
		Created:  curTime,
		Updated:  curTime,
		Accessed: curTime,
	}

	if project.DefaultPermissions == nil {
		project.DefaultPermissions = &models.Permission{
			Read:  true,
			Write: true,
			Admin: false,
		}
	}

	projectId, err := s.repo.Create(project)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"projectId": projectId})
	return r
}

func (s *ProjectService) GetAll(userId int) *models.ApiResponse {
	r := &models.ApiResponse{}
	projects, err := s.repo.GetAll(userId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}
	r.Set(StatusOK, "OK", Map{"projects": projects})
	return r
}

func (s *ProjectService) GetById(userId, projectId int) *models.ApiResponse {
	r := &models.ApiResponse{}

	permissions, err := s.repo.GetPermissions(userId, projectId)
	if err != nil || permissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	project, err := s.repo.GetById(projectId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}
	r.Set(StatusOK, "OK", Map{"project": project})
	return r
}

func (s *ProjectService) Update(userId, projectId int, project *models.UpdateProject) *models.ApiResponse {
	r := &models.ApiResponse{}
	permissions, err := s.repo.GetPermissions(userId, projectId)

	if err != nil || permissions.Write == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}
	curTime := time.Now().Unix()
	project.Datetimes = &models.UpdateDatetimes{
		Updated:  &curTime,
		Accessed: &curTime,
	}

	if err = s.repo.Update(projectId, project); err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{})
	return r
}

func (s *ProjectService) Delete(userId, projectId int) *models.ApiResponse {
	r := &models.ApiResponse{}

	project, err := s.repo.GetById(projectId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}
	if project.OwnerId != userId {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	err = s.repo.Delete(projectId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}
	r.Set(StatusOK, "OK", Map{})
	return r
}

func (s *ProjectService) GetMembers(userId, projectId int) *models.ApiResponse {
	r := &models.ApiResponse{}

	permissions, err := s.repo.GetPermissions(userId, projectId)
	if err != nil || permissions.Read == false {
		r.Error(StatusForbidden, "Forbidden")
		return r
	}

	members, err := s.repo.GetMembers(projectId)
	if err != nil {
		r.Error(StatusInternalServerError, err.Error())
		return r
	}

	r.Set(StatusOK, "OK", Map{"members": members})
	return r
}
