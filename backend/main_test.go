// +build e2e

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/architectv/networking-course-project/backend/pkg/handlers"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/repositories"
	"github.com/architectv/networking-course-project/backend/pkg/services"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const (
	UsernameTestDB = "postgres"
	PasswordTestDB = "1234"
	HostTestDB     = "localhost"
	PortTestDB     = "5432"
	DBnameTestDB   = "yak_test_db"
	SslmodeTestDB  = "disable"
)

func openTestDatabase() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		HostTestDB, PortTestDB, UsernameTestDB, DBnameTestDB, PasswordTestDB, SslmodeTestDB))
	return db, err
}

func prepareTestDatabase() (*sqlx.DB, error) {
	db, err := openTestDatabase()
	schema, err := ioutil.ReadFile("scripts/init.sql")
	if err != nil {
		log.Fatal(err)
	}
	db.MustExec(string(schema))
	return db, err
}

func Test_E2E_App(t *testing.T) {
	db, err := prepareTestDatabase()
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repositories.NewRepository(db)
	services := services.NewService(repos)
	handlers := handlers.NewHandler(services)

	app := fiber.New()
	app.Use(logger.New())
	handlers.RegisterHandlers(app)

	// Register
	expectedStatus := fiber.StatusOK
	expectedUserId := "5"
	expectedBody := fmt.Sprintf(`{"code":200,"message":"OK","data":{"uid":%s}}`, expectedUserId)
	inputNickname := "e2eUser"
	inputEmail := "e2eUser@test.com"
	inputPassword := "qwerty"
	inputBody := fmt.Sprintf(`{"nickname": "%s", "email": "%s", "password": "%s"}`, inputNickname, inputEmail, inputPassword)

	req := httptest.NewRequest(
		"POST",
		"/api/v1/users/signup",
		bytes.NewBufferString(inputBody),
	)
	req.Header.Set("Content-type", "application/json")

	resp, err := app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, string(body))

	// Login
	expectedStatus = fiber.StatusOK
	inputBody = fmt.Sprintf(`{"nickname": "%s", "password": "%s"}`, inputNickname, inputPassword)

	req = httptest.NewRequest(
		"POST",
		"/api/v1/users/signin",
		bytes.NewBufferString(inputBody),
	)
	req.Header.Set("Content-type", "application/json")

	resp, err = app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	apiResp := &models.ApiResponse{}
	err = json.Unmarshal(body, apiResp)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, apiResp.Code)
	token := apiResp.Data.(map[string]interface{})["token"].(string)

	// Get user info
	expectedStatus = fiber.StatusOK

	req = httptest.NewRequest(
		"GET",
		"/api/v1/users/",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-type", "application/json")

	resp, err = app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	apiResp = &models.ApiResponse{}
	err = json.Unmarshal(body, apiResp)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, apiResp.Code)
	userMap := apiResp.Data.(map[string]interface{})["user"].(map[string]interface{})

	userJSON, err := json.Marshal(userMap)
	assert.Nil(t, err)
	user := &models.User{}
	err = json.Unmarshal(userJSON, user)
	assert.Nil(t, err)

	uid, err := strconv.Atoi(expectedUserId)
	assert.Nil(t, err)
	assert.Equal(t, uid, user.Id)
	assert.Equal(t, inputNickname, user.Nickname)
	assert.Equal(t, inputEmail, user.Email)
	assert.Equal(t, inputPassword, user.Password)

	// Create project
	expectedStatus = fiber.StatusOK
	expectedProjectId := "4"
	expectedBody = fmt.Sprintf(`{"code":200,"message":"OK","data":{"projectId":%s}}`, expectedProjectId)
	inputTitle := "E2E Test Project"
	inputDescription := "LW2 testing"
	inputBody = fmt.Sprintf(`{"title":"%s", "description":"%s"}`, inputTitle, inputDescription)

	req = httptest.NewRequest(
		"POST",
		"/api/v1/projects",
		bytes.NewBufferString(inputBody),
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-type", "application/json")

	resp, err = app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, string(body))

	// Get project
	expectedStatus = fiber.StatusOK

	req = httptest.NewRequest(
		"GET",
		"/api/v1/projects/"+expectedProjectId,
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-type", "application/json")

	resp, err = app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	apiResp = &models.ApiResponse{}
	err = json.Unmarshal(body, apiResp)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, apiResp.Code)
	projectMap := apiResp.Data.(map[string]interface{})["project"].(map[string]interface{})

	projectJSON, err := json.Marshal(projectMap)
	assert.Nil(t, err)
	project := &models.Project{}
	err = json.Unmarshal(projectJSON, project)
	assert.Nil(t, err)

	assert.Equal(t, inputTitle, project.Title)
	assert.Equal(t, inputDescription, project.Description)

	// Create board
	expectedStatus = fiber.StatusOK
	expectedBoardId := "4"
	expectedBody = fmt.Sprintf(`{"code":200,"message":"OK","data":{"boardId":%s}}`, expectedBoardId)
	inputBoardTitle := "E2E Test Board"
	inputBody = fmt.Sprintf(`{"title":"%s"}`, inputBoardTitle)

	req = httptest.NewRequest(
		"POST",
		"/api/v1/projects/"+expectedProjectId+"/boards",
		bytes.NewBufferString(inputBody),
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-type", "application/json")

	resp, err = app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, string(body))

	// Get board
	expectedStatus = fiber.StatusOK

	req = httptest.NewRequest(
		"GET",
		"/api/v1/projects/"+expectedProjectId+"/boards/"+expectedBoardId,
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-type", "application/json")

	resp, err = app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	apiResp = &models.ApiResponse{}
	err = json.Unmarshal(body, apiResp)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, apiResp.Code)
	boardMap := apiResp.Data.(map[string]interface{})["board"].(map[string]interface{})

	boardJSON, err := json.Marshal(boardMap)
	assert.Nil(t, err)
	board := &models.Board{}
	err = json.Unmarshal(boardJSON, board)
	assert.Nil(t, err)

	assert.Equal(t, inputBoardTitle, board.Title)
	assert.Equal(t, uid, board.OwnerId)
	pid, err := strconv.Atoi(expectedProjectId)
	assert.Nil(t, err)
	assert.Equal(t, pid, board.ProjectId)

	// Create board
	expectedStatus = fiber.StatusOK
	expectedListId := "4"
	expectedBody = fmt.Sprintf(`{"code":200,"message":"OK","data":{"listId":%s}}`, expectedListId)
	inputListTitle := "todo list"
	inputBody = fmt.Sprintf(`{"title":"%s"}`, inputListTitle)

	req = httptest.NewRequest(
		"POST",
		"/api/v1/projects/"+expectedProjectId+"/boards/"+expectedBoardId+"/lists",
		bytes.NewBufferString(inputBody),
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-type", "application/json")

	resp, err = app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, string(body))

	// Create task 1
	expectedStatus = fiber.StatusOK
	expectedTaskId1 := "7"
	expectedBody = fmt.Sprintf(`{"code":200,"message":"OK","data":{"taskId":%s}}`, expectedTaskId1)
	inputTaskTitle1 := "Task 1"
	inputTaskDescription1 := "some description"
	inputBody = fmt.Sprintf(`{"title":"%s", "description":"%s"}`, inputTaskTitle1, inputTaskDescription1)

	req = httptest.NewRequest(
		"POST",
		"/api/v1/projects/"+expectedProjectId+"/boards/"+expectedBoardId+"/lists/"+expectedListId+"/tasks",
		bytes.NewBufferString(inputBody),
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-type", "application/json")

	resp, err = app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, string(body))

	// Create task 2
	expectedStatus = fiber.StatusOK
	expectedTaskId2 := "8"
	expectedBody = fmt.Sprintf(`{"code":200,"message":"OK","data":{"taskId":%s}}`, expectedTaskId2)
	inputTaskTitle2 := "Task 2"
	inputTaskDescription2 := "some description"
	inputBody = fmt.Sprintf(`{"title":"%s", "description":"%s"}`, inputTaskTitle2, inputTaskDescription2)

	req = httptest.NewRequest(
		"POST",
		"/api/v1/projects/"+expectedProjectId+"/boards/"+expectedBoardId+"/lists/"+expectedListId+"/tasks",
		bytes.NewBufferString(inputBody),
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-type", "application/json")

	resp, err = app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, string(body))

	// Update task 2
	expectedStatus = fiber.StatusOK
	expectedBody = fmt.Sprintf(`{"code":200,"message":"OK","data":{}}`)
	inputBody = fmt.Sprintf(`{"position":0}`)

	req = httptest.NewRequest(
		"PUT",
		"/api/v1/projects/"+expectedProjectId+"/boards/"+expectedBoardId+"/lists/"+expectedListId+"/tasks/"+expectedTaskId2,
		bytes.NewBufferString(inputBody),
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-type", "application/json")

	resp, err = app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, string(body))

	// Get tasks
	expectedStatus = fiber.StatusOK

	req = httptest.NewRequest(
		"GET",
		"/api/v1/projects/"+expectedProjectId+"/boards/"+expectedBoardId+"/lists/"+expectedListId+"/tasks",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-type", "application/json")

	resp, err = app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	apiResp = &models.ApiResponse{}
	err = json.Unmarshal(body, apiResp)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, apiResp.Code)
	tasksInterface := apiResp.Data.(map[string]interface{})["tasks"].([]interface{})
	var tasksSlice []map[string]interface{}
	for _, item := range tasksInterface {
		tasksSlice = append(tasksSlice, item.(map[string]interface{}))
	}
	assert.Equal(t, 2, len(tasksSlice))
	var tasks []*models.Task
	for _, item := range tasksSlice {
		taskJSON, err := json.Marshal(item)
		assert.Nil(t, err)
		task := &models.Task{}
		err = json.Unmarshal(taskJSON, task)
		assert.Nil(t, err)
		tasks = append(tasks, task)
	}
	lid, err := strconv.Atoi(expectedListId)
	assert.Nil(t, err)
	inputTaskTitle := []string{inputTaskTitle2, inputTaskTitle1}
	inputTaskDescription := []string{inputTaskDescription2, inputTaskDescription1}
	for i, item := range tasks {
		assert.Equal(t, i, item.Position)
		assert.Equal(t, lid, item.ListId)
		assert.Equal(t, inputTaskTitle[i], item.Title)
		assert.Equal(t, inputTaskDescription[i], item.Description)
	}

	// Delete task 2
	expectedStatus = fiber.StatusOK
	expectedBody = fmt.Sprintf(`{"code":200,"message":"OK","data":{}}`)

	req = httptest.NewRequest(
		"DELETE",
		"/api/v1/projects/"+expectedProjectId+"/boards/"+expectedBoardId+"/lists/"+expectedListId+"/tasks/"+expectedTaskId2,
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-type", "application/json")

	resp, err = app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, string(body))

	// Get tasks after delete task 2
	expectedStatus = fiber.StatusOK

	req = httptest.NewRequest(
		"GET",
		"/api/v1/projects/"+expectedProjectId+"/boards/"+expectedBoardId+"/lists/"+expectedListId+"/tasks",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-type", "application/json")

	resp, err = app.Test(req, -1)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	apiResp = &models.ApiResponse{}
	err = json.Unmarshal(body, apiResp)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, apiResp.Code)
	tasksInterface = apiResp.Data.(map[string]interface{})["tasks"].([]interface{})
	tasksSlice = []map[string]interface{}{}
	for _, item := range tasksInterface {
		tasksSlice = append(tasksSlice, item.(map[string]interface{}))
	}
	assert.Equal(t, 1, len(tasksSlice))
	tasks = []*models.Task{}
	for _, item := range tasksSlice {
		taskJSON, err := json.Marshal(item)
		assert.Nil(t, err)
		task := &models.Task{}
		err = json.Unmarshal(taskJSON, task)
		assert.Nil(t, err)
		tasks = append(tasks, task)
	}
	inputTaskTitle = []string{inputTaskTitle1}
	inputTaskDescription = []string{inputTaskDescription1}
	for i, item := range tasks {
		assert.Equal(t, i, item.Position)
		assert.Equal(t, lid, item.ListId)
		assert.Equal(t, inputTaskTitle[i], item.Title)
		assert.Equal(t, inputTaskDescription[i], item.Description)
	}
}
