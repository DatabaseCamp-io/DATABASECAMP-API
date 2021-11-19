package router

import (
	"DatabaseCamp/controllers"
	"DatabaseCamp/database"
	"DatabaseCamp/handlers"
	"DatabaseCamp/middleware"
	"DatabaseCamp/repositories"
	"DatabaseCamp/services"

	"github.com/gofiber/fiber/v2"
)

type router struct {
	app *fiber.App
	api fiber.Router
}

var instantiated *router = nil

func New(app *fiber.App) *router {
	if instantiated == nil {
		instantiated = &router{app: app}
		instantiated.init()
	}
	return instantiated
}

func (r *router) init() {
	db := database.New()
	userRepo := repositories.NewUserRepository(db)
	jwt := middleware.NewJwtMiddleware(userRepo)
	r.api = r.app.Group("api/v1")
	r.setupLearning(db, userRepo, jwt)
	r.setupUser(db, userRepo, jwt)
	r.setupExam(db, userRepo, jwt)
}

func (r *router) setupExam(db database.IDatabase, userRepo repositories.IUserRepository, jwt middleware.IJwt) {
	repo := repositories.NewExamRepository(db)
	controller := controllers.NewExamController(repo, userRepo)
	examHandler := handlers.NewExamHandler(controller)
	group := r.api.Group("exam")
	group.Use(jwt.JwtVerify)
	{
		group.Get("/proposition/:id", examHandler.GetExam)
		group.Get("/overview", examHandler.GetExamOverview)
		group.Get("/result/:id", examHandler.GetExamResult)
		group.Post("/check", examHandler.CheckExam)
	}
}

func (r *router) setupLearning(db database.IDatabase, userRepo repositories.IUserRepository, jwt middleware.IJwt) {

	repo := repositories.NewLearningRepository(db)
	service := services.GetAwsServiceInstance()
	controller := controllers.NewLearningController(repo, userRepo, service)
	learningHandler := handlers.NewLearningHandler(controller)
	group := r.api.Group("learning")
	group.Use(jwt.JwtVerify)
	{
		group.Get("/video/:id", learningHandler.GetVideo)
		group.Get("/overview", learningHandler.GetOverview)
		group.Get("/content/roadmap/:id", learningHandler.GetContentRoadmap)
		group.Get("/activity/:id", learningHandler.GetActivity)
		group.Post("/activity/hint/:id", learningHandler.UseHint)
		group.Post("/activity/matching/check-answer", learningHandler.CheckMatchingAnswer)
		group.Post("/activity/multiple/check-answer", learningHandler.CheckMultipleAnswer)
		group.Post("/activity/completion/check-answer", learningHandler.CheckCompletionAnswer)
	}
}

func (r *router) setupUser(db database.IDatabase, repo repositories.IUserRepository, jwt middleware.IJwt) {
	controller := controllers.NewUserController(repo)
	userHandler := handlers.NewUserHandler(controller, jwt)
	group := r.api.Group("user")
	{
		group.Get("/info", jwt.JwtVerify, userHandler.GetOwnProfile)
		group.Get("/profile/:id", jwt.JwtVerify, userHandler.GetProfile)
		group.Get("/ranking", jwt.JwtVerify, userHandler.GetUserRanking)
		group.Put("/profile", jwt.JwtVerify, userHandler.Edit)
		group.Post("/register", userHandler.Register)
		group.Post("/login", userHandler.Login)
	}
}
