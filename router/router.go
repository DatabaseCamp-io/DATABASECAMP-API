package router

import (
	"DatabaseCamp/handlers"
	"DatabaseCamp/middleware"

	"github.com/gofiber/fiber/v2"
)

type router struct {
	App    *fiber.App
	Router fiber.Router

	ExamHandler     handlers.IExamHandler
	LearningHandler handlers.ILearningHandler
	UserHandler     handlers.IUserHandler

	Jwt middleware.IJwt
}

var instantiated *router = nil

func New(
	app *fiber.App,
	examHandler handlers.IExamHandler,
	learningHandler handlers.ILearningHandler,
	userHandler handlers.IUserHandler,
	jwt middleware.IJwt,
) *router {
	if instantiated == nil {
		instantiated = &router{
			App:             app,
			Router:          app.Group("api/v1"),
			ExamHandler:     examHandler,
			LearningHandler: learningHandler,
			UserHandler:     userHandler,
			Jwt:             jwt,
		}
		instantiated.init()
	}
	return instantiated
}

// Set up init
func (r *router) init() {
	r.setupLearning()
	r.setupUser()
	r.setupExam()
}

func (r *router) setupExam() {
	group := r.Router.Group("exam")
	group.Use(r.Jwt.JwtVerify)
	{
		group.Get("/proposition/:id", r.ExamHandler.GetExam)
		group.Get("/overview", r.ExamHandler.GetExamOverview)
		group.Get("/result/:id", r.ExamHandler.GetExamResult)
		group.Post("/check", r.ExamHandler.CheckExam)
	}
}

func (r *router) setupLearning() {
	group := r.Router.Group("learning")
	group.Use(r.Jwt.JwtVerify)
	{
		group.Get("/video/:id", r.LearningHandler.GetVideo)
		group.Get("/overview", r.LearningHandler.GetOverview)
		group.Get("/content/roadmap/:id", r.LearningHandler.GetContentRoadmap)
		group.Get("/activity/:id", r.LearningHandler.GetActivity)
		group.Post("/activity/hint/:id", r.LearningHandler.UseHint)
		group.Post("/activity/matching/check-answer", r.LearningHandler.CheckMatchingAnswer)
		group.Post("/activity/multiple/check-answer", r.LearningHandler.CheckMultipleAnswer)
		group.Post("/activity/completion/check-answer", r.LearningHandler.CheckCompletionAnswer)
	}
}

func (r *router) setupUser() {
	group := r.Router.Group("user")
	{
		group.Get("/info", r.Jwt.JwtVerify, r.UserHandler.GetOwnProfile)
		group.Get("/profile/:id", r.Jwt.JwtVerify, r.UserHandler.GetProfile)
		group.Get("/ranking", r.Jwt.JwtVerify, r.UserHandler.GetUserRanking)
		group.Put("/profile", r.Jwt.JwtVerify, r.UserHandler.Edit)
		group.Post("/register", r.UserHandler.Register)
		group.Post("/login", r.UserHandler.Login)
	}
}
