package r2

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Service struct {
	client       *s3.Client
	presigner    *s3.PresignClient
	bucketName   string
	publicDomain string
}

func NewR2Service(accountID, accessKeyID, secretAccessKey, bucketName, publicDomain string) (*R2Service, error) {
	var endpoint string

	// Use MinIO endpoint if accountID is "minio" (for local development)
	if accountID == "minio" {
		endpoint = "http://minio:9000"
	} else {
		endpoint = fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)
	}

	cfg := aws.Config{
		Region:      "auto",
		Credentials: credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	presigner := s3.NewPresignClient(client)

	return &R2Service{
		client:       client,
		presigner:    presigner,
		bucketName:   bucketName,
		publicDomain: publicDomain,
	}, nil
}

func (r *R2Service) GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, duration time.Duration) (string, error) {
	req, err := r.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = duration
	})
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (r *R2Service) GeneratePresignedDeleteURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	req, err := r.presigner.PresignDeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = duration
	})
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (r *R2Service) GeneratePresignedDownloadURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	req, err := r.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = duration
	})
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (r *R2Service) GetPublicURL(key string) string {
	// Remove leading slash if present
	key = strings.TrimPrefix(key, "/")
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(r.publicDomain, "/"), key)
}
