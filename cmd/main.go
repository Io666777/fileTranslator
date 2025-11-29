package main

import (
	"filetranslation/pkg/handler"
	"log"
)

func main() {
	handlers := handler.NewHandler()
	
	router := handlers.InitRoutes()
	
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("error occurred while running server: %s", err.Error())
	}
}