package usecase

import (
	"archive-api/internal/domain"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"gopkg.in/gomail.v2"
)

func (uc *AnalyzeUseCase) SendMail(file *multipart.FileHeader, emails string) error {
	emailList := strings.Split(emails, ",")

	f, err := file.Open()
	if err != nil {
		return &domain.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "file_read_error",
			Message: "Failed to read file",
		}
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return &domain.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "file_read_error",
			Message: "Failed to read file content",
		}
	}
	mime := mimetype.Detect(content)
	if !domain.ALlowedTypes.EmailTypes[mime.String()] {
		return &domain.APIError{
			Status:  http.StatusBadRequest,
			Code:    "invalid_file_type",
			Message: "Unsupported file type for email",
		}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_FROM"))
	m.SetHeader("To", emailList...)
	m.SetHeader("Subject", "File Attachment")
	m.SetBody("text/plain", "Please find the attached file.")
	m.Attach(file.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(content)
		return err
	}))

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
	)

	return d.DialAndSend(m)
}
