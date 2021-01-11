package v1

import (
	"fmt"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"github.com/architectv/networking-course-project/backend/pkg/services"
	"runtime"

	"github.com/gofiber/fiber/v2"
)

func implementMe() {
	pc, fn, line, _ := runtime.Caller(1)
	fmt.Printf("Implement me in %s[%s:%d]\n", runtime.FuncForPC(pc).Name(), fn, line)
}

type ApiV1 struct {
	services *services.Service
}

func NewApiV1(services *services.Service) *ApiV1 {
	return &ApiV1{services: services}
}

func (apiVX *ApiV1) RegisterHandlers(router fiber.Router) {
	v1 := router.Group("/v1")
	apiVX.registerBoardPermsHandlers(v1)
	apiVX.registerBoardsHandlers(v1)
	apiVX.registerListsHandlers(v1)
	apiVX.registerProjectPermsHandlers(v1)
	apiVX.registerProjectsHandlers(v1)
	apiVX.registerTasksHandlers(v1)
	apiVX.registerUsersHandlers(v1)
	apiVX.registerLabelsHandlers(v1)
}

func Send(ctx *fiber.Ctx, r *models.ApiResponse) error {
	ctx.Status(r.Code)
	return ctx.JSON(r)
}
