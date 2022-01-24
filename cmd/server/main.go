package main

import (
	"context"
	"crypto/tls"
	"database-camp/internal/infrastructure/application"
	"database-camp/internal/infrastructure/database"
	"database-camp/internal/infrastructure/environment"
	"database-camp/internal/logs"
	"database-camp/internal/registry"
	"database-camp/internal/routes"
	"fmt"
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

	err := environment.New().Load(".env")
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

	liveName := fmt.Sprintf("tmp/live%d", os.Getpid())

	live, err := os.Create(liveName)
	if err != nil {
		logs.GetInstance().Error(err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {

		dir, _ := os.Getwd()

		// Create tls certificate
		cer, err := tls.LoadX509KeyPair(dir+"/backend.crt", dir+"/backend.key")
		if err != nil {
			log.Fatal(err)
		}
		config := &tls.Config{Certificates: []tls.Certificate{cer}}

		ln, err := tls.Listen("tcp", ":443", config)
		if err != nil {
			panic(err)
		}

		if err := app.Listener(ln); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()

	err = live.Close()
	if err != nil {
		logs.GetInstance().Error(err)
		return
	}

	err = os.Remove(liveName)
	if err != nil {
		logs.GetInstance().Error(err)
		return
	}

	err = app.Shutdown()
	if err != nil {
		logs.GetInstance().Error(err)
		return
	}
}
