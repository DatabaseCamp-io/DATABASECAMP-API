package router

// router.go
/**
 * 	This file used to define path of the API
 */

import (
	"DatabaseCamp/handlers"
	"DatabaseCamp/middleware"

	"github.com/gofiber/fiber/v2"
)

/**
 * This class manage API route of the application
 */
type router struct {
	App    *fiber.App   // Web framework application
	Router fiber.Router // Router of the web framework application

	ExamHandler     handlers.IExamHandler     // Exam handler for handle exam requests
	LearningHandler handlers.ILearningHandler // Learning handler for handle learning requests
	UserHandler     handlers.IUserHandler     // User handler for handle user requests

	Jwt middleware.IJwt // JWT middleware for verification
}

// Instance of router class for singleton pattern
var instantiated *router = nil

var (
	BuildCommit string
	BuildTime   string
)

/**
 * Constructor creates a new router instance or geting a router instance
 *
 * @param 	app 				Web framework application
 * @param 	examHandler 		Exam handler for handle exam requests
 * @param 	learningHandler 	Learning handler for handle learning requests
 * @param 	userHandler 		User handler for handle user requests
 * @param 	jwt 				JWT middleware for verification
 *
 * @return 	instance of router
 */
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

/**
 * Setup route path for each module
 */
func (r *router) init() {
	r.setup()
	r.setupLearning()
	r.setupUser()
	r.setupExam()
}

func (r *router) setup() {
	r.Router.Get("/x", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"build_commit": BuildCommit,
			"build_time":   BuildTime,
		})
	})

	r.Router.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("OK")
	})
}

/**
 * Setup route path for exam module
 */
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

/**
 * Setup route path for learning module
 */
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

/**
 * Setup route path for user module
 */
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
