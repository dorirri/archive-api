package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, handler *Handler) {
	router.SetTrustedProxies(nil)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hi :)",
		})
	})

	api := router.Group("api")
	{
		api.POST("/archive/information", handler.AnalyzeArchive)
		api.POST("/archive/files", handler.CreateArchive)
		api.POST("/mail/file", handler.SendMail)
	}
}
