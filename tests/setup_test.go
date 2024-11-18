package tests

import (
	"archive-api/internal/delivery/http"
	"archive-api/internal/repository"
	"archive-api/internal/usecase"
	"archive/zip"
	"bytes"
	"mime/multipart"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type testServer struct {
	router  *gin.Engine
	useCase *usecase.AnalyzeUseCase
}

func setupTestServer(t *testing.T) *testServer {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repo := repository.NewZipRepository()
	useCase := usecase.NewAnalyzeUseCase(repo)
	handler := http.NewHandler(useCase)

	router := gin.New()
	router.POST("/api/archive/information", handler.AnalyzeArchive)
	router.POST("/api/archive/files", handler.CreateArchive)
	router.POST("/api/mail/file", handler.SendMail)

	return &testServer{
		router:  router,
		useCase: useCase,
	}
}

func createTestZip(t *testing.T, files map[string][]byte) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	for name, content := range files {
		f, err := w.Create(name)
		require.NoError(t, err)
		_, err = f.Write(content)
		require.NoError(t, err)
	}

	err := w.Close()
	require.NoError(t, err)
	return buf.Bytes()
}

type testFile struct {
	name    string
	content []byte
}

func createMultipartFile(t *testing.T, tf testFile) *multipart.FileHeader {
	t.Helper()

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	part, err := writer.CreateFormFile("file", tf.name)
	require.NoError(t, err)

	_, err = part.Write(tf.content)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	reader := multipart.NewReader(&buffer, writer.Boundary())
	form, err := reader.ReadForm(32 << 20)
	require.NoError(t, err)

	return form.File["file"][0]
}

func createMultipartFiles(t *testing.T, files []testFile) []*multipart.FileHeader {
	t.Helper()
	var result []*multipart.FileHeader
	for _, tf := range files {
		result = append(result, createMultipartFile(t, tf))
	}
	return result
}
