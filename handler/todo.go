package handler

import (
	"DatabaseCamp/controller"
	"net/http"

	"github.com/labstack/echo/v4"
)

type todoHandler struct {
	controller controller.ITodoController
}

type ITodoHandler interface {
	GetAll(c echo.Context) error
}

func NewTodoHandler(controller controller.ITodoController) todoHandler {
	return todoHandler{controller: controller}
}

func (h todoHandler) GetAll(c echo.Context) error {
	todo, err := h.controller.GetAll()
	if err != nil {
		return HandleError(c, err)
	}
	return c.JSON(http.StatusOK, todo)
}
