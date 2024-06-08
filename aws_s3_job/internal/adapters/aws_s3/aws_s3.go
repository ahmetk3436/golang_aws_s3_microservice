package aws_s3

import (
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/aws_s3_job/config"
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/aws_s3_job/internal/application/core/domain"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"path/filepath"
)

type AwsAdapter struct {
}

func NewAdapter() *AwsAdapter {
	return &AwsAdapter{}
}

func (a *AwsAdapter) LoginAws() (*session.Session, error) {
	awsRegion := config.GetAwsRegion()
	awsAccessKeyID := config.GetAwsAccessKey()
	awsSecretAccessKey := config.GetAwsSecretAccessKey()

	awsCredentials := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(awsRegion),
			Credentials: awsCredentials,
		},
		SharedConfigState: session.SharedConfigEnable,
	}))

	return sess, nil
}

func (a *AwsAdapter) Download(sess *session.Session, bucketName string, fileDataChan chan<- *domain.FileData) error {
	svc := s3.New(sess)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucketName)})
	if err != nil {
		return err
	}

	downloadDir := "./downloads"
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return err
	}

	downloader := s3manager.NewDownloader(sess)

	for _, item := range resp.Contents {
		key := item.Key
		downloadFilePath := filepath.Join(downloadDir, *key)

		file, err := os.Create(downloadFilePath)
		if err != nil {
			fileDataChan <- &domain.FileData{Key: key, Success: false, Error: err}
			continue
		}

		_, err = downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    key,
			})
		if err != nil {
			fileDataChan <- &domain.FileData{Key: key, Success: false, Error: err}
			err := file.Close()
			if err != nil {
				fmt.Printf("Error closing file: %v\n", err)
				return err
			}
			err = os.Remove(downloadFilePath)
			if err != nil {
				fmt.Printf("Error removing file: %s\n", downloadFilePath)
				return err
			}
			continue
		}

		err = file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", err)
			return err
		}

		fileDataChan <- &domain.FileData{Key: key, Success: true, Error: nil}
	}

	return nil
}
