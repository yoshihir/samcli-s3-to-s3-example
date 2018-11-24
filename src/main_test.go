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
		println("Test Handler...")
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
		println("Test S3Download...")
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

//
//
//t.Run("Unable to get IP", func(t *testing.T) {
//	DefaultHTTPGetAddress = "http://127.0.0.1:12345"
//
//	_, err := handler(events.APIGatewayProxyRequest{})
//	if err == nil {
//		t.Fatal("Error failed to trigger with an invalid request")
//	}
//})
//
//t.Run("Non 200 Response", func(t *testing.T) {
//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(500)
//	}))
//	defer ts.Close()
//
//	DefaultHTTPGetAddress = ts.URL
//
//	_, err := handler(events.APIGatewayProxyRequest{})
//	if err != nil && err.Error() != ErrNon200Response.Error() {
//		t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
//	}
//})
//
//t.Run("Unable decode IP", func(t *testing.T) {
//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(500)
//	}))
//	defer ts.Close()
//
//	DefaultHTTPGetAddress = ts.URL
//
//	_, err := handler(events.APIGatewayProxyRequest{})
//	if err == nil {
//		t.Fatal("Error failed to trigger with an invalid HTTP response")
//	}
//})
//
//t.Run("Successful Request", func(t *testing.T) {
//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(200)
//		fmt.Fprintf(w, "127.0.0.1")
//	}))
//	defer ts.Close()
//
//	DefaultHTTPGetAddress = ts.URL
//
//	_, err := handler(events.APIGatewayProxyRequest{})
//	if err != nil {
//		t.Fatal("Everything should be ok")
//	}
//})
