package registry

import (
	"database-camp/internal/handler"
	"database-camp/internal/infrastructure/cache"
	"database-camp/internal/infrastructure/database"
	"database-camp/internal/middleware/jwt"
	"database-camp/internal/repositories"
	"database-camp/internal/services"
)

type middlewares struct {
	Jwt jwt.Jwt
}

type handlers struct {
	UserHandler     handler.UserHandler
	LearningHandler handler.LearningHandler
	ExamHandler     handler.ExamHandler
}

type Registry interface {
	GetMiddlewares() *middlewares
	GetHandlers() *handlers
}

type registry struct {
	middlewares middlewares
	handlers    handlers
}

func Regis() *registry {

	db := database.GetMySqlDBInstance()

	cache := cache.NewRedisClient()

	userRepo := repositories.NewUserRepository(db, cache)
	learningRepo := repositories.NewLearningRepository(db, cache)
	examRepo := repositories.NewExamRepository(db, cache)

	userService := services.NewUserService(userRepo, learningRepo)
	learningService := services.NewLearningService(learningRepo, userRepo)
	examService := services.NewExamService(examRepo, userRepo, learningRepo, cache)

	userHandler := handler.NewUserHandler(userService)
	learningHandler := handler.NewLearningHandler(learningService)
	examHandler := handler.NewExamHandler(examService)

	jwt := jwt.New(userRepo)

	return &registry{
		middlewares: middlewares{
			Jwt: jwt,
		},
		handlers: handlers{
			UserHandler:     userHandler,
			LearningHandler: learningHandler,
			ExamHandler:     examHandler,
		},
	}
}

func (r registry) GetMiddlewares() *middlewares {
	return &r.middlewares
}

func (r registry) GetHandlers() *handlers {
	return &r.handlers
}
