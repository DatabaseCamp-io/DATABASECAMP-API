package handler

import (
	"DatabaseCamp/controller"
	"DatabaseCamp/models"
	"DatabaseCamp/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	controller controller.IUserController
	jwt        IJwt
}

type IUserHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Edit(c *fiber.Ctx) error
	GetProfile(c *fiber.Ctx) error
	GetInfo(c *fiber.Ctx) error
	GetBadge(c *fiber.Ctx) error
	GetUserRanking(c *fiber.Ctx) error
}

func NewUserHandler(controller controller.IUserController, jwt IJwt) userHandler {
	return userHandler{controller: controller, jwt: jwt}
}

func (h userHandler) Register(c *fiber.Ctx) error {
	request := models.UserRequest{}

	err := bindRequest(c, &request)
	if err != nil {
		return handleError(c, err)
	}

	err = request.ValidateRegister()
	if err != nil {
		return handleError(c, err)
	}

	response, err := h.controller.Register(request)
	if err != nil {
		return handleError(c, err)
	}

	token, err := h.jwt.JwtSign(response.ID)
	if err != nil {
		return handleError(c, err)
	}

	response.AccessToken = token

	return c.Status(http.StatusOK).JSON(response)
}

func (h userHandler) Login(c *fiber.Ctx) error {
	request := models.UserRequest{}

	err := bindRequest(c, &request)
	if err != nil {
		return handleError(c, err)
	}

	err = request.ValidateLogin()
	if err != nil {
		return handleError(c, err)
	}

	response, err := h.controller.Login(request)
	if err != nil {
		return handleError(c, err)
	}

	token, err := h.jwt.JwtSign(response.ID)
	if err != nil {
		return handleError(c, err)
	}

	response.AccessToken = token

	return c.Status(http.StatusOK).JSON(response)
}

func (h userHandler) GetProfile(c *fiber.Ctx) error {
	id := c.Params("id")
	response, err := h.controller.GetProfile(utils.NewType().ParseInt(id))
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h userHandler) GetInfo(c *fiber.Ctx) error {
	id := c.Locals("id")
	response, err := h.controller.GetProfile(utils.NewType().ParseInt(id))
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h userHandler) GetUserRanking(c *fiber.Ctx) error {
	id := c.Locals("id")
	response, err := h.controller.GetRanking(utils.NewType().ParseInt(id))
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h userHandler) Edit(c *fiber.Ctx) error {
	request := models.EditProfileRequest{}
	userID := utils.NewType().ParseInt(c.Locals("id"))

	err := bindRequest(c, &request)
	if err != nil {
		return handleError(c, err)
	}

	err = request.Validate()
	if err != nil {
		return handleError(c, err)
	}

	response, err := h.controller.EditProfile(userID, request)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}
