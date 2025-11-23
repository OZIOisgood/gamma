package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/OZIOisgood/gamma/internal/tools"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Storage struct {
	Client        *s3.Client
	PresignClient *s3.PresignClient
	Bucket        string
}

func New() *Storage {
	endpoint := tools.GetEnv("S3_ENDPOINT")
	accessKey := tools.GetEnv("S3_ACCESS_KEY")
	secretKey := tools.GetEnv("S3_SECRET_KEY")
	bucket := tools.GetEnv("S3_BUCKET")
	region := tools.GetEnv("S3_REGION")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               endpoint,
					HostnameImmutable: true,
				}, nil
			},
		)),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)

	return &Storage{
		Client:        client,
		PresignClient: presignClient,
		Bucket:        bucket,
	}
}

func (s *Storage) EnsureBucketExists(ctx context.Context) error {
	_, err := s.Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.Bucket),
	})
	if err == nil {
		return nil // Bucket exists
	}

	// If bucket doesn't exist, create it
	_, err = s.Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(s.Bucket),
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}
	return nil
}

func (s *Storage) EnsureBucketNotification(ctx context.Context) error {
	// Define the ARN for the NATS target configured in MinIO
	// Format: arn:minio:sqs::<REGION>:<ID>:nats
	// We used ID "gamma" in docker-compose
	arn := "arn:minio:sqs::gamma:nats"

	_, err := s.Client.PutBucketNotificationConfiguration(ctx, &s3.PutBucketNotificationConfigurationInput{
		Bucket: aws.String(s.Bucket),
		NotificationConfiguration: &types.NotificationConfiguration{
			QueueConfigurations: []types.QueueConfiguration{
				{
					QueueArn: aws.String(arn),
					Events: []types.Event{
						types.EventS3ObjectCreatedPut,
					},
					Filter: &types.NotificationConfigurationFilter{
						Key: &types.S3KeyFilter{
							FilterRules: []types.FilterRule{
								{
									Name:  types.FilterRuleNameSuffix,
									Value: aws.String(".mp4"),
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to set bucket notification: %w", err)
	}
	return nil
}

func (s *Storage) GetPresignedURL(ctx context.Context, key string) (string, error) {
	req, err := s.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return req.URL, nil
}
