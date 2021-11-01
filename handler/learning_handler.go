package handler

import (
	"DatabaseCamp/controller"
	"DatabaseCamp/models"
	"DatabaseCamp/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type learningHandler struct {
	controller controller.ILearningController
}

type ILearningHandler interface {
	GetVideo(c *fiber.Ctx) error
	GetOverview(c *fiber.Ctx) error
	GetActivity(c *fiber.Ctx) error
	CheckMatchingAnswer(c *fiber.Ctx) error
	UseHint(c *fiber.Ctx) error
}

func NewLearningHandler(controller controller.ILearningController) learningHandler {
	return learningHandler{controller: controller}
}

func (h learningHandler) GetVideo(c *fiber.Ctx) error {
	contentID := c.Params("id")
	response, err := h.controller.GetVideoLecture(utils.NewType().ParseInt(contentID))
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h learningHandler) GetOverview(c *fiber.Ctx) error {
	id := c.Locals("id")
	response, err := h.controller.GetOverview(utils.NewType().ParseInt(id))
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h learningHandler) GetActivity(c *fiber.Ctx) error {
	id := c.Params("id")
	response, err := h.controller.GetActivity(utils.NewType().ParseInt(id))
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h learningHandler) CheckMatchingAnswer(c *fiber.Ctx) error {
	id := c.Locals("id")
	request := models.MatchingChoiceAnswerRequest{}

	err := bindRequest(c, &request)
	if err != nil {
		return handleError(c, err)
	}

	err = request.Validate()
	if err != nil {
		return handleError(c, err)
	}

	response, err := h.controller.CheckMatchingAnswer(utils.NewType().ParseInt(id), request)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func (h learningHandler) UseHint(c *fiber.Ctx) error {
	userID := utils.NewType().ParseInt(c.Locals("id"))
	activityID := utils.NewType().ParseInt(c.Params("id"))

	response, err := h.controller.UseHint(userID, activityID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func (h learningHandler) CheckCompletionAnswer(c *fiber.Ctx) error {
	id := c.Locals("id")
	request := models.CompletionAnswerRequest{}

	err := bindRequest(c, &request)
	if err != nil {
		return handleError(c, err)
	}

	err = request.Validate()
	if err != nil {
		return handleError(c, err)
	}

	response, err := h.controller.CheckCompletionAnswer(utils.NewType().ParseInt(id), request)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}
