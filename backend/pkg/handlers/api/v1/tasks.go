package v1

import (
	"fmt"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

func (apiVX *ApiV1) registerTasksHandlers(router fiber.Router) {
	group := router.Group("/projects/:pid/boards/:bid/lists/:lid/tasks", apiVX.userIdentity)
	group.Get("/", apiVX.urlIdsValidation, apiVX.getTasks)
	group.Post("/", apiVX.urlIdsValidation, apiVX.createTask)
	group.Get("/:tid", apiVX.urlIdsValidation, apiVX.getTask)
	group.Put("/:tid", apiVX.urlIdsValidation, apiVX.updateTask)
	group.Delete("/:tid", apiVX.urlIdsValidation, apiVX.deleteTask)
}

func (apiVX *ApiV1) getTasks(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	userId, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusInternalServerError, err.Error())
		return Send(ctx, response)
	}

	projectId, err := strconv.Atoi(ctx.Params("pid"))
	if err != nil || projectId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid projectId")
		return Send(ctx, response)
	}

	boardId, err := strconv.Atoi(ctx.Params("bid"))
	if err != nil || boardId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid boardId")
		return Send(ctx, response)
	}

	listId, err := strconv.Atoi(ctx.Params("lid"))
	if err != nil || listId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid listId")
		return Send(ctx, response)
	}

	response = apiVX.services.Task.GetAll(userId, projectId, boardId, listId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) getTask(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	userId, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusInternalServerError, err.Error())
		return Send(ctx, response)
	}

	projectId, err := strconv.Atoi(ctx.Params("pid"))
	if err != nil || projectId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid projectId")
		return Send(ctx, response)
	}

	boardId, err := strconv.Atoi(ctx.Params("bid"))
	if err != nil || boardId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid boardId")
		return Send(ctx, response)
	}

	listId, err := strconv.Atoi(ctx.Params("lid"))
	if err != nil || listId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid listId")
		return Send(ctx, response)
	}

	taskId, err := strconv.Atoi(ctx.Params("tid"))
	if err != nil || taskId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid taskId")
		return Send(ctx, response)
	}

	response = apiVX.services.Task.GetById(userId, projectId, boardId,
		listId, taskId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) createTask(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	userId, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusInternalServerError, err.Error())
		return Send(ctx, response)
	}

	projectId, err := strconv.Atoi(ctx.Params("pid"))
	if err != nil || projectId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid projectId")
		return Send(ctx, response)
	}

	boardId, err := strconv.Atoi(ctx.Params("bid"))
	if err != nil || boardId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid boardId")
		return Send(ctx, response)
	}

	listId, err := strconv.Atoi(ctx.Params("lid"))
	if err != nil || listId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid listId")
		return Send(ctx, response)
	}

	input := &models.Task{}
	if err := ctx.BodyParser(input); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	response = apiVX.services.Task.Create(userId, projectId, boardId, listId, input)
	return Send(ctx, response)
}

func (apiVX *ApiV1) updateTask(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	userId, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusInternalServerError, err.Error())
		return Send(ctx, response)
	}

	task := &models.UpdateTask{}
	if err := ctx.BodyParser(task); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	if _, err := govalidator.ValidateStruct(task); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	projectId, err := strconv.Atoi(ctx.Params("pid"))
	if err != nil {
		response.Error(fiber.StatusInternalServerError, err.Error())
		return Send(ctx, response)
	}

	boardId, err := strconv.Atoi(ctx.Params("bid"))
	if err != nil || boardId == 0 {
		response.Error(fiber.StatusBadRequest, "Empty boardId")
		return Send(ctx, response)
	}

	listId, err := strconv.Atoi(ctx.Params("lid"))
	if err != nil || listId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid listId")
		return Send(ctx, response)
	}

	taskId, err := strconv.Atoi(ctx.Params("tid"))
	if err != nil || taskId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid taskId")
		return Send(ctx, response)
	}
	fmt.Println("----")
	response = apiVX.services.Task.Update(userId, projectId, boardId, listId, taskId, task)
	return Send(ctx, response)
}

func (apiVX *ApiV1) deleteTask(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	userId, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusInternalServerError, err.Error())
		return Send(ctx, response)
	}

	projectId, err := strconv.Atoi(ctx.Params("pid"))
	if err != nil || projectId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid projectId")
		return Send(ctx, response)
	}

	boardId, err := strconv.Atoi(ctx.Params("bid"))
	if err != nil || boardId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid boardId")
		return Send(ctx, response)
	}

	listId, err := strconv.Atoi(ctx.Params("lid"))
	if err != nil || listId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid listId")
		return Send(ctx, response)
	}

	taskId, err := strconv.Atoi(ctx.Params("tid"))
	if err != nil || taskId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid taskId")
		return Send(ctx, response)
	}

	response = apiVX.services.Task.Delete(userId, projectId, boardId, listId, taskId)
	return Send(ctx, response)
}
