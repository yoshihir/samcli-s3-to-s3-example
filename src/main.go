package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type SampleData struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
}

type SampleConvertData struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
	Time  string `json:"time"`
}

func createSession() *session.Session {
	var sess = session.Must(session.NewSession(&aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(os.Getenv("REGION")),
		Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
	}))
	return sess
}

func s3Upload(buf bytes.Buffer) (*s3manager.UploadOutput, error) {
	sess := createSession()

	var uploader = s3manager.NewUploader(sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("TARGET_S3")),
		Key:    aws.String("example-convert.json.gz"),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		fmt.Println("failed to upload file")
		return nil, err
	}

	return result, err
}

func compress(w io.Writer, convertData []SampleConvertData) error {
	b, _ := json.Marshal(convertData)
	gw, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	gw.Write(b)
	defer gw.Close()
	return err
}

func convert(data []SampleData, time string) ([]SampleConvertData, error) {
	var dataConvert []SampleConvertData
	for _, d := range data {
		dataConvert = append(dataConvert, SampleConvertData{
			Id:    d.Id,
			Value: d.Value,
			Time:  time,
		})
	}
	return dataConvert, nil
}

func extract(file *os.File) ([]SampleData, error) {
	gzipReader, _ := gzip.NewReader(file)
	defer gzipReader.Close()

	raw, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		fmt.Println(err.Error())
	}

	var data []SampleData
	err = json.Unmarshal(raw, &data)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data, err
}

func s3Download(bucket string, key string) (f *os.File, err error) {
	sess := createSession()

	tmpFile, _ := ioutil.TempFile("/tmp", "srctmp_")
	defer os.Remove(tmpFile.Name())

	var downloader = s3manager.NewDownloader(sess)

	_, err = downloader.Download(
		tmpFile,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		fmt.Println("file download error")
	}

	return tmpFile, err
}

func handler(ctx context.Context, req events.S3Event) error {
	bucketName := req.Records[0].S3.Bucket.Name
	key := req.Records[0].S3.Object.Key
	file, err := s3Download(bucketName, key)
	if err != nil {
		fmt.Println("Error failed to s3 download")
		return err
	}
	data, err := extract(file)
	if err != nil {
		fmt.Println("Error failed to extract")
		return err
	}
	timeNow := time.Now().String()
	convertData, err := convert(data, timeNow)
	if err != nil {
		fmt.Println("Error failed to convert")
		return err
	}
	var buf bytes.Buffer
	err = compress(&buf, convertData)
	if err != nil {
		fmt.Println("Error failed compress")
		return err
	}
	_, err = s3Upload(buf)
	if err != nil {
		fmt.Println("Error failed to s3 upload")
		return err
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
