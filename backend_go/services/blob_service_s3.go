package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewSupabaseS3Client(accessKey string, secretKey string,
	region string, endpoint string) (*s3.Client, error) {
	// log.Println("Attempting to load env ...")
	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	log.Fatalf("Failed to load env...\nError: %v", err)
	// }
	// accessKey := os.Getenv("SUPABASE_S3_ACCESS_KEY")
	// secretKey := os.Getenv("SUPABASE_S3_SECRET_KEY")
	// region := os.Getenv("SUPABASE_S3_REGION")
	// endpoint := os.Getenv("SUPABASE_S3_ENDPOINT")
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region), config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")))
	if err != nil {
		return nil, fmt.Errorf("Failed to load S3 config...: %w", err)
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})
	return client, nil
}

func Upload(ctx context.Context, s *s3.Client, bucketName string, supabase_url string,
    filePath string, objectPath string, ) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        log.Println("Unable to open file to upload:", err)
        return "", err
    }
    defer file.Close()
    _, err = s.PutObject(ctx, &s3.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(objectPath),
        Body:   file,
    })
    if err != nil {
        log.Println("Unable to upload file:", err)
        return "", err
    }

    blobURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s",
        supabase_url, bucketName, objectPath,)

    return blobURL, nil
}
