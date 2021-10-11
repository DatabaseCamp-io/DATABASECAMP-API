package handler

import (
	"DatabaseCamp/controller"
	"DatabaseCamp/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	controller controller.IUserController
	jwt        IJwt
}

type IUserHandler interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
}

func NewUserHandler(controller controller.IUserController, jwt IJwt) userHandler {
	return userHandler{controller: controller, jwt: jwt}
}

func (h userHandler) Register(c echo.Context) error {
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

	return c.JSON(http.StatusOK, response)
}

func (h userHandler) Login(c echo.Context) error {
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

	return c.JSON(http.StatusOK, response)
}
