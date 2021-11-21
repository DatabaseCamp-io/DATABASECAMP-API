package handlers

import (
	"DatabaseCamp/controllers"
	"DatabaseCamp/models/request"
	"DatabaseCamp/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type learningHandler struct {
	Controller controllers.ILearningController
}

type ILearningHandler interface {
	GetContentRoadmap(c *fiber.Ctx) error
	GetVideo(c *fiber.Ctx) error
	GetOverview(c *fiber.Ctx) error
	GetActivity(c *fiber.Ctx) error
	UseHint(c *fiber.Ctx) error
	CheckMatchingAnswer(c *fiber.Ctx) error
	CheckMultipleAnswer(c *fiber.Ctx) error
	CheckCompletionAnswer(c *fiber.Ctx) error
}

func NewLearningHandler(controller controllers.ILearningController) learningHandler {
	return learningHandler{Controller: controller}
}

func (h learningHandler) GetContentRoadmap(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	userID := utils.NewType().ParseInt(c.Locals("id"))
	contentID := utils.NewType().ParseInt(c.Params("id"))
	response, err := h.Controller.GetContentRoadmap(userID, contentID)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h learningHandler) GetVideo(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	contentID := c.Params("id")
	response, err := h.Controller.GetVideoLecture(utils.NewType().ParseInt(contentID))
	if err != nil {
		return handleUtil.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h learningHandler) GetOverview(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	id := c.Locals("id")
	response, err := h.Controller.GetOverview(utils.NewType().ParseInt(id))
	if err != nil {
		return handleUtil.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h learningHandler) GetActivity(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	userID := utils.NewType().ParseInt(c.Locals("id"))
	activityID := utils.NewType().ParseInt(c.Params("id"))
	response, err := h.Controller.GetActivity(userID, activityID)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}
	return c.Status(http.StatusOK).JSON(response)
}

func (h learningHandler) UseHint(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	userID := utils.NewType().ParseInt(c.Locals("id"))
	activityID := utils.NewType().ParseInt(c.Params("id"))

	response, err := h.Controller.UseHint(userID, activityID)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func (h learningHandler) CheckMatchingAnswer(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	userID := utils.NewType().ParseInt(c.Locals("id"))
	request := request.MatchingChoiceAnswerRequest{}

	err := handleUtil.BindRequest(c, &request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	err = request.Validate()
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	response, err := h.Controller.CheckAnswer(userID, *request.ActivityID, 1, request.Answer)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func (h learningHandler) CheckMultipleAnswer(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	userID := utils.NewType().ParseInt(c.Locals("id"))
	request := request.MultipleChoiceAnswerRequest{}

	err := handleUtil.BindRequest(c, &request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	err = request.Validate()
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	response, err := h.Controller.CheckAnswer(userID, *request.ActivityID, 2, *request.Answer)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func (h learningHandler) CheckCompletionAnswer(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()
	userID := utils.NewType().ParseInt(c.Locals("id"))
	request := request.CompletionAnswerRequest{}

	err := handleUtil.BindRequest(c, &request)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	err = request.Validate()
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	response, err := h.Controller.CheckAnswer(userID, *request.ActivityID, 3, request.Answer)
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}
