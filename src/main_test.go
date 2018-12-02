package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
	os.Setenv("TARGET_S3", "bucket-example-convert")

	var sess = session.Must(session.NewSession(&aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(os.Getenv("REGION")),
		Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
	}))
	var creater = s3.New(sess)
	var uploader = s3manager.NewUploader(sess)

	_, err := creater.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String("bucket-example"),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}
	_, err = creater.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String("bucket-example-convert"),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	up, err := os.Open("./testdata/example.json.gz")
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

func TestS3Upload(t *testing.T) {
	t.Run("upload", func(t *testing.T) {
		var buf bytes.Buffer
		result, err := s3Upload(buf)
		if err != nil {
			t.Fatal("Error failed to s3upload")
		}
		if result.Location == "" {
			t.Errorf("got: %v\nwant: %v", result.UploadID, "")
		}
		fmt.Println("Test s3upload...")
	})
}

func TestCompress(t *testing.T) {
	t.Run("compress", func(t *testing.T) {
		data := []SampleConvertData{
			{12345678, "abcdefgh", time.Now().String()},
			{23456781, "bcdefgha", time.Now().String()}}

		var buf bytes.Buffer
		err := compress(&buf, data)
		if err != nil {
			t.Fatal("Error failed to compress")
		}
		if len(buf.Bytes()) == 0 {
			t.Fatal("Error failed to compress")
			t.Errorf("got: %v\nwant: %v", buf.Bytes(), 0)
		}
		fmt.Println("Test compress...")
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
		fmt.Println("Test convert...")
	})
}

func TestExtract(t *testing.T) {
	t.Run("extract", func(t *testing.T) {
		file, _ := os.Open("./testdata/example.json.gz")
		defer file.Close()
		actual, err := extract(file)
		if err != nil {
			t.Fatal("Error failed to extract")
		}
		expected := "abcdefgh"
		if actual[0].Value != expected {
			t.Errorf("got: %v\nwant: %v", actual[0].Value, expected)
		}
		fmt.Println("Test extract...")
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
		fmt.Println("Test s3Download...")
	})
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
		fmt.Println("Test handler...")
	})
}
