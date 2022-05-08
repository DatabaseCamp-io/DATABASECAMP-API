package handler

import (
	"database-camp/internal/infrastructure/application"
	"database-camp/internal/models/request"
	"database-camp/internal/services"
	"database-camp/internal/utils"
	"net/http"
)

type UserHandler interface {
	Register(c application.Context)
	Login(c application.Context)
	GetProfile(c application.Context)
	GetOwnProfile(c application.Context)
	GetUserRanking(c application.Context)
	Edit(c application.Context)
}

type userHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) userHandler {
	return userHandler{service: service}
}

func (h userHandler) Register(c application.Context) {
	request := request.UserRequest{}

	err := c.Bind(&request)
	if err != nil {
		c.Error(err)
		return
	}

	err = request.ValidateRegister()
	if err != nil {
		c.Error(err)
		return
	}

	response, err := h.service.Register(request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h userHandler) Login(c application.Context) {
	request := request.UserRequest{}

	err := c.Bind(&request)
	if err != nil {
		c.Error(err)
		return
	}

	err = request.ValidateLogin()
	if err != nil {
		c.Error(err)
		return
	}

	response, err := h.service.Login(request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h userHandler) GetProfile(c application.Context) {
	id := utils.ParseInt(c.Params("id"))

	response, err := h.service.GetProfile(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h userHandler) GetOwnProfile(c application.Context) {
	id := utils.ParseInt(c.Locals("id"))

	response, err := h.service.GetProfile(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h userHandler) GetUserRanking(c application.Context) {
	id := utils.ParseInt(c.Locals("id"))

	response, err := h.service.GetRanking(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h userHandler) Edit(c application.Context) {
	request := request.UserRequest{}
	userID := utils.ParseInt(c.Locals("id"))

	err := c.Bind(&request)
	if err != nil {
		c.Error(err)
		return
	}

	err = request.ValidateEdit()
	if err != nil {
		c.Error(err)
		return
	}

	response, err := h.service.EditProfile(userID, request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}
