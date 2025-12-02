package main

import (
	"filetranslation/internal"
	"filetranslation/pkg/handler"
	"filetranslation/pkg/repository"
	"filetranslation/pkg/service"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer db.Close()

	repos := repository.NewRepository(db)

	// БЕРЕМ URL ИЗ КОНФИГА
	translationService := service.NewTranslationService(
		os.Getenv("API_URL"),
	)

	services := &service.Service{
		Authorization: service.NewAuthService(repos.Authorization),
		File:          service.NewFileService(repos.File),
		Translation:   translationService,
	}

	handlers := handler.NewHandler(services)

	srv := new(internal.Server)
	port := "8080"

	if err := srv.Run(port, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occurred while running http server: %s", err.Error())
	}
}
