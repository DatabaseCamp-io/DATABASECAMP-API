package main

import (
	"DatabaseCamp/database"
	"DatabaseCamp/logs"
	"DatabaseCamp/router"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func setupTimeZone() error {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return err
	}
	time.Local = location
	return nil
}

func getConfig() fiber.Config {
	return fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Database Camp",
	}
}

func setupFiber() error {
	app := fiber.New(getConfig())

	app.Use(cors.New())
	app.Use(recover.New())

	router.New(app)
	err := app.Listen(":" + os.Getenv("PORT"))

	return err
}

func main() {
	err := godotenv.Load()
	if err != nil {
		logs.New().Error(err)
		return
	}

	err = setupTimeZone()
	if err != nil {
		logs.New().Error(err)
		return
	}

	db := database.New()
	err = db.OpenConnection()
	if err != nil {
		logs.New().Error(err)
		return
	}
	defer db.CloseDB()

	err = setupFiber()
	if err != nil {
		logs.New().Error(err)
		return
	}
}
