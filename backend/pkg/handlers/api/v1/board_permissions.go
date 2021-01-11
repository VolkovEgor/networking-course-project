package v1

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

func (apiVX *ApiV1) registerBoardPermsHandlers(router fiber.Router) {
	group := router.Group("/projects/:pid/boards/:bid/permissions/:member_id", apiVX.userIdentity)
	group.Post("/", apiVX.urlIdsValidation, apiVX.createBoardPerms)
	group.Get("/", apiVX.urlIdsValidation, apiVX.getBoardPerms)
	group.Put("/", apiVX.urlIdsValidation, apiVX.updateBoardPerms)
	group.Delete("/", apiVX.urlIdsValidation, apiVX.deleteBoardPerms)
}

func (apiVX *ApiV1) getBoardPerms(ctx *fiber.Ctx) error {
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
	if err != nil || projectId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid projectId")
		return Send(ctx, response)
	}

	memberId, err := strconv.Atoi(ctx.Params("member_id"))
	if err != nil || memberId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid memberId")
		return Send(ctx, response)
	}

	response = apiVX.services.BoardPerms.Get(userId, projectId, boardId, memberId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) createBoardPerms(ctx *fiber.Ctx) error {
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
	if err != nil || projectId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid projectId")
		return Send(ctx, response)
	}

	memberNickname := ctx.Params("member_id")
	if memberNickname == "" {
		response.Error(fiber.StatusBadRequest, "memberNickname is empty")
		return Send(ctx, response)
	}

	permissions := &models.Permission{}
	if err := ctx.BodyParser(permissions); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	if _, err := govalidator.ValidateStruct(permissions); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	response = apiVX.services.BoardPerms.Create(userId, projectId, boardId, memberNickname,
		permissions)
	return Send(ctx, response)
}

func (apiVX *ApiV1) deleteBoardPerms(ctx *fiber.Ctx) error {
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
	if err != nil || projectId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid projectId")
		return Send(ctx, response)
	}

	memberId, err := strconv.Atoi(ctx.Params("member_id"))
	if err != nil || memberId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid memberId")
		return Send(ctx, response)
	}

	response = apiVX.services.BoardPerms.Delete(userId, projectId, boardId, memberId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) updateBoardPerms(ctx *fiber.Ctx) error {
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
	if err != nil || projectId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid projectId")
		return Send(ctx, response)
	}

	memberId, err := strconv.Atoi(ctx.Params("member_id"))
	if err != nil || memberId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid memberId")
		return Send(ctx, response)
	}

	permissions := &models.UpdatePermission{}
	if err := ctx.BodyParser(permissions); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	if _, err := govalidator.ValidateStruct(permissions); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	response = apiVX.services.BoardPerms.Update(userId, projectId, boardId, memberId,
		permissions)
	return Send(ctx, response)
}
