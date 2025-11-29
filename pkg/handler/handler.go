package handler

import "github.com/gin-gonic/gin"

type Handler struct {
	// позже добавим сервисы
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/signUp", h.signUp)//h.signIn undefined (type *Handler has no field or method signIn)
		auth.POST("/signIn", h.signIn)//h.signIn undefined (type *Handler has no field or method signIn)
	}

	api := router.Group("/api")
	{
		files := api.Group("/files")
		{
			files.POST("/", h.createFile)
			files.GET("/", h.getAllFiles) // теперь будет работать
			files.GET("/:id", h.getFileById) // теперь будет работать
			files.POST("/:id/translations", h.createTranslation)
			files.DELETE("/:id", h.deleteFile)
		}
	}

	return router
}