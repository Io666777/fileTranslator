package main

import (
	"fileTranslator/cmd/handlers"
	"fileTranslator/cmd/storage"

	"github.com/labstack/echo/v4"
)

func main()  {
    e:= echo.New()
    e.GET("/", handlers.Home)
    storage.InitDB()
    e.Logger.Fatal(e.Start(":5500"))
}