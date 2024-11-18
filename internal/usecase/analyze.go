package usecase

import (
	"archive-api/internal/domain"
	"archive-api/pkg/validator"
)

type AnalyzeUseCase struct {
	repo domain.ArchiveRepository
}

func NewAnalyzeUseCase(repo domain.ArchiveRepository) *AnalyzeUseCase {
	return &AnalyzeUseCase{repo: repo}
}

func (uc *AnalyzeUseCase) Execute(content []byte, filename string) (*domain.ArchiveInfo, error) {
	if err := validator.ValidateZipFile(filename, content); err != nil {
		return nil, err
	}
	return uc.repo.AnalyzeArchive(content, filename)
}
