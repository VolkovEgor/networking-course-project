package v1

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

func (apiVX *ApiV1) registerBoardsHandlers(router fiber.Router) {
	group := router.Group("/projects/:pid/boards", apiVX.userIdentity)
	group.Get("/", apiVX.urlIdsValidation, apiVX.getBoards)
	group.Post("/", apiVX.urlIdsValidation, apiVX.createBoard)
	group.Get("/:bid", apiVX.urlIdsValidation, apiVX.getBoard)
	group.Get("/:bid/members", apiVX.urlIdsValidation, apiVX.getBoardMembers)
	group.Put("/:bid", apiVX.urlIdsValidation, apiVX.updateBoard)
	group.Delete("/:bid", apiVX.urlIdsValidation, apiVX.deleteBoard)
}

func (apiVX *ApiV1) getBoards(ctx *fiber.Ctx) error {
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

	response = apiVX.services.Board.GetAll(userId, projectId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) getBoard(ctx *fiber.Ctx) error {
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

	response = apiVX.services.Board.GetById(userId, projectId, boardId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) createBoard(ctx *fiber.Ctx) error {
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

	input := &models.Board{}
	if err := ctx.BodyParser(input); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	response = apiVX.services.Board.Create(userId, projectId, input)
	return Send(ctx, response)
}

func (apiVX *ApiV1) updateBoard(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	userId, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusInternalServerError, err.Error())
		return Send(ctx, response)
	}

	board := &models.UpdateBoard{}
	if err := ctx.BodyParser(board); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	if _, err := govalidator.ValidateStruct(board); err != nil {
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
		response.Error(fiber.StatusBadRequest, "Invalid boardId")
		return Send(ctx, response)
	}

	response = apiVX.services.Board.Update(userId, projectId, boardId, board)
	return Send(ctx, response)
}

func (apiVX *ApiV1) deleteBoard(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	userId, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusInternalServerError, err.Error())
		return Send(ctx, response)
	}

	projectId, err := strconv.Atoi(ctx.Params("pid"))
	if err != nil {
		response.Error(fiber.StatusInternalServerError, err.Error())
		return Send(ctx, response)
	}

	boardId, err := strconv.Atoi(ctx.Params("bid"))
	if err != nil || boardId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid boardId")
		return Send(ctx, response)
	}

	response = apiVX.services.Board.Delete(userId, projectId, boardId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) getBoardMembers(ctx *fiber.Ctx) error {
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

	response = apiVX.services.Board.GetMembers(userId, projectId, boardId)
	return Send(ctx, response)
}
