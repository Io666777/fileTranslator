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
	
	// Прямое подключение без конфига
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

	// Инициализация репозиториев, сервисов и хендлеров
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	// Запуск сервера
	srv := new(internal.Server)

	// Получаем порт из конфига с дефолтным значением
	port := viper.GetString("port")
	if port == "" {
		port = "8080"
	}
	
	if err := srv.Run(port, handlers.InitRoutes()); err != nil { // исправлено
		logrus.Fatalf("error occurred while running http server: %s", err.Error()) // исправлено
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}