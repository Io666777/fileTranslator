package handlers

import (
	"fileTranslator/cmd/models"
	"fileTranslator/cmd/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateFile(c echo.Context) error {
    file := models.Filestr{}
    if err := c.Bind(&file); err != nil {
        return c.JSON(http.StatusBadRequest, "Invalid input")
    }
    
    newFile, err := repositories.CreateFile(file)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err.Error())
    }
    return c.JSON(http.StatusCreated, newFile)
}