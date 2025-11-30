package main

import (
	"filetranslation/internal"
	"filetranslation/pkg/handler"
	"filetranslation/pkg/repository"
	"filetranslation/pkg/service"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initialization config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer db.Close()

	repos := repository.NewRepository(db)

	// БЕРЕМ URL ИЗ КОНФИГА
	translationService := service.NewTranslationService(
		viper.GetString("translation.api_url"),
	)

	services := &service.Service{
		Authorization: service.NewAuthService(repos.Authorization),
		File:          service.NewFileService(repos.File),
		Translation:   translationService,
	}

	handlers := handler.NewHandler(services)

	srv := new(internal.Server)
	port := viper.GetString("port")
	if port == "" {
		port = "8080"
	}

	if err := srv.Run(port, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occurred while running http server: %s", err.Error())
	}
}

// ДОБАВИТЬ ЭТУ ФУНКЦИЮ
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
