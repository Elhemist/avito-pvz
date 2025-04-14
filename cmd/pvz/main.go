package main

import (
	"os"
	"pvz-test/internal/app"
	"pvz-test/internal/handler"
	"pvz-test/internal/repository"
	"pvz-test/internal/service"
	"pvz-test/pkg/httpserver"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Loading env variables error: %s", err.Error())
	}

	app.Migrations()

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Usename:  os.Getenv("POSTGRES_USERNAME"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DATABASE"),
		SSLmode:  "disable",
	})
	if err != nil {
		logrus.Fatalf("DB init fail: %s", err.Error())
	}
	logrus.Info("DB init complete")

	repos := repository.NewRepository(db)
	service := service.NewService(repos)
	handlers := handler.NewHandler(service)

	srv := new(httpserver.Server)
	if err := srv.Start(os.Getenv("SERVER_ADDRESS"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Running error: %s", err.Error())
	}

}
