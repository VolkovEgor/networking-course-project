package v1

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

func (apiVX *ApiV1) registerProjectsHandlers(router fiber.Router) {
	group := router.Group("/projects", apiVX.userIdentity)
	group.Get("/", apiVX.getProjects)
	group.Post("/", apiVX.createProject)
	group.Get("/:pid", apiVX.getProject)
	group.Get("/:pid/members", apiVX.getProjectMembers)
	group.Put("/:pid", apiVX.updateProject)
	group.Delete("/:pid", apiVX.deleteProject)
}

func (apiVX *ApiV1) createProject(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	userId, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusInternalServerError, err.Error())
		return Send(ctx, response)
	}

	project := &models.Project{}

	if err := ctx.BodyParser(project); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	if _, err := govalidator.ValidateStruct(project); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	response = apiVX.services.Project.Create(userId, project)
	return Send(ctx, response)
}

func (apiVX *ApiV1) getProjects(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	userId, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusInternalServerError, err.Error())
		return Send(ctx, response)
	}

	response = apiVX.services.Project.GetAll(userId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) getProject(ctx *fiber.Ctx) error {
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

	response = apiVX.services.Project.GetById(userId, projectId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) updateProject(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	userId, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusInternalServerError, err.Error())
		return Send(ctx, response)
	}

	project := &models.UpdateProject{}
	if err := ctx.BodyParser(project); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	if _, err := govalidator.ValidateStruct(project); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	projectId, err := strconv.Atoi(ctx.Params("pid"))
	if err != nil || projectId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid projectId")
		return Send(ctx, response)
	}

	response = apiVX.services.Project.Update(userId, projectId, project)
	return Send(ctx, response)
}

func (apiVX *ApiV1) deleteProject(ctx *fiber.Ctx) error {
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

	response = apiVX.services.Project.Delete(userId, projectId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) getProjectMembers(ctx *fiber.Ctx) error {
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

	response = apiVX.services.Project.GetMembers(userId, projectId)
	return Send(ctx, response)
}
