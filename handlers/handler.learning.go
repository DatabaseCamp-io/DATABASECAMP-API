package handlers

// handler.learning.go
/**
 * 	This file is a part of handler, used to handle request of the learning
 */

import (
	"DatabaseCamp/controllers"
	"DatabaseCamp/models/request"
	"DatabaseCamp/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

/**
 * This class handle request of the learning
 */
type learningHandler struct {
	Controller controllers.ILearningController
}

/**
 * Constructor creates a new learningHandler instance
 *
 * @param   controller    	Learning controller for doing business logic of the learning
 *
 * @return 	instance of learningHandler
 */
func NewLearningHandler(controller controllers.ILearningController) learningHandler {
	return learningHandler{Controller: controller}
}

/**
 * Get content roadmap
 *
 * @param 	c  context of the web framework
 *
 * @return the error of getting exam
 */
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

/**
 * Get video lecture of the content
 *
 * @param 	c  context of the web framework
 *
 * @return the error of getting exam
 */
func (h learningHandler) GetVideo(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()

	contentID := c.Params("id")

	response, err := h.Controller.GetVideoLecture(utils.NewType().ParseInt(contentID))
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}

/**
 * Get content overview
 *
 * @param 	c  context of the web framework
 *
 * @return the error of getting exam
 */
func (h learningHandler) GetOverview(c *fiber.Ctx) error {
	handleUtil := utils.NewHandle()

	id := c.Locals("id")

	response, err := h.Controller.GetOverview(utils.NewType().ParseInt(id))
	if err != nil {
		return handleUtil.HandleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(response)
}

/**
 * Get activity for user to do
 *
 * @param 	c  context of the web framework
 *
 * @return the error of getting exam
 */
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

/**
 * Use hint of the activity
 *
 * @param 	c  context of the web framework
 *
 * @return the error of getting exam
 */
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

/**
 * Check matching choice answer of the activity
 *
 * @param 	c  context of the web framework
 *
 * @return the error of getting exam
 */
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

/**
 * Check multiple choice answer of the activity
 *
 * @param 	c  context of the web framework
 *
 * @return the error of getting exam
 */
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

/**
 * Check completion choice answer of the activity
 *
 * @param 	c  context of the web framework
 *
 * @return the error of getting exam
 */
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
