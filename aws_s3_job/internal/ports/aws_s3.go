package ports

import (
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/aws_s3_job/internal/application/core/domain"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Aws interface {
	Download(sess *session.Session, bucketName string, fileDataChan chan<- *domain.FileData) error
	LoginAws() (*session.Session, error)
}
