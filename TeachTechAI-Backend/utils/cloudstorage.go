package utils

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"teach-tech-ai/helpers"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

var (
	bucketName     = helpers.MustGetenv("STORAGE_BUCKET_NAME")
	serviceKeyPath = helpers.MustGetenv("STORAGE_BUCKET_KEY_PATH")
)

func generateObjectName(ext string, userID uuid.UUID) string {
	t := time.Now().UTC()
	return fmt.Sprintf("upload-%d%d%d-%d%d%d-%s%s", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), userID, ext)
}

func UploadFileToCloud(ctx context.Context, localFilePath string, userID uuid.UUID) (string, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(serviceKeyPath))
	if err != nil {
		return "", err
	}
	defer client.Close()

	file, err := os.Open(localFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(localFilePath)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", err
	}
	objectName := generateObjectName(ext, userID)

	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)

	// Determine content type based on file extension
	wc.ContentType = mime.TypeByExtension(ext)
	if wc.ContentType == "" {
		wc.ContentType = "application/octet-stream" // Default to binary stream if unknown
	}

	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	return objectName, nil
}

func DownloadFileFromCloud(ctx context.Context, objectName, localFilePath string) error {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(serviceKeyPath))
	if err != nil {
		return err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	file, err := os.Create(localFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()

	if _, err := io.Copy(file, rc); err != nil {
		return err
	}

	return nil
}

func DeleteFileFromCloud(ctx context.Context, fileName string) error {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(serviceKeyPath))
	if err != nil {
		return err
	}
	defer client.Close()

	object := client.Bucket(bucketName).Object(fileName)
	if err := object.Delete(ctx); err != nil {
		return err
	}

	return nil
}
