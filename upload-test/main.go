package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
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
	if err != nil {
		log.Fatal(err)
	}
	session, err := client.CreateSession(context.TODO(), &s3.CreateSessionInput{
		Bucket: azBucket,
	})
	if err != nil {
		log.Fatal(err)
	}
	sessionCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			aws.ToString(session.Credentials.AccessKeyId),
			aws.ToString(session.Credentials.SecretAccessKey),
			aws.ToString(session.Credentials.SessionToken))),
	)

	sessionClient := s3.NewFromConfig(sessionCfg)
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

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go uploadRegionBucket(files, client, wg, regionalBucket)
	go uploadAzBucket(files, sessionClient, wg, azBucket)
	wg.Wait()

}

func uploadRegionBucket(files map[string]*os.File, client *s3.Client, wg *sync.WaitGroup, bucket *string) {
	log.Println("upload file to bucket " + aws.ToString(bucket))

	for fileName, file := range files {
		startTimeUploadRegionBucket := time.Now()
		client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: bucket,
			Key:    aws.String(fileName),
			Body:   file,
		})
		endTimeUploadRegionBucket := time.Now()
		uploadTime := endTimeUploadRegionBucket.Sub(startTimeUploadRegionBucket)

		log.Printf("upload file %s to %s successful time：%s\n", fileName, aws.ToString(bucket), uploadTime)
	}

	defer wg.Done()
}

func uploadAzBucket(files map[string]*os.File, sessionClient *s3.Client, wg *sync.WaitGroup, bucket *string) {
	log.Println("upload file to bucket ", aws.ToString(bucket))

	for fileName, file := range files {
		startTimeUploadAzBucket := time.Now()
		sessionClient.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: bucket,
			Key:    aws.String(fileName),
			Body:   file,
		})
		endTimeUploadAzBucket := time.Now()
		uploadTime := endTimeUploadAzBucket.Sub(startTimeUploadAzBucket)

		log.Printf("upload file %s to %s successful time：%s\n", fileName, aws.ToString(bucket), uploadTime)
	}

	defer wg.Done()
}
