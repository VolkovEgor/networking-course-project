package services

import (
	"errors"
	"github.com/architectv/networking-course-project/backend/pkg/builders"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/repositories/postgres"
	"testing"

	mock_repositories "github.com/architectv/networking-course-project/backend/pkg/repositories/mocks"

	"github.com/golang/mock/gomock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Create(t *testing.T) {
	db, err := prepareTestDatabase()
	if err != nil {
		t.Fatalf(err.Error())
	}

	type args struct {
		user *models.User
	}

	tests := []struct {
		name                string
		input               args
		expectedApiResponse *models.ApiResponse
	}{
		{
			name: "Ok",
			input: args{
				user: builders.NewUserBuilder().WithNickname("User Builder").Build(),
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusOK,
				Data: Map{"uid": 10001},
			},
		},
		{
			name: "Already Exists",
			input: args{
				user: builders.NewUserBuilder().WithNickname("User Builder").Build(),
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusConflict,
			},
		},
		{
			name: "Repo Error",
			input: args{
				user: builders.NewUserBuilder().WithNickname("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA").Build(),
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			repo := postgres.NewUserPg(db)
			s := &UserService{repo: repo}

			got := s.Create(test.input.user)
			assert.Equal(t, test.expectedApiResponse.Code, got.Code)
			if test.expectedApiResponse.Code == StatusOK {
				assert.Equal(t, test.expectedApiResponse.Data, got.Data)
			}
		})
	}
}

func TestUserServiceMock_Create(t *testing.T) {
	type args struct {
		user *models.User
	}
	type mockCheck func(r *mock_repositories.MockUser, nickname string)
	type mockCreate func(r *mock_repositories.MockUser, user *models.User)

	tests := []struct {
		name                string
		input               args
		mockCheck           mockCheck
		mockCreate          mockCreate
		expectedApiResponse *models.ApiResponse
	}{
		{
			name: "Ok",
			input: args{
				user: builders.NewUserBuilder().WithNickname("User Builder").Build(),
			},
			mockCheck: func(r *mock_repositories.MockUser, nickname string) {
				r.EXPECT().GetByNickname(nickname).Return(nil, errors.New("new nickname"))
			},
			mockCreate: func(r *mock_repositories.MockUser, user *models.User) {
				r.EXPECT().Create(user).Return(1, nil)
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusOK,
				Data: Map{"uid": 1},
			},
		},
		{
			name: "Already Exists",
			input: args{
				user: builders.NewUserBuilder().WithNickname("User Builder").Build(),
			},
			mockCheck: func(r *mock_repositories.MockUser, nickname string) {
				r.EXPECT().GetByNickname(nickname).Return(nil, nil)
			},
			mockCreate: func(r *mock_repositories.MockUser, user *models.User) {},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusConflict,
			},
		},
		{
			name: "Repo Error",
			input: args{
				user: builders.NewUserBuilder().WithNickname("User Builder").Build(),
			},
			mockCheck: func(r *mock_repositories.MockUser, nickname string) {
				r.EXPECT().GetByNickname(nickname).Return(nil, errors.New("new nickname"))
			},
			mockCreate: func(r *mock_repositories.MockUser, user *models.User) {
				r.EXPECT().Create(user).Return(0, errors.New("repo error"))
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

			repo := mock_repositories.NewMockUser(c)
			test.mockCheck(repo, test.input.user.Nickname)
			test.mockCreate(repo, test.input.user)
			s := &UserService{repo: repo}

			got := s.Create(test.input.user)
			assert.Equal(t, test.expectedApiResponse.Code, got.Code)
			if test.expectedApiResponse.Code == StatusOK {
				assert.Equal(t, test.expectedApiResponse.Data, got.Data)
			}
		})
	}
}

func TestUserService_Get(t *testing.T) {
	db, err := prepareTestDatabase()
	if err != nil {
		t.Fatalf(err.Error())
	}

	type args struct {
		id int
	}

	tests := []struct {
		name                string
		input               args
		expectedApiResponse *models.ApiResponse
	}{
		{
			name: "Ok",
			input: args{
				id: 1,
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusOK,
				Data: Map{"user": &models.User{
					Id:       1,
					Nickname: "test",
					Email:    "test@.mail.ru",
					Password: "qwerty",
					Avatar:   "photo1",
				}},
			},
		},
		{
			name: "Repo err",
			input: args{
				id: 3,
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			repo := postgres.NewUserPg(db)
			s := &UserService{repo: repo}

			got := s.Get(test.input.id)
			assert.Equal(t, test.expectedApiResponse.Code, got.Code)
			if test.expectedApiResponse.Code == StatusOK {
				assert.Equal(t, test.expectedApiResponse.Data, got.Data)
			}
		})
	}
}

func TestUserService_Update(t *testing.T) {
	db, err := prepareTestDatabase()
	if err != nil {
		t.Fatalf(err.Error())
	}

	type args struct {
		id   int
		user *models.UpdateUser
	}

	tests := []struct {
		name                string
		input               args
		expectedApiResponse *models.ApiResponse
	}{
		{
			name: "Ok",
			input: args{
				id:   1,
				user: builders.NewUpdUserBuilder().Build(),
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusOK,
				Data: Map{},
			},
		},
		{
			name: "User is not exists",
			input: args{
				id:   3,
				user: builders.NewUpdUserBuilder().Build(),
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusConflict,
			},
		},
		{
			name: "Repo err",
			input: args{
				id:   3,
				user: builders.NewUpdUserBuilder().WithNickname("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA").Build(),
			},
			expectedApiResponse: &models.ApiResponse{
				Code: StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			repo := postgres.NewUserPg(db)
			s := &UserService{repo: repo}

			got := s.Update(test.input.id, test.input.user)
			assert.Equal(t, test.expectedApiResponse.Code, got.Code)
			if test.expectedApiResponse.Code == StatusOK {
				assert.Equal(t, test.expectedApiResponse.Data, got.Data)
			}
		})
	}
}
