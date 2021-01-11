package api

import (
	v1 "github.com/architectv/networking-course-project/backend/pkg/handlers/api/v1"
	"github.com/architectv/networking-course-project/backend/pkg/services"

	"github.com/gofiber/fiber/v2"
)

type Api struct {
	services *services.Service
}

func NewApi(services *services.Service) *Api {
	return &Api{services: services}
}

func (a *Api) RegisterHandlers(router fiber.Router) {
	api := router.Group("/api")
	apiV1 := v1.NewApiV1(a.services)
	apiV1.RegisterHandlers(api)
}
