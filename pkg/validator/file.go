package validator

import (
	"net/http"
	"path/filepath"

	"archive-api/internal/domain"
)

const (
	maxFileSize = 100 * 1024 * 1024
)

func ValidateZipFile(filename string, content []byte) error {
	if ext := filepath.Ext(filename); ext != ".zip" {
		return &domain.APIError{
			Status:  http.StatusBadRequest,
			Code:    "invalid_file_type",
			Message: "Only ZIP files are supported",
		}
	}

	if len(content) == 0 {
		return &domain.APIError{
			Status:  http.StatusBadRequest,
			Code:    "empty_file",
			Message: "The uploaded file is empty",
		}
	}

	if len(content) > maxFileSize {
		return &domain.APIError{
			Status:  http.StatusBadRequest,
			Code:    "file_too_large",
			Message: "File exceeds the max limit of 100mb",
		}
	}

	if len(content) < 4 {
		return &domain.APIError{
			Status:  http.StatusBadRequest,
			Code:    "invalid_zip_format",
			Message: "Invalid ZIP file format",
		}
	}
	signature := content[:4]
	isValidSignature := (signature[0] == 0x50 && signature[1] == 0x4B &&
		((signature[2] == 0x03 && signature[3] == 0x04) || //norm
			(signature[2] == 0x05 && signature[3] == 0x06) || //pustoi
			(signature[2] == 0x07 && signature[3] == 0x08)))
	if !isValidSignature {
		return &domain.APIError{
			Status:  http.StatusBadRequest,
			Code:    "invalid_zip_format",
			Message: "Invalid ZIP file format",
		}
	}

	return nil
}
