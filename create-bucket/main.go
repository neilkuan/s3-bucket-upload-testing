package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"flag"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func main() {
	var bucket = flag.String("bucket", "s3-bucket-upload-testing", "create bucket")
	flag.Parse()
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	cfg.Region = "ap-northeast-1"
	runCreateBucket(s3.NewFromConfig(cfg), aws.ToString(bucket))
	runCreateAzBucket(s3.NewFromConfig(cfg), aws.ToString(bucket))
	fmt.Println("Execute the following instructions to declare variables, REGIONAL_BUCKET_NAME and AZ_BUCKET_NAME:")
	fmt.Printf("export REGIONAL_BUCKET_NAME=%s\n", aws.ToString(bucket))
	fmt.Printf("export AZ_BUCKET_NAME=%s\n", aws.ToString(bucket)+"--apne1-az4--x-s3")
}

func runCreateAzBucket(c *s3.Client, bucket string) {
	resp, err := c.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(bucket + "--apne1-az4--x-s3"),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			Location: &types.LocationInfo{
				Name: aws.String("apne1-az4"),
				Type: types.LocationTypeAvailabilityZone,
			},
			Bucket: &types.BucketInfo{
				DataRedundancy: types.DataRedundancySingleAvailabilityZone,
				Type:           types.BucketTypeDirectory,
			},
		},
	})
	var terr *types.BucketAlreadyOwnedByYou
	if errors.As(err, &terr) {
		fmt.Printf("Bucket %s AlreadyOwnedByYou: %s\n", bucket+"--apne1-az4--x-s3", aws.ToString(terr.Message))
		fmt.Printf("...\n")
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("bucket created at %s\n", aws.ToString(resp.Location))
}

func runCreateBucket(c *s3.Client, bucket string) {
	resp, err := c.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint("ap-northeast-1"),
		},
	})
	var terr *types.BucketAlreadyOwnedByYou
	if errors.As(err, &terr) {
		fmt.Printf("Bucket %s AlreadyOwnedByYou: %s\n", bucket, aws.ToString(terr.Message))
		fmt.Printf("...\n")
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("bucket created at %s\n", aws.ToString(resp.Location))
}
