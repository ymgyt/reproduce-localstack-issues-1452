package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	bucket = "ymgyt-localstack-repro"
	region =  "ap-northeast-1"
	endpoint = "http://localhost:4572"
	uploadFile = "gopher.png"
)

func main() {
	client := newS3Client()
	createBucketIfNotExist(client)
	putObject(client)

	if err := deleteObject(client); err != nil {
		fmt.Fprintf(os.Stderr,"delete error %s\n", err.Error())
	}
}

func createBucketIfNotExist(client *s3.S3) {
	output, err := client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(region),
		},
	})

	if err == nil {
		fmt.Printf("bucket %s successfully created at %s 👌\n", bucket, aws.StringValue(output.Location))
	} else {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) {
			switch awsErr.Code() {
				case s3.ErrCodeBucketAlreadyExists, s3.ErrCodeBucketAlreadyOwnedByYou:
				err = nil // nolint
			}
		}
		if err != nil {
			panic(err)
		}
	}
}

func putObject(client *s3.S3) {
	f, err := os.Open(uploadFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = client.PutObject(&s3.PutObjectInput{
		Body:   f,
		Bucket: aws.String(bucket),
		Key:    aws.String(uploadFile),
	})
	if err != nil {
		panic(err)
	}
}

func deleteObject(client *s3.S3) error {
	_, err := client.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &s3.Delete{
			Objects:[]*s3.ObjectIdentifier{
				{ Key:       aws.String(uploadFile)},
			},
		},
	})
	return err
}

func newS3Client() *s3.S3 {
	cfg := aws.NewConfig().
		WithRegion(region).
		WithCredentials(credentials.NewStaticCredentials(
			"dummy",
			"dummy",
			"",
		)).
		WithEndpoint(endpoint).
		WithDisableSSL(true).
		WithS3ForcePathStyle(true)

	sess, err := session.NewSession(cfg)
	if err != nil {
		panic(err)
	}
	return s3.New(sess)
}