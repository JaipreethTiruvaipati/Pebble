// Package main (s3_upload.go) provides S3 receipt storage used before scoring-service
// downloads the object for Google Vision OCR.
package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// S3Uploader uploads receipt images to the configured AWS S3 bucket and region.
type S3Uploader struct {
	bucketName string
	region     string
}

// NewS3Uploader returns an uploader bound to bucketName in region.
func NewS3Uploader(bucketName, region string) *S3Uploader {
	return &S3Uploader{
		bucketName: bucketName,
		region:     region,
	}
}

// UploadImage stores fileBytes under receipts/{uuid}.jpg and returns the S3 object key.
func (u *S3Uploader) UploadImage(ctx context.Context, fileBytes []byte, contentType string) (string, error) {
	// Generate a unique filename
	key := fmt.Sprintf("receipts/%s.jpg", uuid.New().String())
	
	log.Info().Str("bucket", u.bucketName).Str("key", key).Msg("uploading file to S3")
	
	// TODO: Phase 1 - Implement aws-sdk-go-v2 S3 PutObject
	
	return key, nil
}
