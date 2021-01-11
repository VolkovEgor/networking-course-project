package services

import (
	"errors"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"testing"

	mock_repositories "github.com/architectv/networking-course-project/backend/pkg/repositories/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProjectService_Create(t *testing.T) {
	type args struct {
		userId  int
		project *models.Project
	}
	type mockBehavior func(r *mock_repositories.MockProject, project *models.Project)

	tests := []struct {
		name                string
		input               args
		mock                mockBehavior
		expectedApiResponse *models.ApiResponse
	}{
		{
			name: "Ok",
			input: args{
				userId: 1,
				project: &models.Project{
					OwnerId: 1,
					DefaultPermissions: &models.Permission{
						Read:  true,
						Write: true,
						Admin: true,
					},
					Datetimes: &models.Datetimes{
						Created:  1,
						Updated:  1,
						Accessed: 1,
					},
					Title:       "New Test Project",
					Description: "Some Description",
				},
			},
			mock: func(r *mock_repositories.MockProject, project *models.Project) {
				r.EXPECT().Create(project).Return(1, nil)
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusOK,
				Data: Map{"projectId": 1},
			},
		},
		{
			name: "Repo Error",
			input: args{
				userId: 1,
				project: &models.Project{
					OwnerId: 1,
					DefaultPermissions: &models.Permission{
						Read:  true,
						Write: true,
						Admin: true,
					},
					Datetimes: &models.Datetimes{
						Created:  1,
						Updated:  1,
						Accessed: 1,
					},
					Title:       "New Test Project",
					Description: "Some Description",
				},
			},
			mock: func(r *mock_repositories.MockProject, project *models.Project) {
				r.EXPECT().Create(project).Return(0, errors.New("some error"))
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},
		{
			name: "Null Default Permissions",
			input: args{
				userId: 1,
				project: &models.Project{
					OwnerId: 1,
					Datetimes: &models.Datetimes{
						Created:  1,
						Updated:  1,
						Accessed: 1,
					},
					Title:       "New Test Project",
					Description: "Some Description",
				},
			},
			mock: func(r *mock_repositories.MockProject, project *models.Project) {
				r.EXPECT().Create(project).Return(1, nil)
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusOK,
				Data: Map{"projectId": 1},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repositories.NewMockProject(c)
			test.mock(repo, test.input.project)
			s := &ProjectService{repo: repo}

			got := s.Create(test.input.userId, test.input.project)
			assert.Equal(t, test.expectedApiResponse.Code, got.Code)
			if test.expectedApiResponse.Code == StatusOK {
				assert.Equal(t, test.expectedApiResponse.Data, got.Data)
			}
		})
	}
}
