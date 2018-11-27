package main

import (
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
	"io/ioutil"
	"os"
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

func s3Upload(file *os.File) error {
	var sess = session.Must(session.NewSession(&aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(os.Getenv("REGION")),
		Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
	}))

	var uploader = s3manager.NewUploader(sess)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("bucket-example-convert"),
		Key:    aws.String("example-convert.json.gz"),
		Body:   file,
	})
	if err != nil {
		fmt.Println("failed to upload file")
		return err
	}

	return nil
}

func compress(convertData []SampleConvertData) (*os.File, error) {
	b, _ := json.Marshal(convertData)

	tmpfile, _ := ioutil.TempFile("/tmp", "srctmp_")
	defer os.Remove(tmpfile.Name())

	writer := gzip.NewWriter(tmpfile)
	writer.Write([]byte(b))
	writer.Close()

	return tmpfile, nil
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
		fmt.Println("file download error")
	}

	return tmpfile, err
}

func handler(ctx context.Context, req events.S3Event) error {
	return nil
}

func main() {
	lambda.Start(handler)
}
