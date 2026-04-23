package storage

import (
	"bytes"
	"context"
	"fmt"
	"net/url"

	"poster-backend/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3 struct {
	client        *s3.Client
	bucket        string
	bucketPreview string
	endpoint      string
}

func NewS3(cfg *config.Config) (*S3, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		u, _ := url.Parse(cfg.S3.Endpoint)
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           u.String(),
			SigningRegion: cfg.S3.Region,
		}, nil
	})

	awsConf, err := awsCfg.LoadDefaultConfig(context.Background(),
		awsCfg.WithRegion(cfg.S3.Region),
		awsCfg.WithEndpointResolverWithOptions(customResolver),
		awsCfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.S3.AccessKey, cfg.S3.SecretKey, "")),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsConf, func(o *s3.Options) {
		o.UsePathStyle = true // MinIO
	})

	return &S3{
		client:        client,
		bucket:        cfg.S3.Bucket,
		bucketPreview: cfg.S3.BucketPreview,
		endpoint:      cfg.S3.Endpoint,
	}, nil
}

func (s *S3) UploadOrig(ctx context.Context, data []byte, contentType string) (string, error) {
	key := uuid.New().String()
	return s.upload(ctx, s.bucket, key, data, contentType)
}

func (s *S3) UploadPreview(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	return s.upload(ctx, s.bucketPreview, key, data, contentType)
}

func (s *S3) upload(ctx context.Context, bucket, key string, data []byte, contentType string) (string, error) {
	ext := ".png"
	if contentType == "image/jpeg" {
		ext = ".jpg"
	}

	key = key + ext
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", s.endpoint, bucket, key), nil
}
