package services

// service.aws.go
/**
 * 	This file is a part of services, used to connect to the Amason S3
 */

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

/**
 * This class manage connection to Amazon S3
 */
type awsService struct {
	Sess *session.Session
}

// Instance of awsService class for singleton pattern
var awsInstantiated *awsService = nil

/**
 * Constructor creates a new awsService instance
 *
 * @return 	instance of awsService
 */
func GetAwsServiceInstance() *awsService {
	if awsInstantiated == nil {
		sess, _ := session.NewSession(getConfig())
		awsInstantiated = &awsService{Sess: sess}
	}
	return awsInstantiated
}

/**
 * Get credential of the amazon s3
 *
 * @return credential of the amazon s3
 */
func getCredentials() *credentials.Credentials {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	return credentials.NewStaticCredentials(
		accessKeyID,
		secretAccessKey,
		"",
	)
}

/**
 * Get config of the amazon s3
 *
 * @return config of the amazon s3
 */
func getConfig() *aws.Config {
	region := os.Getenv("AWS_REGION")
	return &aws.Config{
		Region:      aws.String(region),
		Credentials: getCredentials(),
	}
}

/**
 * Get object Input of the amazon s3
 *
 * @param 	imagekey  object key for getting file link
 *
 * @return object Input of the amazon s3
 */
func (s awsService) getObjectInput(imageKey string) *s3.GetObjectInput {
	bucket := os.Getenv("AWS_BUCKET_NAME")
	return &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(imageKey),
	}
}

/**
 * Get File link from amazon s3
 *
 * @param 	imagekey  object key for getting file link
 *
 * @return file link
 * @return the error of getting file link
 */
func (s awsService) GetFileLink(imageKey string) (string, error) {
	svc := s3.New(s.Sess)
	req, _ := svc.GetObjectRequest(s.getObjectInput(imageKey))
	expireTime := 60 * time.Minute
	url, err := req.Presign(expireTime)
	return url, err
}
