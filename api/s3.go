package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"io"
	"log"
)

const Region = "us-west-2"
const Bucket = "pbcp"

var config aws.Config = aws.Config{
	Region: aws.String(Region),
}

func upload(key string, content io.Reader) error {
	uploader := s3manager.NewUploader(session.New(&config))
	res, err := uploader.Upload(&s3manager.UploadInput{
		Body:   content,
		Bucket: aws.String(Bucket),
		Key:    aws.String(key),
	})
	if err == nil {
		log.Printf("Upload: %s", res.Location)
	}
	return err
}

func delete(key string) error {
	svc := s3.New(session.New(&config))
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(Bucket),
		Key: aws.String(key),
	}
	_, err := svc.DeleteObject(params)
	return err
}
