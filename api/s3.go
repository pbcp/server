package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"io"
	"log"
)

func upload(key string, content io.Reader) error {
	uploader := s3manager.NewUploader(
		session.New(&aws.Config{Region: aws.String("us-west-2")}),
	)
	res, err := uploader.Upload(&s3manager.UploadInput{
		Body:   content,
		Bucket: aws.String("pbcp"),
		Key:    aws.String(key),
	})
	if err == nil {
		log.Printf("Upload: %s", res.Location)
	}
	return err
}
