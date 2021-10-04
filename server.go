package main

import (
	"DatabaseCamp/database"
	"DatabaseCamp/logs"
	"DatabaseCamp/router"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var allowOrigin []string

func setupTimeZone() error {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return err
	}
	time.Local = location
	return nil
}

func getCORSConfig() middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins: allowOrigin,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPut,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
	}
}

func setupEcho() error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(getCORSConfig()))
	e.Pre(middleware.AddTrailingSlash())

	err := e.Start(":" + os.Getenv("PORT"))
	if err != nil {
		return err
	}

	router.New(e)
	return nil
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
	defer db.Close()

	err = setupEcho()
	if err != nil {
		logs.New().Error(err)
		return
	}
}
