package v1

import (
	"errors"
	"github.com/architectv/networking-course-project/backend/pkg/models"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (apiVX *ApiV1) registerUsersHandlers(router fiber.Router) {
	group := router.Group("/users")
	group.Get("/", apiVX.userIdentity, apiVX.getUser)
	// group.Get("/", apiVX.getUsers)
	group.Post("/signup", apiVX.signUp)
	group.Post("/signin", apiVX.signIn)
	group.Get("/signout", apiVX.userIdentity, apiVX.signOut)
	group.Put("/update", apiVX.userIdentity, apiVX.update)
}

func (apiVX *ApiV1) getUsers(ctx *fiber.Ctx) error {
	users, err := apiVX.services.User.GetAll()
	if err != nil {
		return err
	}
	return ctx.JSON(users)
}

func (apiVX *ApiV1) getUser(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	id, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusUnauthorized, err.Error())
		return Send(ctx, response)
	}

	response = apiVX.services.User.Get(id)
	return Send(ctx, response)
}

func (apiVX *ApiV1) update(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	id, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusUnauthorized, err.Error())
		return Send(ctx, response)
	}

	input := &models.UpdateUser{}
	if err := ctx.BodyParser(&input); err != nil {
		logrus.Println("body parser!")
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}
	if _, err := govalidator.ValidateStruct(input); err != nil {
		logrus.Println("govalid!")
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	response = apiVX.services.User.Update(id, input)
	return Send(ctx, response)
}

func (apiVX *ApiV1) signIn(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	type signInInput struct {
		Nickname string `json:"nickname" valid:"length(3|32)"`
		Password string `json:"password" valid:"length(6|32)"`
	}
	var input signInInput

	if err := ctx.BodyParser(&input); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}
	if _, err := govalidator.ValidateStruct(input); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}

	response = apiVX.services.User.GenerateToken(input.Nickname, input.Password)
	return Send(ctx, response)
}

func (apiVX *ApiV1) signUp(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	input := &models.User{}

	if err := ctx.BodyParser(input); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}
	if _, err := govalidator.ValidateStruct(input); err != nil {
		response.Error(fiber.StatusBadRequest, err.Error())
		return Send(ctx, response)
	}
	response = apiVX.services.User.Create(input)
	return Send(ctx, response)
}

func (apiVX *ApiV1) signOut(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}
	_, err := getUserId(ctx)
	if err != nil {
		response.Error(fiber.StatusUnauthorized, err.Error())
		return Send(ctx, response)
	}

	token, err := getToken(ctx)
	if err != nil {
		response.Error(fiber.StatusUnauthorized, err.Error())
		return Send(ctx, response)
	}

	response = apiVX.services.User.SignOut(token)
	return Send(ctx, response)
}

const (
	authorizationHeader = "Authorization"
	userCtx             = "_id"
)

func getToken(ctx *fiber.Ctx) (string, error) {
	header := ctx.Get(authorizationHeader)
	if header == "" {
		return "", errors.New("Empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("Invalid token")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("Empty token")
	}

	return headerParts[1], nil
}

func (apiVX *ApiV1) userIdentity(ctx *fiber.Ctx) error {
	response := &models.ApiResponse{}

	token, err := getToken(ctx)
	if err != nil {
		response.Error(fiber.StatusUnauthorized, err.Error())
		return Send(ctx, response)
	}

	userId, err := apiVX.services.User.ParseToken(token)
	if err != nil {
		response.Error(fiber.StatusUnauthorized, err.Error())
		return Send(ctx, response)
	}

	ctx.Request().Header.Set(userCtx, strconv.Itoa(userId))
	return ctx.Next()
}

func getUserId(ctx *fiber.Ctx) (int, error) {
	id := ctx.Get(userCtx)
	if id == "" {
		return 0, errors.New("user id not found")
	}

	intId, err := strconv.Atoi(id)
	if err != nil {
		return 0, errors.New("user id is of invalid type")
	}

	return intId, nil
}
