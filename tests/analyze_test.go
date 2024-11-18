package tests

import (
	"archive-api/internal/domain"
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalyzeArchive(t *testing.T) {
	server := setupTestServer(t)

	tests := []struct {
		name         string
		zipFiles     map[string][]byte
		expectedCode int
		validate     func(*testing.T, *domain.ArchiveInfo)
	}{
		{
			name: "valid zip with multiple files",
			zipFiles: map[string][]byte{
				"test.docx": []byte("docx content"),
				"test.jpg":  []byte{0xFF, 0xD8, 0xFF, 0xE0}, 
			},
			expectedCode: http.StatusOK,
			validate: func(t *testing.T, info *domain.ArchiveInfo) {
				require.Equal(t, float64(2), info.TotalFiles)
				require.Len(t, info.Files, 2)
			},
		},
		{
			name:         "invalid zip",
			zipFiles:     nil,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var content []byte
			if tt.zipFiles != nil {
				content = createTestZip(t, tt.zipFiles)
			} else {
				content = []byte("invalid zip")
			}

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("file", "test.zip")
			require.NoError(t, err)
			_, err = part.Write(content)
			require.NoError(t, err)
			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/api/archive/information", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			w := httptest.NewRecorder()

			server.router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var response domain.ArchiveInfo
				err = json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
				tt.validate(t, &response)
			}
		})
	}
}
