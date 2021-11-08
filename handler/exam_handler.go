package handler

import (
	"DatabaseCamp/controller"
	"DatabaseCamp/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type examHandler struct {
	controller controller.IExamController
}

type IExamHandler interface {
	GetExam(c *fiber.Ctx) error
	CheckExam(c *fiber.Ctx) error
	GetExamOverview(c *fiber.Ctx) error
}

func NewExamHandler(controller controller.IExamController) examHandler {
	return examHandler{controller: controller}
}

func (h examHandler) GetExam(c *fiber.Ctx) error {
	examID := utils.NewType().ParseInt(c.Params("id"))
	userID := utils.NewType().ParseInt(c.Locals("id"))
	response, err := h.controller.GetExam(examID, userID)
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h examHandler) CheckExam(c *fiber.Ctx) error {
	return nil
}

func (h examHandler) GetExamOverview(c *fiber.Ctx) error {
	return nil
}
