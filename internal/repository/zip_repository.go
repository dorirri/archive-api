package repository

import (
	"archive-api/internal/domain"
	"archive/zip"
	"bytes"
	"io"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
)

type ZipRepository struct{}

func NewZipRepository() *ZipRepository {
	return &ZipRepository{}
}

func (r *ZipRepository) AnalyzeArchive(content []byte, filename string) (*domain.ArchiveInfo, error) {
	reader := bytes.NewReader(content)
	zipReader, err := zip.NewReader(reader, int64(len(content)))
	if err != nil {
		return nil, &domain.APIError{
			Status:  http.StatusBadRequest,
			Code:    "invalid_archive",
			Message: "The uploaded file is not a valid ZIP archive",
		}
	}

	var files []domain.FileInfo
	var totalSize float64

	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			continue
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			continue
		}
		mime := mimetype.Detect(content)
		size := float64(file.UncompressedSize64)
		totalSize += size
		files = append(files, domain.FileInfo{
			FilePath: file.Name,
			Size:     size,
			Mimetype: mime.String(),
		})
	}

	return &domain.ArchiveInfo{
		FileName:    filename,
		ArchiveSize: float64(len(content)),
		TotalSize:   totalSize,
		TotalFiles:  float64(len(files)),
		Files:       files,
	}, nil
}
