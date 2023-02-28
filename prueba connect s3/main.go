package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Client struct {
	Region string
	Sess   *session.Session
	Svc    *s3.S3
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func NewS3Client(region string) *S3Client {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		exitErrorf("Failed to create AWS session: %v", err)
	}

	return &S3Client{
		Region: region,
		Sess:   sess,
		Svc:    s3.New(sess),
	}
}

func (c *S3Client) ListBuckets() {
	result, err := c.Svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Failed to list buckets: %v", err)
	}

	fmt.Println("Buckets:")
	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}

func (c *S3Client) UploadFile(localFile, bucket, keyName string) error {
	f, err := os.Open(localFile)
	if err != nil {
		return fmt.Errorf("Failed to open file %s: %v", localFile, err)
	}
	defer f.Close()

	uploader := s3manager.NewUploader(c.Sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(keyName),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("Failed to upload file %s: %v", localFile, err)
	}

	return nil
}

// func (c *S3Client) GetObjectURL(bucket, keyName string) (string, error) {
// 	req, _ := c.Svc.GetObjectRequest(&s3.GetObjectInput{
// 		Bucket: aws.String(bucket),
// 		Key:    aws.String(keyName),
// 	})
// 	urlStr, err := req.Presign(0)
// 	if err != nil {
// 		return "", fmt.Errorf("Failed to generate URL for %s/%s: %v", bucket, keyName, err)
// 	}

//		return urlStr, nil
//	}
func main() {
	region := "us-east-2"
	bucket := "practica1-g8-imagenes"
	keyName := "Fotos_Perfil/mydog2.jpg"
	// create a new S3 client
	s3Client := NewS3Client(region)

	// list all S3 buckets
	s3Client.ListBuckets()

	// upload a local file to S3\
	fmt.Println("Uploading %s to S3 bucket %s...\n", keyName, bucket)
	err := s3Client.UploadFile("C:/Users/socop/OneDrive/Escritorio/mydog.jpg", bucket, keyName)
	if err != nil {
		exitErrorf("Failed to upload file to S3: %v", err)
	}
	fmt.Println("Upload successful!")

	url := "https://practica1-g8-imagenes.s3.us-east-2.amazonaws.com/" + keyName
	if err != nil {
		exitErrorf("Failed to generate URL for S3 object: %v", err)
	}
	fmt.Println(url)
}
