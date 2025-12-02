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

	router.Static("/css", "./front/css")
	router.Static("/js", "./front/js")
	router.Static("/assets", "./front/assets")

	router.GET("/", func(c *gin.Context) {
		c.File("./front/index.html")
	})

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		files := api.Group("/files")
		{
			files.POST("/upload", h.uploadFile)
			files.GET("/", h.getAllFiles)
			files.GET("/:id/download", h.downloadFile)
			files.POST("/:id/translate", h.createTranslation)
			files.DELETE("/:id", h.deleteFile)
		}
	}

	return router
}
