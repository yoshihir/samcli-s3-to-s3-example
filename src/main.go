package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io/ioutil"
	"os"
)

type SampleData struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
}

func extract(file *os.File) ([]SampleData, error) {
	gzipReader, _ := gzip.NewReader(file)
	defer gzipReader.Close()

	raw, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		println(err.Error())
	}

	var data []SampleData
	err = json.Unmarshal(raw, &data)
	if err != nil {
		println(err.Error())
	}

	return data, nil
}

func s3Download(bucket string, key string) (f *os.File, err error) {
	var sess = session.Must(session.NewSession(&aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(os.Getenv("REGION")),
		Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
	}))

	tmpfile, _ := ioutil.TempFile("/tmp", "srctmp_")
	defer os.Remove(tmpfile.Name())

	// ダウンロード処理
	var downloader = s3manager.NewDownloader(sess)

	_, err = downloader.Download(
		tmpfile,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		println("file download error")
	}

	return tmpfile, err
}

func handler(ctx context.Context, req events.S3Event) error {
	return nil
}

func main() {
	lambda.Start(handler)
}
