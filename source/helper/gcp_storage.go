package helper

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"sarana-dafa-ai-service/storage/env"
	"time"

	"cloud.google.com/go/storage"
)

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

var uploader *ClientUploader

func NewGcpStorage() {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	uploader = &ClientUploader{
		cl:         client,
		bucketName: os.Getenv(env.GOOGLE_BUCKET_NAME),
		projectID:  os.Getenv(env.GOOGLE_PROJECT_ID),
		uploadPath: os.Getenv(env.GOOGLE_BUCKET_FOLDER_NAME),
	}
}

// UploadFile uploads an object
func (c *ClientUploader) UploadFileToBucket(file multipart.File, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

func (c *ClientUploader) DeleteFileFromBucket(object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(object)

	attrs, err := wc.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %w", err)
	}
	wc = wc.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := wc.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %w", object, err)
	}

	return nil
}
