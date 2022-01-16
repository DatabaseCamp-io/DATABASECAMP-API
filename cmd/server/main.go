package main

import (
	"context"
	"database-camp/internal/infrastructure/application"
	"database-camp/internal/infrastructure/database"
	"database-camp/internal/infrastructure/environment"
	"database-camp/internal/logs"
	"database-camp/internal/registry"
	"database-camp/internal/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setupTimeZone() error {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return err
	}
	time.Local = location
	return nil
}

func main() {

	live, err := os.Create("tmp/live")
	if err != nil {
		logs.GetInstance().Error(err)
		return
	}

	err = environment.New().Load(".env")
	if err != nil {
		logs.GetInstance().Error(err)
		return
	}

	err = setupTimeZone()
	if err != nil {
		logs.GetInstance().Error(err)
		return
	}

	db := database.GetMySqlDBInstance()
	err = db.OpenConnection()
	if err != nil {
		logs.GetInstance().Error(err)
		return
	}
	defer db.CloseConnection()

	app := application.NewFiberApp()
	regis := registry.Regis()

	routes.NewRouter(app, regis)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := app.Listen(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	err = app.Shutdown()
	if err != nil {
		logs.GetInstance().Error(err)
		return
	}

	err = live.Close()
	if err != nil {
		logs.GetInstance().Error(err)
		return
	}

	err = os.Remove("tmp/live")
	if err != nil {
		logs.GetInstance().Error(err)
		return
	}
}
