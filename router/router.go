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
	userRepo := repository.NewUserRepository(db)
	jwt := handler.NewJwtMiddleware(userRepo)
	r.setupLearning(db, userRepo, jwt)
	r.setupUser(db, userRepo, jwt)
}

func (r router) setupLearning(db database.IDatabase, userRepo repository.IUserRepository, jwt handler.IJwt) {

	repo := repository.NewLearningRepository(db)
	service := services.GetAwsServiceInstance()
	controller := controller.NewLearningController(repo, userRepo, service)
	learningHandler := handler.NewLearningHandler(controller)
	group := r.app.Group("learning")
	group.Use(jwt.JwtVerify)
	{
		group.Get("/video", learningHandler.GetVideo)
		group.Get("/overview", learningHandler.GetOverview)
		group.Get("/overview", learningHandler.GetActivity)
	}
}

func (r router) setupUser(db database.IDatabase, repo repository.IUserRepository, jwt handler.IJwt) {
	controller := controller.NewUserController(repo)
	userHandler := handler.NewUserHandler(controller, jwt)
	group := r.app.Group("user")
	{
		group.Post("/register", userHandler.Register)
		group.Post("/login", userHandler.Login)
		group.Get("/info", jwt.JwtVerify, userHandler.GetInfo)
		group.Get("/profile/:id", jwt.JwtVerify, userHandler.GetProfile)
	}
}
