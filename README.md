# Demo Directory bucket vs regional bucket put object
- go version: `v1.21.x`


## [Directory bucket](https://docs.aws.amazon.com/zh_tw/AmazonS3/latest/userguide/directory-buckets-overview.html)
There are two types of Amazon S3 buckets, general purpose buckets and directory buckets. Choose the bucket type that best fits your application and performance requirements:

General purpose buckets are the original S3 bucket type and are recommended for most use cases and access patterns. General purpose buckets also allow objects that are stored across all storage classes, except S3 Express One Zone.

Directory buckets use the S3 Express One Zone storage class, which is recommended if your application is performance sensitive and benefits from single-digit millisecond PUT and GET latencies.

You can create up to 10 directory buckets in each of your AWS accounts, with no limit on the number of objects you can store in a bucket. Your bucket quota is applied to each Region in your AWS account. If your application requires increasing this limit, contact AWS support. For more information, visit the [Service Quotas Console](https://console.aws.amazon.com/servicequotas/home/services/s3/quotas/).

## Directory bucket names
A directory bucket name consists of a base name that you provide and a suffix that contains the ID of the Availability Zone that your bucket is located in. Directory bucket names must use the following format and follow the naming rules for directory buckets:
```bash
base-name--azid--x-s3
```
For more information, see Directory bucket naming rules.

### Create bucket to upload
```bash
cd create-bucket/
go mod download

Example:
go run main.go -bucket neil-demo-s3-bucket-upload-testing

--- example output ---
bucket created at http://neil-demo-s3-bucket-upload-testing.s3.amazonaws.com/
bucket created at https://neil-demo-s3-bucket-upload-testing--apne1-az4--x-s3.s3express-apne1-az4.ap-northeast-1.amazonaws.com/
Execute the following instructions to declare variables, REGIONAL_BUCKET_NAME and AZ_BUCKET_NAME:
export REGIONAL_BUCKET_NAME=neil-demo-s3-bucket-upload-testing
export AZ_BUCKET_NAME=neil-demo-s3-bucket-upload-testing--apne1-az4--x-s3
```

### Create file to upload
- [mkfile](https://ss64.com/bash/mkfile.html): mkfile creates one or more files that are suitable for use as NFS- mounted swap areas. The sticky bit is set, and the file is padded with zeroes by default.

```bash
cd upload-test/
mkfile 1m 1mb-file.txt
mkfile 5m 5mb-file.txt
mkfile 10m 10mb-file.txt
mkfile 50m 50mb-file.txt

```

### Upload file testing
```bash
cd upload-test/

# check file created
ls -l *.txt

--- example output ---
.rw------- neil.kuan staff  10 MB Mon Dec 11 12:31:35 2023  10mb-file.txt
.rw------- neil.kuan staff 1.0 MB Mon Dec 11 12:31:27 2023  1mb-file.txt
.rw------- neil.kuan staff  50 MB Mon Dec 11 12:44:54 2023  50mb-file.txt
.rw------- neil.kuan staff 5.0 MB Mon Dec 11 12:31:24 2023  5mb-file.txt
--- end ---

# run testing
go run main.go -regional-bucket ${REGIONAL_BUCKET_NAME} -az-bucket ${AZ_BUCKET_NAME}

Example:
go run main.go -regional-bucket neil-demo-s3-bucket-upload-testing -az-bucket neil-demo-s3-bucket-upload-testing--apne1-az4--x-s3

--- example output ---
2023/12/11 14:23:41 opening file 1mb-file.txt
2023/12/11 14:23:41 opening file 5mb-file.txt
2023/12/11 14:23:41 opening file 10mb-file.txt
2023/12/11 14:23:41 opening file 50mb-file.txt
2023/12/11 14:23:41 upload file to bucket  neil-demo-s3-bucket-upload-testing--apne1-az4--x-s3
2023/12/11 14:23:41 upload file to bucket neil-demo-s3-bucket-upload-testing
2023/12/11 14:23:41 upload file 1mb-file.txt to neil-demo-s3-bucket-upload-testing--apne1-az4--x-s3 successful time：2.572792ms
2023/12/11 14:23:41 upload file 5mb-file.txt to neil-demo-s3-bucket-upload-testing--apne1-az4--x-s3 successful time：972.708µs
2023/12/11 14:23:41 upload file 10mb-file.txt to neil-demo-s3-bucket-upload-testing--apne1-az4--x-s3 successful time：1.464167ms
2023/12/11 14:23:41 upload file 50mb-file.txt to neil-demo-s3-bucket-upload-testing--apne1-az4--x-s3 successful time：1.39325ms
2023/12/11 14:23:46 upload file 5mb-file.txt to neil-demo-s3-bucket-upload-testing successful time：4.8290915s
2023/12/11 14:23:50 upload file 10mb-file.txt to neil-demo-s3-bucket-upload-testing successful time：3.714354625s
2023/12/11 14:24:07 upload file 50mb-file.txt to neil-demo-s3-bucket-upload-testing successful time：17.285714167s
2023/12/11 14:24:07 upload file 1mb-file.txt to neil-demo-s3-bucket-upload-testing successful time：534.065542ms
--- end ---
```


# 🤪🤪🤪 Also take a look 🤪🤪🤪: [AWS S3 Pricing](https://aws.amazon.com/tw/s3/pricing/)