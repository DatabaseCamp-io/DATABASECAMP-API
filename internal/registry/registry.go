package registry

import (
	"database-camp/internal/handler"
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

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	learningRepo := repositories.NewLearningRepository(db)
	learningService := services.NewLearningService(learningRepo, userRepo)
	learningHandler := handler.NewLearningHandler(learningService)

	examRepo := repositories.NewExamRepository(db)
	examService := services.NewExamService(examRepo, userRepo, learningRepo)
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
