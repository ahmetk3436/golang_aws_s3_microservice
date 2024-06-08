package ports

import (
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/aws_s3_job/internal/application/core/domain"
	"context"
)

type FileJobPort interface {
	ReadFile(ctx context.Context, filePath string) ([]*domain.Product, error)
	DeleteFile(filePath string)
}
