package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// S3Uploader handles uploading receipt images to AWS S3.
type S3Uploader struct {
	bucketName string
	region     string
}

// NewS3Uploader creates a new uploader instance.
func NewS3Uploader(bucketName, region string) *S3Uploader {
	return &S3Uploader{
		bucketName: bucketName,
		region:     region,
	}
}

// UploadImage uploads the image bytes and returns the S3 Key.
func (u *S3Uploader) UploadImage(ctx context.Context, fileBytes []byte, contentType string) (string, error) {
	// Generate a unique filename
	key := fmt.Sprintf("receipts/%s.jpg", uuid.New().String())
	
	log.Info().Str("bucket", u.bucketName).Str("key", key).Msg("uploading file to S3")
	
	// TODO: Phase 1 - Implement aws-sdk-go-v2 S3 PutObject
	
	return key, nil
}
