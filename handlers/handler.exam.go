package handlers

// handler.exam.go
/**
 * 	This file is a part of handler, used to handle request of the exam
 */

import (
	"DatabaseCamp/controllers"
	"DatabaseCamp/models/request"
	"DatabaseCamp/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

/**
 * This class handle request of the exam
 */
type examHandler struct {
	Controller controllers.IExamController // Exam controller for doing business logic of the exam
}

/**
 * Constructor creates a new examHandler instance
 *
 * @param   controller    	Exam controller for doing business logic of the exam
 *
 * @return 	instance of examHandler
 */
func NewExamHandler(controller controllers.IExamController) examHandler {
	return examHandler{Controller: controller}
}

/**
 * Get the exam to use for the test
 *
 * @param 	c  context of the web framework
 *
 * @return the error of getting exam
 */
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

/**
 * Check answer of the exam
 *
 * @param 	c  context of the web framework
 *
 * @return the error of getting exam
 */
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

/**
 * Get overview of the exam
 *
 * @param 	c  context of the web framework
 *
 * @return the error of getting exam
 */
func (h examHandler) GetExamOverview(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()

	userID := utils.NewType().ParseInt(c.Locals("id"))

	response, err := h.Controller.GetOverview(userID)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}

/**
 * Get exam result of the user
 *
 * @param 	c  context of the web framework
 *
 * @return the error of getting exam
 */
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
