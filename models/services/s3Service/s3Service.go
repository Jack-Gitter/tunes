package s3Service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct { 
    client *s3.Client
}

type Is3Service interface {
    UploadToBucket() error
    InitClient() 
}

func(s3Service *S3Service) InitClient() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
        panic(err)
    }
    fmt.Println(cfg)
    client := s3.NewFromConfig(cfg)
    s3Service.client = client
}

func(s3Service *S3Service) UploadToBucket() error {
    reader := strings.NewReader("hello world!")
    putObjectInput := &s3.PutObjectInput{
        Bucket: aws.String("tunes-profile-pictures"),
        Key: aws.String("hi"),
        Body: reader,
    }

    putObjectOutput, err := s3Service.client.PutObject(context.Background(), putObjectInput)
    if err != nil {
        return err
    }

    fmt.Println(putObjectOutput)

    return nil
}
