package handlers

import (
	"DatabaseCamp/controllers"
	"DatabaseCamp/models/request"
	"DatabaseCamp/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type examHandler struct {
	Controller controllers.IExamController
}

type IExamHandler interface {
	GetExam(c *fiber.Ctx) error
	CheckExam(c *fiber.Ctx) error
	GetExamOverview(c *fiber.Ctx) error
	GetExamResult(c *fiber.Ctx) error
}

func NewExamHandler(controller controllers.IExamController) examHandler {
	return examHandler{Controller: controller}
}

func (h examHandler) GetExam(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	examID := utils.NewType().ParseInt(c.Params("id"))
	userID := utils.NewType().ParseInt(c.Locals("id"))
	response, err := h.Controller.GetExam(examID, userID)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h examHandler) CheckExam(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	userID := utils.NewType().ParseInt(c.Locals("id"))
	request := request.ExamAnswerRequest{}

	err := handleUtil.BindRequest(c, &request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	err = request.Validate()
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	response, err := h.Controller.CheckExam(userID, request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func (h examHandler) GetExamOverview(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	userID := utils.NewType().ParseInt(c.Locals("id"))
	response, err := h.Controller.GetOverview(userID)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h examHandler) GetExamResult(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	userID := utils.NewType().ParseInt(c.Locals("id"))
	examResultID := utils.NewType().ParseInt(c.Params("id"))
	response, err := h.Controller.GetExamResult(userID, examResultID)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}
