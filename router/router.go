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
	r.setupExam(db, userRepo, jwt)
}

func (r router) setupExam(db database.IDatabase, userRepo repository.IUserRepository, jwt handler.IJwt) {
	repo := repository.NewExamRepository(db)
	controller := controller.NewExamController(repo, userRepo)
	examHandler := handler.NewExamHandler(controller)
	group := r.app.Group("exam")
	group.Use(jwt.JwtVerify)
	{
		group.Get("/proposition/:id", examHandler.GetExam)
		group.Get("/overview", examHandler.GetExamOverview)
		group.Get("/result/:id", examHandler.GetExamResult)
		group.Post("/check", examHandler.CheckExam)
	}
}

func (r router) setupLearning(db database.IDatabase, userRepo repository.IUserRepository, jwt handler.IJwt) {

	repo := repository.NewLearningRepository(db)
	service := services.GetAwsServiceInstance()
	controller := controller.NewLearningController(repo, userRepo, service)
	learningHandler := handler.NewLearningHandler(controller)
	group := r.app.Group("learning")
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

func (r router) setupUser(db database.IDatabase, repo repository.IUserRepository, jwt handler.IJwt) {
	controller := controller.NewUserController(repo)
	userHandler := handler.NewUserHandler(controller, jwt)
	group := r.app.Group("user")
	{
		group.Post("/register", userHandler.Register)
		group.Post("/login", userHandler.Login)
		group.Get("/info", jwt.JwtVerify, userHandler.GetInfo)
		group.Get("/profile/:id", jwt.JwtVerify, userHandler.GetProfile)
		group.Get("/ranking", jwt.JwtVerify, userHandler.GetUserRanking)
	}
}
