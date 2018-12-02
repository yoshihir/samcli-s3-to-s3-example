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

func s3Upload(buf bytes.Buffer) (*s3manager.UploadOutput, error) {
	region := os.Getenv("REGION")
	endpoint := os.Getenv("S3_ENDPOINT")

	var sess = session.Must(session.NewSession(&aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
	}))

	var uploader = s3manager.NewUploader(sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		// TODO: 検証環境のbucket名をbucket-example-convert-stagingに変えたいので、環境変数から取れるようにする
		Bucket: aws.String("bucket-example-convert-staging"),
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
	// Write gzipped data to the client
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
	region := os.Getenv("REGION")
	endpoint := os.Getenv("S3_ENDPOINT")
	var sess = session.Must(session.NewSession(&aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
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
	bucketName := req.Records[0].S3.Bucket.Name
	key := req.Records[0].S3.Object.Key
	file, err := s3Download(bucketName, key)
	if err != nil {
		fmt.Println("Error failed to s3 download")
		return err
	}
	data, err := extract(file)
	time := time.Now().String()
	convertData, err := convert(data, time)
	var buf bytes.Buffer
	err = compress(&buf, convertData)
	if err != nil {
		fmt.Println("Error failed compress")
		return err
	}
	_, err = s3Upload(buf)

	fmt.Println("Success!!")
	return nil
}

func main() {
	lambda.Start(handler)
}
