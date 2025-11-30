package handler

import (
	"filetranslation/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp) // исправлено на kebab-case
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api")
	{
		files := api.Group("/files")
		{
			files.POST("/", h.createFile)
			files.GET("/", h.getAllFiles)
			files.GET("/:id", h.getFileById)
			files.POST("/:id/translations", h.createTranslation)
			files.DELETE("/:id", h.deleteFile)
		}
	}

	return router
}