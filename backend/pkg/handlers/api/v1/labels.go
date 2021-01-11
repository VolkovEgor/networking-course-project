package v1

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (apiVX *ApiV1) registerLabelsHandlers(router fiber.Router) {
	taskGroup := router.Group("/projects/:pid/boards/:bid/lists/:lid/tasks/:tid/labels", apiVX.userIdentity)
	taskGroup.Get("/", apiVX.urlIdsValidation, apiVX.getLabelsInTask)
	taskGroup.Post("/:tlid", apiVX.urlIdsValidation, apiVX.createLabelInTask)
	taskGroup.Delete("/:tlid", apiVX.urlIdsValidation, apiVX.deleLabelteInTask)

	boardGroup := router.Group("/projects/:pid/boards/:bid/labels", apiVX.userIdentity)
	boardGroup.Get("/", apiVX.urlIdsValidation, apiVX.getLabels)
	boardGroup.Post("/", apiVX.urlIdsValidation, apiVX.createLabel)
	boardGroup.Get("/:tlid", apiVX.urlIdsValidation, apiVX.getLabel)
	boardGroup.Put("/:tlid", apiVX.urlIdsValidation, apiVX.updateLabel)
	boardGroup.Delete("/:tlid", apiVX.urlIdsValidation, apiVX.deleteLabel)
}

func (apiVX *ApiV1) getLabelsInTask(ctx *fiber.Ctx) error {
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
		logrus.Println("boardId:", boardId)
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

	response = apiVX.services.Label.GetAllInTask(userId, projectId, boardId, taskId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) getLabels(ctx *fiber.Ctx) error {
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

	response = apiVX.services.Label.GetAll(userId, projectId, boardId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) getLabel(ctx *fiber.Ctx) error {
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

	labelId, err := strconv.Atoi(ctx.Params("tlid"))
	if err != nil || labelId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid labelId")
		return Send(ctx, response)
	}

	response = apiVX.services.Label.GetById(userId, projectId, boardId, labelId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) createLabel(ctx *fiber.Ctx) error {
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

	input := &models.Label{}
	if err := ctx.BodyParser(input); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	response = apiVX.services.Label.Create(userId, projectId, boardId, input)
	return Send(ctx, response)
}

func (apiVX *ApiV1) createLabelInTask(ctx *fiber.Ctx) error {
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

	labelId, err := strconv.Atoi(ctx.Params("tlid"))
	if err != nil || labelId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid labelId")
		return Send(ctx, response)
	}

	response = apiVX.services.Label.CreateInTask(userId, projectId, boardId, taskId, labelId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) updateLabel(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	userId, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusInternalServerError, err.Error())
		return Send(ctx, response)
	}

	label := &models.UpdateLabel{}
	if err := ctx.BodyParser(label); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	if _, err := govalidator.ValidateStruct(label); err != nil {
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

	labelId, err := strconv.Atoi(ctx.Params("tlid"))
	if err != nil || labelId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid labelId")
		return Send(ctx, response)
	}

	response = apiVX.services.Label.Update(userId, projectId, boardId, labelId, label)
	return Send(ctx, response)
}

func (apiVX *ApiV1) deleLabelteInTask(ctx *fiber.Ctx) error {
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

	labelId, err := strconv.Atoi(ctx.Params("tlid"))
	if err != nil || labelId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid labelId")
		return Send(ctx, response)
	}

	response = apiVX.services.Label.DeleteInTask(userId, projectId, boardId, taskId, labelId)
	return Send(ctx, response)
}

func (apiVX *ApiV1) deleteLabel(ctx *fiber.Ctx) error {
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

	labelId, err := strconv.Atoi(ctx.Params("tlid"))
	if err != nil || labelId == 0 {
		response.Error(fiber.StatusBadRequest, "Invalid labelId")
		return Send(ctx, response)
	}

	response = apiVX.services.Label.Delete(userId, projectId, boardId, labelId)
	return Send(ctx, response)
}
