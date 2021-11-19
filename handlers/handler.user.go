package handlers

import (
	"DatabaseCamp/controllers"
	"DatabaseCamp/middleware"
	"DatabaseCamp/models"
	"DatabaseCamp/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	controller controllers.IUserController
	jwt        middleware.IJwt
}

func NewUserHandler(controller controllers.IUserController, jwt middleware.IJwt) userHandler {
	return userHandler{controller: controller, jwt: jwt}
}

func (h userHandler) Register(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	request := models.UserRequest{}

	err := handleUtil.BindRequest(c, &request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	err = request.ValidateRegister()
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	response, err := h.controller.Register(request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	token, err := h.jwt.JwtSign(response.ID)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	response.AccessToken = token

	return c.Status(http.StatusOK).JSON(response)
}

func (h userHandler) Login(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	request := models.UserRequest{}

	err := handleUtil.BindRequest(c, &request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	err = request.ValidateLogin()
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	response, err := h.controller.Login(request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	token, err := h.jwt.JwtSign(response.ID)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	response.AccessToken = token

	return c.Status(http.StatusOK).JSON(response)
}

func (h userHandler) GetProfile(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	id := c.Params("id")
	response, err := h.controller.GetProfile(utils.NewType().ParseInt(id))
	if err != nil {
		return handleUtil.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h userHandler) GetOwnProfile(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	id := c.Locals("id")
	response, err := h.controller.GetProfile(utils.NewType().ParseInt(id))
	if err != nil {
		return handleUtil.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h userHandler) GetUserRanking(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	id := c.Locals("id")
	response, err := h.controller.GetRanking(utils.NewType().ParseInt(id))
	if err != nil {
		return handleUtil.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h userHandler) Edit(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	request := models.UserRequest{}
	userID := utils.NewType().ParseInt(c.Locals("id"))

	err := handleUtil.BindRequest(c, &request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	err = request.ValidateEdit()
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	response, err := h.controller.EditProfile(userID, request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}
