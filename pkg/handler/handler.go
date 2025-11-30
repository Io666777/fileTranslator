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

	router.Static("/front", "./front")
	router.GET("/", func(c *gin.Context) {
		c.File("./front/index.html")
	})

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp) // исправлено на kebab-case
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		files := api.Group("/files")
		{
			files.POST("/upload", h.uploadFile)               // ЗАГРУЗКА
			files.GET("/", h.getAllFiles)                     // СПИСОК
			files.GET("/:id/download", h.downloadFile)        // СКАЧИВАНИЕ
			files.POST("/:id/translate", h.createTranslation) // ПЕРЕВОД
			files.DELETE("/:id", h.deleteFile)                // УДАЛЕНИЕ
		}
	}

	return router
}
