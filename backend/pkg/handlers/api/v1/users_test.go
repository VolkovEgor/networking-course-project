package v1

import (
	"bytes"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/services"
	mock_services "github.com/architectv/networking-course-project/backend/pkg/services/mocks"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUsersHadlers_signUp(t *testing.T) {
	type mockBehavior func(r *mock_services.MockUser, user *models.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            *models.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"nickname": "nickname", "email": "nick@test.com", "password": "qwerty"}`,
			inputUser: &models.User{
				Nickname: "nickname",
				Email:    "nick@test.com",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_services.MockUser, user *models.User) {
				r.EXPECT().Create(user).Return(&models.ApiResponse{
					Code:    200,
					Message: "OK",
					Data: fiber.Map{
						"userId": "1",
					},
				})
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"code":200,"message":"OK","data":{"userId":"1"}}`,
		},
		{
			name:                 "Wrong Email",
			inputBody:            `{"nickname": "nickname", "email": "nick", "password": "qwerty"}`,
			inputUser:            &models.User{},
			mockBehavior:         func(r *mock_services.MockUser, user *models.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":400,"message":"email: nick does not validate as email"}`,
		},
		{
			name:                 "Wrong Nickname (Too Short)",
			inputBody:            `{"nickname": "n", "email": "nick@test.com", "password": "qwerty"}`,
			inputUser:            &models.User{},
			mockBehavior:         func(r *mock_services.MockUser, user *models.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":400,"message":"nickname: n does not validate as length(3|32)"}`,
		},
		{
			name:                 "Wrong Password (Too Short)",
			inputBody:            `{"nickname": "nickname", "email": "nick@test.com", "password": "q"}`,
			inputUser:            &models.User{},
			mockBehavior:         func(r *mock_services.MockUser, user *models.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":400,"message":"password: q does not validate as length(6|32)"}`,
		},
		{
			name:                 "Wrong Request Body",
			inputBody:            ``,
			inputUser:            &models.User{},
			mockBehavior:         func(r *mock_services.MockUser, user *models.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":400,"message":"json: unexpected end of JSON input: "}`,
		},
		{
			name:      "User Already Exists",
			inputBody: `{"nickname": "nickname", "email": "nick@test.com", "password": "qwerty"}`,
			inputUser: &models.User{
				Nickname: "nickname",
				Email:    "nick@test.com",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_services.MockUser, user *models.User) {
				r.EXPECT().Create(user).Return(&models.ApiResponse{
					Code:    409,
					Message: "User already exists",
				})
			},
			expectedStatusCode:   409,
			expectedResponseBody: `{"code":409,"message":"User already exists"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"nickname": "nickname", "email": "nick@test.com", "password": "qwerty"}`,
			inputUser: &models.User{
				Nickname: "nickname",
				Email:    "nick@test.com",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_services.MockUser, user *models.User) {
				r.EXPECT().Create(user).Return(&models.ApiResponse{
					Code:    500,
					Message: "Something went wrong",
				})
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":500,"message":"Something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_services.NewMockUser(c)
			test.mockBehavior(repo, test.inputUser)

			services := &services.Service{User: repo}
			handler := ApiV1{services}

			// Init Endpoint
			r := fiber.New()
			// api := r.Group("/api")
			// apiV1 := NewApiV1(handler.services)
			// apiV1.RegisterHandlers(api)

			url := "/api/auth/users/signup"
			r.Post(url, handler.signUp)

			// Create Request
			req := httptest.NewRequest(
				"POST",
				url,
				bytes.NewBufferString(test.inputBody),
			)
			req.Header.Set("Content-type", "application/json")

			// Make Request
			w, err := r.Test(req, -1)
			assert.Nil(t, err)

			bytesBody, err := ioutil.ReadAll(w.Body)
			assert.Nil(t, err)

			body := string(bytesBody)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.StatusCode)
			assert.Equal(t, test.expectedResponseBody, body)
		})
	}
}
