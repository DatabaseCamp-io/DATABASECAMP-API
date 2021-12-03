package services

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type awsService struct {
	Sess *session.Session
}

type IAwsService interface {
	GetFileLink(imagekey string) (string, error)
}

var awsInstantiated *awsService = nil

// Get AWS service instance
func GetAwsServiceInstance() *awsService {
	if awsInstantiated == nil {
		sess, _ := session.NewSession(getConfig())
		awsInstantiated = &awsService{Sess: sess}
	}
	return awsInstantiated
}

// Get cradentials
func getCredentials() *credentials.Credentials {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	return credentials.NewStaticCredentials(
		accessKeyID,
		secretAccessKey,
		"",
	)
}

// Get config
func getConfig() *aws.Config {
	region := os.Getenv("AWS_REGION")
	return &aws.Config{
		Region:      aws.String(region),
		Credentials: getCredentials(),
	}
}

// Get object input
func (s awsService) getObjectInput(imageKey string) *s3.GetObjectInput {
	bucket := os.Getenv("AWS_BUCKET_NAME")
	return &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(imageKey),
	}
}

// Get file link
func (s awsService) GetFileLink(imageKey string) (string, error) {
	svc := s3.New(s.Sess)
	req, _ := svc.GetObjectRequest(s.getObjectInput(imageKey))
	expireTime := 60 * time.Minute
	url, err := req.Presign(expireTime)
	return url, err
}
