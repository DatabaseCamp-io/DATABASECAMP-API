package controller

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/repository"
	"DatabaseCamp/utils"
)

type todoController struct {
	repository repository.ITodoRepository
}

func NewTodoController(todoRepository repository.ITodoRepository) todoController {
	return todoController{repository: todoRepository}
}

func (c todoController) GetAll() ([]TodoResponse, error) {
	todoResponses := make([]TodoResponse, 0)

	todo, err := c.repository.GetAll()
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบข้อมูล", "Items Not Found")
	}

	err = utils.NewType().StructToStruct(todo, &todoResponses)
	if err != nil {
		logs.New().Error(err)
		return nil, err
	}
	return todoResponses, err
}
