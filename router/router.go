package router

import (
	"DatabaseCamp/controller"
	"DatabaseCamp/database"
	"DatabaseCamp/handler"
	"DatabaseCamp/repository"
	"DatabaseCamp/services"

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
	db := database.New()
	r.setupLearning(db)
	r.setupUser(db)
}

func (r router) setupLearning(db database.IDatabase) {
	repo := repository.NewLearningRepository(db)
	service := services.GetAwsServiceInstance()
	controller := controller.NewLearningController(repo, service)
	learningHandler := handler.NewLearningHandler(controller)
	group := r.app.Group("learning")
	{
		group.Get("/video/", learningHandler.GetVideo)
	}
}

func (r router) setupUser(db database.IDatabase) {
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
