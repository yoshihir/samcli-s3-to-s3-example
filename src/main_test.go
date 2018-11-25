package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// test前処理
	println("before all...")

	os.Setenv("REGION", "ap-northeast-1")
	os.Setenv("S3_ENDPOINT", "http://localhost:4572")

	var sess = session.Must(session.NewSession(&aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(os.Getenv("REGION")),
		Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
	}))
	var creater = s3.New(sess)
	var uploader = s3manager.NewUploader(sess)

	input := &s3.CreateBucketInput{
		Bucket: aws.String("bucket-example"),
	}

	_, err := creater.CreateBucket(input)
	if err != nil {
		//if aerr, ok := err.(awserr.Error); ok {
		//	switch aerr.Code() {
		//	case s3.ErrCodeBucketAlreadyExists:
		//		fmt.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())
		//	case s3.ErrCodeBucketAlreadyOwnedByYou:
		//		fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
		//	default:
		//		fmt.Println(aerr.Error())
		//	}
		//} else {
		//	fmt.Println(err.Error())
		//}
		//return
	}

	up, err := os.Open("../example.json.gz")
	if err != nil {
		fmt.Println("failed to open file")
		return
	}

	gzip.NewWriter(up).Flush()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("bucket-example"),
		Key:    aws.String("example.json.gz"),
		Body:   up,
	})
	if err != nil {
		fmt.Println("failed to upload file")
		return
	}

	// test実行
	code := m.Run()
	// test後実行
	println("after all...")
	os.Exit(code)
}

func TestHandler(t *testing.T) {
	t.Run("handler input test", func(t *testing.T) {
		raw, err := ioutil.ReadFile("../event_file.json")
		if err != nil {
			t.Fatal("Error failed to event file load")
		}
		var event events.S3Event
		json.Unmarshal(raw, &event)
		err = handler(context.Background(), event)
		if err != nil {
			t.Fatal("Error failed to s3 event")
		}
		println("Test handler...")
	})
}

func TestS3Download(t *testing.T) {
	t.Run("s3 download test", func(t *testing.T) {
		tmpFile, err := s3Download("bucket-example", "example.json.gz")
		if err != nil {
			t.Fatal("Error failed to s3 download")
		}
		if tmpFile.Name() == "" {
			t.Errorf("got: %v\nwant: %v", "", "/tmp/srctmp_*********")
		}
		println("Test s3Download...")
	})
}

func TestExtract(t *testing.T) {
	t.Run("extract", func(t *testing.T) {
		file, _ := os.Open("../example.json.gz")
		defer file.Close()
		actual, err := extract(file)
		if err != nil {
			t.Fatal("Error failed to extract")
		}
		expected := "abcdefgh"
		if actual[0].Value != expected {
			t.Errorf("got: %v\nwant: %v", actual[0].Value, expected)
		}
		println("Test extract...")
	})
}

func TestConvert(t *testing.T) {
	t.Run("convert", func(t *testing.T) {
		data := []SampleData{{12345678, "abcdefgh"}}
		expected := time.Now().String()

		convertData, err := convert(data, expected)
		if err != nil {
			t.Fatal("Error failed to convert")
		}
		if convertData[0].Time != expected {
			t.Errorf("got: %v\nwant: %v", convertData[0].Time, expected)
		}

		println("Test convert...")
	})
}

func TestCompress(t *testing.T) {
	t.Run("compress", func(t *testing.T) {
		data := []SampleConvertData{
			{12345678, "abcdefgh", time.Now().String()},
			{23456781, "bcdefgha", time.Now().String()}}
		_, err := compress(data)
		if err != nil {
			t.Fatal("Error failed to compress")
		}

		// ここのtestは誰かに相談する
		println("Test compress...")
	})
}
