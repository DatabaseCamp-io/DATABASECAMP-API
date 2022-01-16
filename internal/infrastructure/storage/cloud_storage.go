package storage

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

type cloudStorageService struct {
	bucketName string
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
		bucketName: os.Getenv("CLOUD_BUCKET_NAME"),
	}
}

func (c cloudStorageService) GetFileLink(objectName string) (string, error) {
	jsonKey, err := ioutil.ReadFile("service_account.json")
	if err != nil {
		return "", err
	}

	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", err
	}

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         http.MethodGet,
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}

	url, err := storage.SignedURL(c.bucketName, objectName, opts)

	return url, err
}
