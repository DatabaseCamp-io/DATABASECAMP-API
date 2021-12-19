package services

import (
	"context"
	"os"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/storage"

	credentialspb "google.golang.org/genproto/googleapis/iam/credentials/v1"
)

type cloudStorageService struct {
	bucketName     string
	googleAccessID string
}

var cloudStorageInstantiated *cloudStorageService = nil

func GetCloudStorageServiceInstance() *cloudStorageService {
	if cloudStorageInstantiated == nil {
		cloudStorageInstantiated = initGoogleCloudStorageService()
	}
	return cloudStorageInstantiated
}

func initGoogleCloudStorageService() *cloudStorageService {
	return &cloudStorageService{
		bucketName:     os.Getenv("CLOUD_BUCKET_NAME"),
		googleAccessID: os.Getenv("CLOUD_GOOGLE_ACCESS_ID"),
	}
}

func (c cloudStorageService) signBytes(b []byte) ([]byte, error) {
	ctx := context.Background()

	client, _ := credentials.NewIamCredentialsClient(ctx)

	req := &credentialspb.SignBlobRequest{
		Payload: b,
		Name:    c.googleAccessID,
	}

	resp, err := client.SignBlob(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.SignedBlob, err
}

func (c cloudStorageService) GetFileLink(objectName string) (string, error) {
	options := storage.SignedURLOptions{
		GoogleAccessID: c.googleAccessID,
		SignBytes:      c.signBytes,
	}

	url, err := storage.SignedURL(c.bucketName, objectName, &options)

	return url, err
}
