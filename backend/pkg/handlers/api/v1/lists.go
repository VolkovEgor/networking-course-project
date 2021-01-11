package v1

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

func (apiVX *ApiV1) registerListsHandlers(router fiber.Router) {
	group := router.Group("/projects/:pid/boards/:bid/lists", apiVX.userIdentity)
	group.Get("/", apiVX.urlIdsValidation, apiVX.getLists)
	group.Post("/", apiVX.urlIdsValidation, apiVX.createList)
	group.Get("/:lid", apiVX.urlIdsValidation, apiVX.getList)
	group.Patch("/:lid", apiVX.urlIdsValidation, apiVX.updateList)
	group.Delete("/:lid", apiVX.urlIdsValidation, apiVX.deleteList)
}

func (apiVX *ApiV1) getLists(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	userId, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusInternalServerError, err.Error())
		return Send(ctx, response)
	}

	projectId, err := strconv.Atoi(ctx.Params("pid"))
	if err != nil || projectId == 0 {
		response.Error(fiber.StatusBadRequest, "Empty projectId")
		return Send(ctx, response)
	}

	boardId, err := strconv.Atoi(ctx.Params("bid"))
	if err != nil || boardId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid boardId")
		return Send(ctx, response)
	}

	response = apiVX.services.TaskList.GetAll(userId, projectId, boardId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) getList(ctx *fiber.Ctx) error {
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

	response = apiVX.services.TaskList.GetById(userId, projectId, boardId, listId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) createList(ctx *fiber.Ctx) error {
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

	input := &models.TaskList{}
	if err := ctx.BodyParser(input); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	response = apiVX.services.TaskList.Create(userId, projectId, boardId, input)
	return Send(ctx, response)
}

func (apiVX *ApiV1) updateList(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	userId, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusInternalServerError, err.Error())
		return Send(ctx, response)
	}

	list := &models.UpdateTaskList{}
	if err := ctx.BodyParser(list); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	if _, err := govalidator.ValidateStruct(list); err != nil {
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

	response = apiVX.services.TaskList.Update(userId, projectId, boardId, listId, list)
	return Send(ctx, response)
}

func (apiVX *ApiV1) deleteList(ctx *fiber.Ctx) error {
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

	response = apiVX.services.TaskList.Delete(userId, projectId, boardId, listId)
	return Send(ctx, response)
}
