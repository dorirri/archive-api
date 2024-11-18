package http

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, handler *Handler) {
	api := router.Group("api")
	{
		api.POST("/archive/information", handler.AnalyzeArchive)
		api.POST("/archive/files", handler.CreateArchive)
		api.POST("/mail/file", handler.SendMail)
	}
}
