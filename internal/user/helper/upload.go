package helper

import (
	"bytes"
	cloud "cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/SemmiDev/blog/config"
	"github.com/SemmiDev/blog/internal/common/logger"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
	"time"
)

// UploadResult is the result of an upload.
type UploadResult struct {
	Path  string
	Error error
}

// UploadImage uploads an image to the cloud storage
func UploadImage(ctx context.Context, cloudStorage *cloud.Client, file *multipart.FileHeader, email string) chan UploadResult {
	result := make(chan UploadResult)
	go func() {
		f, err := file.Open()
		if err != nil {
			result <- UploadResult{
				Path:  "",
				Error: errors.New("failed to open file"),
			}
		}
		defer f.Close()

		if file.Size > 1000000 {
			result <- UploadResult{
				Path:  "",
				Error: errors.New("image is too large"),
			}
		}

		buffer := make([]byte, file.Size)
		f.Read(buffer)

		fileType := http.DetectContentType(buffer)
		if !strings.HasPrefix(fileType, "image") {
			result <- UploadResult{
				Path:  "",
				Error: errors.New("file is not an image"),
			}
		}

		fileName := DefineFileName("profile", email, path.Ext(file.Filename))
		fileBytes := bytes.NewReader(buffer)
		bucketName := config.Env.FirebaseBucketName

		writer := cloudStorage.Bucket(bucketName).Object(fileName).NewWriter(ctx)
		defer writer.Close()

		if _, err := io.Copy(writer, fileBytes); err != nil {
			result <- UploadResult{
				Path:  "",
				Error: errors.New("failed to upload image"),
			}
		}

		uploadedImageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, fileName)
		logger.Log.Info().Msg("image uploaded")
		result <- UploadResult{
			Path:  uploadedImageURL,
			Error: nil,
		}
	}()

	return result
}

// DefineFileName generates a unique file name for the image.
func DefineFileName(kind string, email string, ext string) string {
	var filename bytes.Buffer
	filename.WriteString(email)
	filename.WriteString("_")
	filename.WriteString(kind)
	filename.WriteString("_")
	filename.WriteString(fmt.Sprintf("%d", time.Now().Unix()))
	filename.WriteString(ext)
	return filename.String()
}
