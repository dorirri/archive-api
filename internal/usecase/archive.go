package usecase

import (
	"archive-api/internal/domain"
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
)

func (us *AnalyzeUseCase) CreateArchive(files []*multipart.FileHeader) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			return nil, &domain.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "file_read_error",
				Message: "Failed to read file: " + file.Filename,
			}
		}
		defer f.Close()

		content, err := io.ReadAll(f)
		if err != nil {
			return nil, &domain.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "file_read_error",
				Message: "Failed to read file content: " + file.Filename,
			}
		}
		mime := mimetype.Detect(content)
		if !domain.ALlowedTypes.ArchiveTypes[mime.String()] {
			return nil, &domain.APIError{
				Status:  http.StatusBadRequest,
				Code:    "invalid_file_type",
				Message: "Unsupported file type: " + file.Filename,
			}
		}
		fileWriter, err := zipWriter.Create(file.Filename)
		if err != nil {
			return nil, &domain.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "zip_creation_error",
				Message: "Failed to create zip entry",
			}
		}

		_, err = fileWriter.Write(content)
		if err != nil {
			return nil, &domain.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "zip_write_error",
				Message: "Failed to write to zip",
			}
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return nil, &domain.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "zip_close_error",
			Message: "Failed to finalize zip",
		}
	}
	return buf.Bytes(), nil
}
