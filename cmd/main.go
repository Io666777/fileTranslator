package main

import (
	"filetranslation/internal"
	"filetranslation/pkg/handler"
	"filetranslation/pkg/repository"
	"filetranslation/pkg/service"
	"log"

	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil{
		log.Fatalf("error initialization config: %s", err.Error())
	}
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	
	srv := new(internal.Server) // исправлено file.Server -> internal.Server
	if err := srv.Run(viper.GetSring("8000"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occurred while running server: %s", err.Error())
	}
}

func initConfig() error{
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}