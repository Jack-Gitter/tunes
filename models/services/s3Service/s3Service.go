package s3Service

import (
	"context"
	"log"
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
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
        panic(err)
    }
    client := s3.NewFromConfig(cfg)
    s3Service.client = client
}

func(s3Service *S3Service) UploadToBucket() error {
	output, err := s3Service.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String("tunes-profile-pictures"),
	})

	if err != nil {
		log.Fatal(err)
	}

    log.Println("first page results:")

	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}
    return nil
}
