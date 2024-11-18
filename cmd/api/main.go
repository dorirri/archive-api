package main

import (
	"archive-api/internal/delivery/http"
	"archive-api/internal/repository"
	"archive-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	repo := repository.NewZipRepository()
	useCase := usecase.NewAnalyzeUseCase(repo)
	handler := http.NewHandler(useCase)

	router := gin.Default()
	http.SetupRoutes(router, handler)

	router.Run(":8080")
}
