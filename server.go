package main

// server.go
/**
 * 	This file used to run server
 */

import (
	"DatabaseCamp/controllers"
	"DatabaseCamp/database"
	"DatabaseCamp/handlers"
	"DatabaseCamp/logs"
	"DatabaseCamp/middleware"
	"DatabaseCamp/repositories"
	"DatabaseCamp/router"
	"DatabaseCamp/services"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

/**
 * Set up server time zone
 *
 * @return the error of set up server time zone
 */
func setupTimeZone() error {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return err
	}
	time.Local = location
	return nil
}

/**
 * Config Web Server
 *
 * @return config of the web server
 */
func getConfig() fiber.Config {
	return fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Database Camp",
	}
}

/**
 * Set up web framework
 *
 * @return the error of set up web framework
 */
func setupFiber() error {

	// Create application
	app := fiber.New(getConfig())

	// Use middleware
	app.Use(cors.New())
	app.Use(recover.New())

	// Create database
	db := database.New()

	// Create Storage Service
	service := services.GetCloudStorageServiceInstance()

	// Create Repository
	userRepo := repositories.NewUserRepository(db)
	learningRepo := repositories.NewLearningRepository(db, service)
	examRepo := repositories.NewExamRepository(db)

	// Create Middleware
	jwt := middleware.NewJwtMiddleware(userRepo)

	// Create Controller
	learningController := controllers.NewLearningController(learningRepo, userRepo)
	examController := controllers.NewExamController(examRepo, userRepo)
	userController := controllers.NewUserController(userRepo)

	// Create Handler
	learningHandler := handlers.NewLearningHandler(learningController)
	examHandler := handlers.NewExamHandler(examController)
	userHandler := handlers.NewUserHandler(userController, jwt)

	// Create router
	router.New(app, examHandler, learningHandler, userHandler, jwt)

	// Running application
	err := app.Listen(":" + os.Getenv("PORT"))
	return err
}

/**
 * Main function of the application
 */
func main() {

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		logs.New().Error(err)
		return
	}

	// Setup time zone
	err = setupTimeZone()
	if err != nil {
		logs.New().Error(err)
		return
	}

	// Setup database
	db := database.New()
	err = db.OpenConnection()
	if err != nil {
		logs.New().Error(err)
		return
	}
	defer db.CloseDB()

	// Setup web framework
	err = setupFiber()
	if err != nil {
		logs.New().Error(err)
		return
	}
}
