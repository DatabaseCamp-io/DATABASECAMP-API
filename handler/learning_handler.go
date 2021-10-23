package handler

import (
	"DatabaseCamp/controller"
	"DatabaseCamp/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type learningHandler struct {
	controller controller.ILearningController
}

type ILearningHandler interface {
	GetVideo(c *fiber.Ctx) error
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
