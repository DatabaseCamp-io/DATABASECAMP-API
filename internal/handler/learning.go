package handler

import (
	"database-camp/internal/infrastructure/application"
	"database-camp/internal/models/request"
	"database-camp/internal/services"
	"database-camp/internal/utils"
	"net/http"
)

type LearningHandler interface {
	GetContentRoadmap(c application.Context)
	GetVideo(c application.Context)
	GetOverview(c application.Context)
	GetActivity(c application.Context)
	GetRecommend(c application.Context)
	UseHint(c application.Context)
	CheckAnswer(c application.Context)
}

type learningHandler struct {
	service services.LearningService
}

func NewLearningHandler(service services.LearningService) *learningHandler {
	return &learningHandler{service: service}
}

func (h learningHandler) GetContentRoadmap(c application.Context) {
	userID := utils.ParseInt(c.Locals("id"))
	contentID := utils.ParseInt(c.Params("id"))

	response, err := h.service.GetContentRoadmap(userID, contentID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h learningHandler) GetVideo(c application.Context) {
	contentID := c.Params("id")

	response, err := h.service.GetVideoLecture(utils.ParseInt(contentID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h learningHandler) GetOverview(c application.Context) {
	id := c.Locals("id")

	response, err := h.service.GetOverview(utils.ParseInt(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h learningHandler) GetActivity(c application.Context) {
	userID := utils.ParseInt(c.Locals("id"))
	activityID := utils.ParseInt(c.Params("id"))

	response, err := h.service.GetActivity(userID, activityID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h learningHandler) GetRecommend(c application.Context) {
	userID := utils.ParseInt(c.Locals("id"))

	response, err := h.service.GetRecommend(userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h learningHandler) UseHint(c application.Context) {
	userID := utils.ParseInt(c.Locals("id"))
	activityID := utils.ParseInt(c.Params("id"))

	response, err := h.service.UseHint(userID, activityID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h learningHandler) CheckAnswer(c application.Context) {
	userID := utils.ParseInt(c.Locals("id"))
	request := request.CheckAnswerRequest{}

	err := c.Bind(&request)
	if err != nil {
		c.Error(err)
		return
	}

	err = request.Validate()
	if err != nil {
		c.Error(err)
		return
	}

	response, err := h.service.CheckAnswer(userID, request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}
