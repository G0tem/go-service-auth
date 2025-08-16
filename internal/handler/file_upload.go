package handler

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func (h *Handler) FileUpload(file *multipart.FileHeader) (*minio.UploadInfo, error) {
	ctx := context.Background()
	bucketName := h.cfg.S3AvatarsBucketName

	// Get Buffer from file
	buffer, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer buffer.Close()

	endpoint := h.cfg.S3Endpoint
	accessKeyID := h.cfg.S3AccessKey
	secretAccessKey := h.cfg.S3SecretAccessKey

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		BucketLookup: minio.BucketLookupPath,
		Region:       h.cfg.S3Region,
		Secure:       true,
	})
	if err != nil {
		return nil, err
	}

	objectName := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(file.Filename))
	fileBuffer := buffer
	contentType := file.Header["Content-Type"][0]
	fileSize := file.Size

	// Upload the file with PutObject
	info, err := minioClient.PutObject(ctx, bucketName, objectName, fileBuffer, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"FileName": file.Filename,
		},
	})
	if err != nil {
		return nil, err
	}

	return &info, nil
}
