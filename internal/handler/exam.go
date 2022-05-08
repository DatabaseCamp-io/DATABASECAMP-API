package handler

import (
	"database-camp/internal/infrastructure/application"
	"database-camp/internal/models/request"
	"database-camp/internal/services"
	"database-camp/internal/utils"
	"net/http"
)

type ExamHandler interface {
	GetExam(c application.Context)
	GetExamOverview(c application.Context)
	GetExamResult(c application.Context)
	CheckExam(c application.Context)
}

type examHandler struct {
	service services.ExamService
}

func NewExamHandler(service services.ExamService) *examHandler {
	return &examHandler{service: service}
}

func (h examHandler) GetExam(c application.Context) {
	examID := utils.ParseInt(c.Params("id"))
	userID := utils.ParseInt(c.Locals("id"))

	response, err := h.service.GetExam(examID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h examHandler) CheckExam(c application.Context) {
	userID := utils.ParseInt(c.Locals("id"))
	request := request.ExamAnswerRequest{}

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

	response, err := h.service.CheckExam(userID, request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h examHandler) GetExamOverview(c application.Context) {
	userID := utils.ParseInt(c.Locals("id"))

	response, err := h.service.GetOverview(userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h examHandler) GetExamResult(c application.Context) {
	userID := utils.ParseInt(c.Locals("id"))
	examResultID := utils.ParseInt(c.Params("id"))

	response, err := h.service.GetExamResult(userID, examResultID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}
