package routes

import (
	"database-camp/internal/infrastructure/application"
	"database-camp/internal/registry"
	"net/http"
)

var (
	BuildCommit string
	BuildTime   string
)

type router struct {
	app   application.App
	route application.Router
	regis registry.Registry
}

var instantiated *router = nil

func NewRouter(app application.App, regis registry.Registry) *router {
	if instantiated == nil {
		instantiated = &router{
			app:   app,
			route: app.Group("api/v1"),
			regis: regis,
		}
		instantiated.setup()
	}
	return instantiated
}

func (r *router) setup() {
	r.setupProbe()
	r.setupUser()
	r.setupLearning()
	r.setupExam()
}

func (r *router) setupProbe() {
	r.route.Get("/x", func(c application.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"build_commit": BuildCommit,
			"build_time":   BuildTime,
		})
	})

	r.route.Get("/healthz", func(c application.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"status": "OK",
		})
	})
}

func (r *router) setupUser() {
	jwt := r.regis.GetMiddlewares().Jwt
	handler := r.regis.GetHandlers().UserHandler
	userRoute := r.route.Group("user")
	{
		userRoute.Post("/register", handler.Register)
		userRoute.Post("/login", handler.Login)
	}

	{
		userRoute.Get("/info", jwt.Verify, handler.GetOwnProfile)
		userRoute.Get("/profile/:id", jwt.Verify, handler.GetProfile)
		userRoute.Get("/ranking", jwt.Verify, handler.GetUserRanking)
		userRoute.Put("/profile", jwt.Verify, handler.Edit)
	}
}

func (r *router) setupLearning() {
	jwt := r.regis.GetMiddlewares().Jwt
	handler := r.regis.GetHandlers().LearningHandler
	learningRoute := r.route.Group("learning", jwt.Verify)
	activityRoute := learningRoute.Group("activity")
	{
		learningRoute.Get("/video/:id", handler.GetVideo)
		learningRoute.Get("/overview", handler.GetOverview)
		learningRoute.Get("/content/roadmap/:id", handler.GetContentRoadmap)
		learningRoute.Get("/recommend", handler.GetRecommend)
		learningRoute.Get("/spider", handler.GetSpiderData)

	}

	{
		activityRoute.Get("/:id", handler.GetActivity)
		activityRoute.Post("/hint/:id", handler.UseHint)
		activityRoute.Post("/check-answer", handler.CheckAnswer)
	}
}

func (r *router) setupExam() {
	jwt := r.regis.GetMiddlewares().Jwt
	handler := r.regis.GetHandlers().ExamHandler
	examRoute := r.route.Group("exam", jwt.Verify)
	{
		examRoute.Get("/proposition/:id", handler.GetExam)
		examRoute.Get("/overview", handler.GetExamOverview)
		examRoute.Get("/result/:id", handler.GetExamResult)
		examRoute.Post("/check", handler.CheckExam)
	}
}
