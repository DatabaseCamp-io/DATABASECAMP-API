package router

import (
	"DatabaseCamp/controller"
	"DatabaseCamp/database"
	"DatabaseCamp/handler"
	"DatabaseCamp/repository"

	"github.com/gofiber/fiber/v2"
)

type router struct {
	app *fiber.App
}

var instantiated *router = nil

func New(app *fiber.App) *router {
	if instantiated == nil {
		instantiated = &router{app: app}
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
	userHandler := handler.NewUserHandler(controller, jwt)
	middleware := handler.NewJwtMiddleware(repository)
	group := r.app.Group("user")
	{
		group.Post("/register", userHandler.Register)
		group.Post("/login", userHandler.Login)
		group.Get("/info", middleware.JwtVerify, userHandler.GetInfo)
		group.Get("/profile/:id", middleware.JwtVerify, userHandler.GetProfile)
	}
}
