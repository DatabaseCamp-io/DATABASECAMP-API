package handler

import (
	"DatabaseCamp/controller"
	"DatabaseCamp/models"
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
