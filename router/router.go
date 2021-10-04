package router

import (
	"DatabaseCamp/controller"
	"DatabaseCamp/database"
	"DatabaseCamp/handler"
	"DatabaseCamp/repository"

	"github.com/labstack/echo/v4"
)

type Router struct {
	router *echo.Echo
}

type IRouter interface {
}

var instantiated *Router = nil

func New(router *echo.Echo) *Router {
	if instantiated == nil {
		instantiated = &Router{router: router}
		instantiated.init()
	}
	return instantiated
}

func (r Router) init() {
	r.setupTodo()
}

func (r Router) setupTodo() {
	db := database.New().Get()
	repository := repository.NewTodoRepository(db)
	controller := controller.NewTodoController(repository)
	handler := handler.NewTodoHandler(controller)
	group := r.router.Group("todo")
	{
		group.GET("/list", handler.GetAll)
	}
}
