package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	var regionalBucket = flag.String("regional-bucket", "", "input regional-bucket")
	var azBucket = flag.String("az-bucket", "", "input az-bucket")
	flag.Parse()

	if aws.ToString(regionalBucket) == "" || aws.ToString(azBucket) == "" {
		log.Fatal("missing bucket name\n", "go run main.go -regional-bucket xxxx-regional-bucket -az-bucket xxxx-az-bucket--apne1-az4--x-s3")
	}

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	cfg.Region = "ap-northeast-1"
	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	var fileNames []string = []string{"1mb-file.txt", "5mb-file.txt", "10mb-file.txt", "50mb-file.txt"}
	var files map[string]*os.File = make(map[string]*os.File)

	for _, name := range fileNames {
		file, err := os.Open(name)
		log.Println("opening file", name)
		if err != nil {
			log.Println("Failed opening file", name, err)
		}
		files[name] = file
	}
	for fileName, file := range files {
		uploadAzBucket(file, fileName, uploader, azBucket)
		uploadRegionBucket(file, fileName, uploader, regionalBucket)
	}

}

func uploadRegionBucket(file *os.File, fileName string, uploader *manager.Uploader, bucket *string) {
	startTimeUploadRegionBucket := time.Now()
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: bucket,
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		log.Println("Failed upload file", err)
	}
	endTimeUploadRegionBucket := time.Now()
	uploadTime := endTimeUploadRegionBucket.Sub(startTimeUploadRegionBucket)
	if err == nil {
		log.Printf("upload file %s to %s successful time：%s\n", fileName, aws.ToString(bucket), uploadTime)
	}
}

func uploadAzBucket(file *os.File, fileName string, uploader *manager.Uploader, bucket *string) {
	startTimeUploadAzBucket := time.Now()
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: bucket,
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		log.Println("Failed upload file", err)
	}
	endTimeUploadAzBucket := time.Now()
	uploadTime := endTimeUploadAzBucket.Sub(startTimeUploadAzBucket)

	if err == nil {
		log.Printf("upload file %s to %s successful time：%s\n", fileName, aws.ToString(bucket), uploadTime)
	}
}
