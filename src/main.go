package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// temp location to store image and thumbnail
const tmp = "/tmp/"

// S3 Session to use
var sess = session.Must(session.NewSession())

// Create a downloader with session and default option
var downloader = s3manager.NewDownloader(sess)

func handler(ctx context.Context, req events.S3Event) error {
	return nil
}

func main() {
	lambda.Start(handler)
}
