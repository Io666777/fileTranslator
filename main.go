package main

import (
	"fileTranslator/cmd/handlers"
	"fileTranslator/cmd/storage"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Static("/", "web")

	storage.InitDB()

	// User routes
	e.POST("/user", handlers.CreateUser)

	// File routes
	e.POST("/file", handlers.CreateFile)

	e.Logger.Fatal(e.Start(":5500"))
}

//http://localhost:5500
