package v1

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (apiVX *ApiV1) urlIdsValidation(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}

	projectId, _ := strconv.Atoi(ctx.Params("pid"))
	boardId, _ := strconv.Atoi(ctx.Params("bid"))
	listId, _ := strconv.Atoi(ctx.Params("lid"))
	taskId, _ := strconv.Atoi(ctx.Params("tid"))

	urlIds := &models.UrlIds{
		ProjectId: projectId,
		BoardId:   boardId,
		ListId:    listId,
		TaskId:    taskId,
	}

	response = apiVX.services.UrlValidator.Validation(urlIds)
	if response.Code != fiber.StatusOK {
		return Send(ctx, response)
	}
	return ctx.Next()

}
