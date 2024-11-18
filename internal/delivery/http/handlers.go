package http

import (
	"io"
	"net/http"

	"archive-api/internal/domain"
	"archive-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	analyzeUseCase *usecase.AnalyzeUseCase
}

func NewHandler(analyzeUseCase *usecase.AnalyzeUseCase) *Handler {
	return &Handler{
		analyzeUseCase: analyzeUseCase,
	}
}

func (h *Handler) AnalyzeArchive(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.APIError{
			Code:    "missing_file",
			Message: "No file was uploaded",
		})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIError{
			Code:    "file_read_error",
			Message: "Failed to read uploaded file",
		})
		return
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.APIError{
			Code:    "file_read_error",
			Message: "Failed to read uploaded file",
		})
		return
	}
	info, err := h.analyzeUseCase.Execute(content, file.Filename)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Status, apiErr)
			return
		}
		c.JSON(http.StatusInternalServerError, domain.APIError{
			Code:    "internal_error",
			Message: "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *Handler) CreateArchive(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.APIError{
			Code:    "invalid_request",
			Message: "Invalid multipart form data",
		})
		return
	}

	files := form.File["files[]"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, domain.APIError{
			Code:    "no_files",
			Message: "No files provided",
		})
		return
	}
	archiveBytes, err := h.analyzeUseCase.CreateArchive(files)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Status, apiErr)
			return
		}
		c.JSON(http.StatusInternalServerError, domain.APIError{
			Code:    "internal_error",
			Message: "Failed to create archive",
		})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=archive.zip")
	c.Data(http.StatusOK, "application/zip", archiveBytes)
}

func (h *Handler) SendMail(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.APIError{
			Code:    "missing_file",
			Message: "No file provided",
		})
		return
	}

	emails := c.PostForm("emails")
	if emails == "" {
		c.JSON(http.StatusBadRequest, domain.APIError{
			Code:    "missing_emails",
			Message: "No email addresses provided",
		})
		return
	}
	err = h.analyzeUseCase.SendMail(file, emails)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Status, apiErr)
			return
		}
		c.JSON(http.StatusInternalServerError, domain.APIError{
			Code:    "mail_send_error",
			Message: "Failed to send email",
		})
		return
	}
	c.Status(http.StatusOK)
}
