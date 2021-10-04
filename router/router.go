package router

import (
	"DatabaseCamp/controller"
	"DatabaseCamp/database"
	"DatabaseCamp/handler"
	"DatabaseCamp/repository"

	"github.com/labstack/echo/v4"
)

type router struct {
	echo *echo.Echo
}

type IRouter interface {
}

var instantiated *router = nil

func New(e *echo.Echo) *router {
	if instantiated == nil {
		instantiated = &router{echo: e}
		instantiated.init()
	}
	return instantiated
}

func (r router) init() {
	r.setupTodo()
}

func (r router) setupTodo() {
	db := database.New().Get()
	repository := repository.NewTodoRepository(db)
	controller := controller.NewTodoController(repository)
	handler := handler.NewTodoHandler(controller)
	group := r.echo.Group("todo")
	{
		group.GET("/list", handler.GetAll)
	}
}
