package handlers

import (
	"fileTranslator/cmd/models"
	"fileTranslator/cmd/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)


func CreateUser(c echo.Context) error {
  user := models.User{}
  c.Bind(&user)
  newUser, err := repositories.CreateUser(user)
  if err != nil {
    return c.JSON(http.StatusInternalServerError, err.Error())
  }
  return c.JSON(http.StatusCreated, newUser)
}

func CreateUserSimple(c echo.Context) error {
    return c.JSON(200, map[string]string{
        "id": "1", 
        "name": "Test User",
        "email": "test@example.com",
        "message": "User created (simulated)",
    })
}