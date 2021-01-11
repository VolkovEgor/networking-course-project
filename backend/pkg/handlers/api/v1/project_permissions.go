package v1

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

func (apiVX *ApiV1) registerProjectPermsHandlers(router fiber.Router) {
	group := router.Group("/projects/:pid/permissions/:member_id", apiVX.userIdentity)
	group.Post("/", apiVX.createProjectPerms)
	group.Get("/", apiVX.getProjectPerms)
	group.Put("/", apiVX.updateProjectPerms)
	group.Delete("/", apiVX.deleteProjectPerms)
}

func (apiVX *ApiV1) getProjectPerms(ctx *fiber.Ctx) error {
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

	memberId, err := strconv.Atoi(ctx.Params("member_id"))
	if err != nil || memberId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid memberId")
		return Send(ctx, response)
	}

	response = apiVX.services.ProjectPerms.Get(userId, projectId, memberId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) createProjectPerms(ctx *fiber.Ctx) error {
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

	response = apiVX.services.ProjectPerms.Create(userId, projectId, memberNickname,
		permissions)
	return Send(ctx, response)
}

func (apiVX *ApiV1) deleteProjectPerms(ctx *fiber.Ctx) error {
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

	memberId, err := strconv.Atoi(ctx.Params("member_id"))
	if err != nil || memberId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid memberId")
		return Send(ctx, response)
	}

	response = apiVX.services.ProjectPerms.Delete(userId, projectId, memberId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) updateProjectPerms(ctx *fiber.Ctx) error {
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

	response = apiVX.services.ProjectPerms.Update(userId, projectId, memberId,
		permissions)
	return Send(ctx, response)
}
