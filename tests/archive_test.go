package tests

import (
	"archive-api/internal/domain"
	"archive-api/internal/repository"
	"archive-api/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateArchive(t *testing.T) {
	pngContent := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	jpegContent := []byte{0xFF, 0xD8, 0xFF, 0xE0}

	tests := []struct {
		name          string
		files         []testFile
		expectedError string
	}{
		{
			name: "successful_archive_creation",
			files: []testFile{
				{name: "test1.png", content: pngContent},
				{name: "test2.jpg", content: jpegContent},
			},
			expectedError: "",
		},
		{
			name: "invalid_file_type",
			files: []testFile{
				{name: "test.exe", content: []byte{0x4D, 0x5A}},
			},
			expectedError: "Unsupported file type: test.exe",
		},
		{

			name:          "empty_file_list",
			files:         []testFile{},
			expectedError: "",
		},
		{
			name: "mixed_valid_and_invalid_files",
			files: []testFile{
				{name: "valid.png", content: pngContent},
				{name: "invalid.txt", content: []byte("text file")},
			},
			expectedError: "Unsupported file type: invalid.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := createMultipartFiles(t, tt.files)

			repo := repository.NewZipRepository()
			uc := usecase.NewAnalyzeUseCase(repo)

			result, err := uc.CreateArchive(files)

			if tt.expectedError != "" {
				require.Error(t, err)
				apiErr, ok := err.(*domain.APIError)
				require.True(t, ok, "error should be of type *domain.APIError")
				assert.Contains(t, apiErr.Message, tt.expectedError)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			if len(tt.files) > 0 && tt.expectedError == "" {
				info, err := repo.AnalyzeArchive(result, "test_archive.zip")
				require.NoError(t, err)

				assert.Equal(t, float64(len(tt.files)), info.TotalFiles)

				fileNames := make(map[string]bool)
				for _, f := range info.Files {
					fileNames[f.FilePath] = true
				}

				for _, expectedFile := range tt.files {
					assert.True(t, fileNames[expectedFile.name], "File %s should be present in archive", expectedFile.name)
				}

				var expectedTotalSize float64
				for _, f := range tt.files {
					expectedTotalSize += float64(len(f.content))
				}
				assert.Equal(t, expectedTotalSize, info.TotalSize)
			}
		})
	}
}
