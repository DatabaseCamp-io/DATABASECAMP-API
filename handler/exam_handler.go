package handler

import "DatabaseCamp/controller"

type examHandler struct {
	controller controller.IExamController
}

type IExamHandler interface {
}

func NewExamHandler(controller controller.IExamController) examHandler {
	return examHandler{controller: controller}
}
