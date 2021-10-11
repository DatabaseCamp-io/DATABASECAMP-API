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

var instantiated *router = nil

func New(e *echo.Echo) *router {
	if instantiated == nil {
		instantiated = &router{echo: e}
		instantiated.init()
	}
	return instantiated
}

func (r router) init() {
	r.setupUser()
}

func (r router) setupUser() {
	db := database.New()
	repository := repository.NewUserRepository(db)
	controller := controller.NewUserController(repository)
	jwt := handler.NewJwtMiddleware(repository)
	handler := handler.NewUserHandler(controller, jwt)
	group := r.echo.Group("user")
	{
		group.POST("/register/", handler.Register)
		group.POST("/login/", handler.Login)
	}
}
