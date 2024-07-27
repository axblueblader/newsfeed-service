package storage

import "fmt"

type objectStorage struct {
	// we can have a Cloud Client interface in the future and pass it in here
	AwsClient interface{}
}

func NewObjectStorage(awsClient interface{}) ObjectStorage {
	return &objectStorage{AwsClient: awsClient}
}

type ObjectStorage interface {
	GenerateSignedUrl(bucket string, path string) string
}

func (s objectStorage) GenerateSignedUrl(bucket string, path string) string {
	return fmt.Sprintf("https://cloud.provider/%s/%s?signature=xxx&expiresIn=123", bucket, path)
}
