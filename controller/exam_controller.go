package controller

import "DatabaseCamp/repository"

type examController struct {
	examRepo repository.IExamRepository
}

type IExamController interface {
}

func NewExamController(examRepo repository.IExamRepository) examController {
	return examController{examRepo: examRepo}
}
